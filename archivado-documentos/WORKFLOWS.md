# GitHub Actions Workflows - edugo-shared

Este documento describe todos los workflows de CI/CD y cuándo se ejecutan.

---

## 📋 Resumen de Workflows

| Workflow | Archivo | Triggers | Propósito |
|----------|---------|----------|-----------|
| CI Pipeline | `ci.yml` | PR + Push main | Tests y validación en cambios |
| Tests with Coverage | `test.yml` | Manual + PR | Coverage detallado por módulo |
| Release CI/CD | `release.yml` | Tag v* | Release modular automático |
| Sync Main to Dev | `sync-main-to-dev.yml` | Push main + Tag | Sincronización de ramas |

---

## 🔄 CI Pipeline (`ci.yml`)

**Archivo:** `.github/workflows/ci.yml`

### Cuándo se Ejecuta

```yaml
on:
  pull_request:
    branches: [ main, dev ]
  push:
    branches: [ main ]
```

- ✅ Al abrir/actualizar PR a `main` o `dev`
- ✅ Al hacer push directo a `main`
- ❌ NO se ejecuta en push a otras ramas

### Qué Hace

1. **Tests por módulo** (matriz):
   - Compila cada módulo
   - Ejecuta tests unitarios
   - Valida con `-race` (data races)

2. **Compatibilidad Go**:
   - Prueba con Go 1.23, 1.24, 1.25
   - Asegura compatibilidad hacia atrás

3. **Lint** (opcional):
   - golangci-lint en cada módulo
   - `continue-on-error: true` (no bloquea)

### Duración Típica

~3-4 minutos (todos los módulos en paralelo)

---

## 🧪 Tests with Coverage (`test.yml`)

**Archivo:** `.github/workflows/test.yml`

### Cuándo se Ejecuta

```yaml
on:
  workflow_dispatch:  # Manual desde UI
  pull_request:
    branches: [ main, dev ]
```

**IMPORTANTE:**
```yaml
if: github.event_name != 'push'  # ← Evita "fallos fantasma"
```

- ✅ Manualmente desde GitHub UI
- ✅ En PRs a `main` o `dev`
- ❌ NO en push (condición explícita)

### Qué Hace

1. **Coverage por módulo**:
   - Ejecuta tests con `-cover`
   - Genera `coverage.out` por módulo
   - Calcula porcentaje de cobertura

2. **Reportes**:
   - Sube artifacts con coverage
   - Genera resumen en GitHub Step Summary

### Duración Típica

~5-6 minutos (ejecuta tests más completos)

---

## 🚀 Release CI/CD (`release.yml`)

**Archivo:** `.github/workflows/release.yml`

### Cuándo se Ejecuta

```yaml
on:
  push:
    tags:
      - 'v*'  # Ejemplo: v1.0.0, v0.1.2
```

- ✅ Al crear y pushear tag con formato `v*`
- ❌ NO se ejecuta en tags sin `v` prefix

### Qué Hace

1. **Extrae versión** del tag (ej: v1.0.0 → 1.0.0)
2. **Crea GitHub Release** con changelog
3. **Publica instrucciones** de instalación por módulo

### Crear Release Manualmente

```bash
# 1. Crear tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 2. Push tag
git push origin v1.0.0

# 3. El workflow se ejecuta automáticamente
```

---

## 🔄 Sync Main to Dev (`sync-main-to-dev.yml`)

**Archivo:** `.github/workflows/sync-main-to-dev.yml`

### Cuándo se Ejecuta

```yaml
on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
```

- ✅ Después de merge a `main`
- ✅ Después de crear tag de release
- ❌ NO en push a otras ramas

### Qué Hace

1. Verifica si rama `dev` existe (crea si no)
2. Compara commits entre `main` y `dev`
3. Hace merge automático de `main` → `dev`
4. Maneja conflictos (aborta si hay)

### Condiciones Especiales

```yaml
if: "!contains(github.event.head_commit.message, 'chore: sync')"
```

- ❌ NO se ejecuta si el commit ya es un sync (evita loops)

---

## 🎯 Flujo Típico de Trabajo

### Desarrollo de Feature

```
1. Crear rama: feature/nueva-funcionalidad
2. Hacer cambios y commits
3. Crear PR a dev
   ├─> ✅ ci.yml se ejecuta (tests)
   └─> ✅ test.yml se ejecuta (coverage)
4. Review y merge
```

### Merge a Main (Release)

```
1. Crear PR de dev → main
   ├─> ✅ ci.yml se ejecuta
   └─> ✅ test.yml se ejecuta
2. Merge a main
   └─> ✅ ci.yml se ejecuta de nuevo
3. Crear tag v1.0.0
   ├─> ✅ release.yml se ejecuta (crea release)
   └─> ✅ sync-main-to-dev.yml se ejecuta (sync)
```

### Ejecución Manual de Tests

```
1. Ir a Actions en GitHub
2. Seleccionar "Tests with Coverage"
3. Click "Run workflow"
4. Seleccionar rama
5. Click "Run workflow"
```

---

## 🐛 Troubleshooting

### Workflow no se ejecuta

**Síntoma:** Push hecho pero workflow no aparece en Actions.

**Soluciones:**
1. Verificar que el trigger incluye tu evento:
   ```bash
   # Ver triggers de un workflow
   grep -A 5 "^on:" .github/workflows/ci.yml
   ```

2. Verificar que la rama está en el trigger:
   ```yaml
   on:
     push:
       branches: [ main ]  # ← Solo main
   ```

3. Verificar sintaxis YAML:
   ```bash
   python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
   ```

### "Fallos fantasma" en historial

**Síntoma:** Workflow aparece fallando con 0s de duración.

**Causa:** GitHub intenta ejecutar workflow en evento no configurado.

**Solución:** Agregar condición explícita:
```yaml
jobs:
  mi-job:
    if: github.event_name != 'push'  # O el evento a excluir
```

### Workflow tarda mucho

**Síntoma:** Workflow toma >10 minutos.

**Soluciones:**
1. Verificar que usa matriz para paralelización
2. Optimizar caché de Go:
   ```yaml
   - uses: actions/setup-go@v5
     with:
       cache: true  # ← Importante
   ```
3. Considerar saltar tests de integración en CI:
   ```bash
   go test -short ./...
   ```

---

## 📚 Referencias

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Events that trigger workflows](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows)

---

**Última actualización:** 20 Nov 2025
**Versión:** 1.0
**Autor:** CI/CD Sprint 1
