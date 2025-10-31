# 📊 Coverage Configuration - EduGo Shared

Esta configuración de cobertura está diseñada para **enfocar los esfuerzos de testing en el código que realmente importa** y excluir archivos que son principalmente configuración, constantes o constructores simples.

## 🎯 Filosofía de Cobertura

### ✅ **Paquetes CRÍTICOS** (Requieren alta cobertura)
Estos paquetes contienen **lógica de negocio crítica** y deben tener buena cobertura:

- **`pkg/auth/`** - Autenticación JWT, tokens, validaciones de seguridad
- **`pkg/database/`** - Conexiones, transacciones, manejo de errores de BD
- **`pkg/logger/`** - Configuración de logging, formatos, niveles
- **`pkg/messaging/`** - Publisher/Consumer, manejo de mensajes
- **`pkg/validator/`** - Validaciones de entrada, reglas de negocio

**Meta de cobertura:** 80%+ en estos paquetes

### ❌ **Paquetes EXCLUIDOS** (No afectan cobertura crítica)
Estos archivos son principalmente configuración y no requieren testing exhaustivo:

- **`pkg/config/`** - Solo getters de variables de entorno
- **`pkg/errors/`** - Solo constructores de errores (sin lógica compleja)
- **`pkg/types/enum/`** - Constantes y métodos simples (IsValid, String)

### 🤔 **Paquetes OPCIONALES** (Recomendado pero no crítico)
- **`pkg/types/uuid.go`** - Wrapper de UUID con algo de lógica de serialización

## 🚀 Comandos de Cobertura

### Cobertura Crítica (Recomendado)
```bash
# Solo paquetes críticos - esto es lo que importa
make test-coverage-critical

# O más corto:
make test-coverage
```

### Cobertura Completa (Informativo)
```bash
# Todos los paquetes - solo para información completa
make test-coverage-all
```

### Ver Configuración
```bash
# Mostrar qué se incluye/excluye
make coverage-info
```

## 📈 Interpretación de Resultados

### Ejemplo de Output:
```
📊 Cobertura de paquetes críticos:
total: (statements) 85.2%  ← Esta es la métrica importante

📊 Cobertura completa (incluye config/errors/types):
total: (statements) 45.3%  ← Esta baja por archivos excluidos (normal)
```

### ✅ Lo que significa **85.2% en críticos**:
- El código que **realmente necesita testing** tiene buena cobertura
- La lógica de autenticación, base de datos, etc. está bien probada
- **Esta es la métrica que debes optimizar**

### ⚠️ Lo que significa **45.3% completo**:
- Incluye archivos de configuración y constantes (que no necesitan testing)
- **Normal que sea bajo** - no te preocupes por este número
- **Solo informativo**

## 📋 Archivos de Configuración

### `.testcoverage.yml`
Configuración completa con exclusiones y umbrales:
- Threshold objetivo: 80%
- Mínimo aceptable: 70%
- Lista de exclusiones
- Paquetes críticos definidos

### `Makefile`
Comandos separados para:
- Cobertura crítica (`test-coverage-critical`)
- Cobertura completa (`test-coverage-all`)
- Información de configuración (`coverage-info`)

## 🎪 CI/CD Integration

El workflow de GitHub Actions genera **dos reportes HTML**:
1. **`coverage-critical.html`** - Solo paquetes críticos ⭐
2. **`coverage-all.html`** - Todos los paquetes 📋

**Enfócate en el reporte crítico** para tus revisiones de cobertura.

## 💡 Mejores Prácticas

### ✅ DO:
- Enfocar esfuerzos en alcanzar 80%+ en paquetes críticos
- Usar `make test-coverage` (críticos) para desarrollo diario
- Revisar el reporte `coverage-critical.html` en PRs

### ❌ DON'T:
- Preocuparse por la cobertura baja en el reporte completo
- Intentar testear exhaustivamente archivos de configuración
- Perder tiempo escribiendo tests para constantes/enums simples

---

**🎯 Objetivo:** Alcanzar 80%+ de cobertura en paquetes críticos, ignorando archivos de configuración que solo contienen constantes y getters simples.