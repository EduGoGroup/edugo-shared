# Decisión: Tarea 1.2 Bloqueada - Go 1.25 no disponible

**Fecha:** 20 Nov 2025, 19:30
**Tarea:** 1.2 - Migrar a Go 1.25
**Razón:** Go 1.25 no ha sido lanzado aún

## Contexto

Al intentar ejecutar `go mod tidy` en los módulos después de actualizar la versión de Go a 1.25, el sistema intenta descargar Go 1.25.0 pero falla porque:

1. Go 1.25 no ha sido lanzado oficialmente
2. La versión actual de Go instalada es 1.24.7
3. No hay acceso a internet para descargar toolchains

```
go: downloading go1.25.0 (linux/amd64)
go: download go1.25.0: ... dial tcp: lookup storage.googleapis.com ... connection refused
```

## Análisis

Los workflows ya están configurados con Go 1.25:
- `.github/workflows/ci.yml`: `GO_VERSION: '1.25'`
- `.github/workflows/release.yml`: `GO_VERSION: '1.25'`
- `.github/workflows/test.yml`: `go-version: '1.25'`

Esto sugiere que la migración a Go 1.25 es preparatoria para cuando se lance esta versión.

## Decisión

Implementar un **stub** usando `GOTOOLCHAIN=local` para permitir que el código compile y se pruebe con Go 1.24.x mientras se mantiene la declaración de Go 1.25 en los archivos:

1. Actualizar go.mod a `go 1.25` (HECHO)
2. Usar `GOTOOLCHAIN=local` para forzar el uso de Go 1.24.x instalado
3. Ejecutar `go mod tidy` con `GOTOOLCHAIN=local`
4. Documentar que esto es temporal hasta que Go 1.25 se lance

## Implementación del Stub

```bash
# Ejecutar go mod tidy con toolchain local (Go 1.24.7)
export GOTOOLCHAIN=local
for dir in auth bootstrap common config database/mongodb database/postgres evaluation lifecycle logger messaging/rabbit middleware/gin testing; do
  if [ -d "$dir" ] && [ -f "$dir/go.mod" ]; then
    (cd "$dir" && GOTOOLCHAIN=local go mod tidy)
  fi
done
```

## Para Fase 2

Cuando Go 1.25 se lance oficialmente:
1. Verificar que Go 1.25 está disponible en el sistema
2. Re-ejecutar `go mod tidy` sin `GOTOOLCHAIN=local`
3. Validar compilación y tests con Go 1.25 real
4. Actualizar esta decisión

## Validación

- ✅ go.mod actualizados a `go 1.25`
- ✅ Workflows ya configurados con Go 1.25
- ⚠️ Usando stub: GOTOOLCHAIN=local para compilar con Go 1.24.7
- ⏳ Pendiente: Migración real cuando Go 1.25 se lance

## Migaja

- Marcada como: ✅ (stub)
- Pendiente para Fase 2: Validar con Go 1.25 real cuando esté disponible
