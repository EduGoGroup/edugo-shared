# Priority Module Deep Dive: `bootstrap`

## Por qué `bootstrap` es el módulo donde más conviene invertir

`bootstrap` concentra el orquestado de recursos críticos (logger, PostgreSQL, MongoDB, RabbitMQ, S3) y es dependencia directa del `edugo-worker`.

Señales objetivas:

- Alto volumen: 26 archivos Go (entre los más altos del repo).
- Mezcla de responsabilidades en raíz del módulo: factories, init, cleanup, health checks, extractores de config.
- Superficie operativa sensible: si falla aquí, el Worker no inicia correctamente.

## Qué implica invertir en `bootstrap`

## 1) Beneficio técnico esperado

- Menor complejidad de mantenimiento.
- Menos duplicación en inicialización/cleanup.
- Errores de startup más consistentes y trazables.
- Mejor base para agregar nuevos recursos sin crecer deuda técnica.

## 2) Implicación para consumidores

### `edugo-worker` (impacto principal)

- Es el consumidor directo más sensible.
- Requiere validación obligatoria en cada paso de refactor.
- Si se preserva API (`Bootstrap`, `Factories`, `Resources`, `BootstrapOptions`), el impacto debe ser compatible.

### IAM/Admin/Mobile

- No dependen directamente de `bootstrap` según matriz actual.
- Impacto esperado nulo, salvo cambios transversales en módulos compartidos que `bootstrap` use.

### `kmp_new` / `apple_new`

- Impacto indirecto solo si el Worker cambia comportamiento funcional observable.

## 3) Ejemplos reales donde invertir primero

## A) Inicializadores repetidos

Archivos:

- `bootstrap/init_postgresql.go`
- `bootstrap/init_mongodb.go`
- `bootstrap/init_rabbitmq.go`
- `bootstrap/init_s3.go`

Patrón actual: flujo casi idéntico (validación, extracción de config, init, log, cleanup).

Acción:

- Crear template interno de inicialización para reducir duplicación.

## B) Cleanup registrars repetidos

Archivo:

- `bootstrap/cleanup_registrars.go`

Acción:

- Abstraer lógica común de registro tipado y mantener wrappers específicos.

## C) Wrappers de config potencialmente desalineados

Archivo:

- `bootstrap/interfaces.go` (`PostgreSQLConfig`, `MongoDBConfig`, `RabbitMQConfig`)

Comparar con:

- `database/postgres/config.go`
- `database/mongodb/config.go`
- `messaging/rabbit/config.go`

Acción:

- Definir adaptadores explícitos y evitar drift de configuración.

## Ruta de implementación recomendada para `bootstrap`

1. Reorganización interna a `internal/` sin tocar API pública.
2. Consolidación de patrones init/cleanup.
3. Adaptadores de config explícitos (`bootstrap -> módulos`).
4. Tests de regresión dirigidos a startup/cleanup del Worker.

## Riesgo y control

Riesgo:

- Medio, por centralidad operativa.

Controles:

- Mantener contrato público estable.
- Validar con tests de integración del módulo y del Worker.
- Release con changelog detallado y checklist de compatibilidad.
