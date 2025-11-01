# 🔄 Workflows de CI/CD - edugo-shared

## 🎯 Estrategia de Ejecución por Branch

Esta tabla muestra **qué workflows se ejecutan en cada tipo de branch** para evitar ejecuciones innecesarias y notificaciones de falsos positivos:

| Workflow | feature/* | dev | main | PR a dev | PR a main | Tags v* | Manual |
|----------|-----------|-----|------|----------|-----------|---------|--------|
| **ci.yml** | ❌ | ✅ (push) | ✅ (push) | ✅ | ✅ | ❌ | ❌ |
| **test.yml** | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ✅ |
| **release.yml** | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| **sync-main-to-dev.yml** | ❌ | ❌ | ✅ (push) | ❌ | ❌ | ✅ | ❌ |

### 📌 Resumen por Escenario

```bash
# Push a feature/* → SIN workflows automáticos
git push origin feature/mi-feature
# ✅ Sin ejecuciones, sin notificaciones

# Crear PR desde feature/* a dev → CI completo
gh pr create --base dev --head feature/mi-feature
# ✅ ci.yml (tests en 7 módulos)
# ✅ test.yml (cobertura por módulo)
# ✅ Copilot code review

# Merge PR a dev → Solo CI
# ✅ ci.yml se ejecuta
# ✅ Código integrado en dev

# Crear PR de dev a main → CI completo + Preparación para release
gh pr create --base main --head dev
# ✅ ci.yml (validación completa)
# ✅ test.yml (cobertura)
# ⚠️  TÚ DEBES AGREGAR TAGS MANUALMENTE en la descripción del PR

# Merge PR a main → Sin workflows (hasta que crees el tag)
# Espera a crear el tag manualmente

# Crear tag manualmente → Release automático
git tag v2.1.0
git push origin v2.1.0
# ✅ release.yml (valida 7 módulos + crea GitHub Release)
# ✅ sync-main-to-dev.yml (sincroniza main → dev)

# Release manual de módulo específico
git tag middleware/gin/v0.1.0
git push origin middleware/gin/v0.1.0
# ✅ release.yml se ejecuta igual (valida todos los módulos)
```

### ⚠️ Diferencia Clave con edugo-api-mobile

edugo-shared **NO tiene auto-versionado**. El versionado es manual para permitir:
- Control total sobre cuándo y cómo versionar
- Versionado híbrido (global + por módulo)
- Mejor coordinación con proyectos consumidores

---

## 📋 Workflows Configurados

### 1️⃣ **ci.yml** - Pipeline de Integración Continua

**Trigger:**
- ✅ Pull Requests a `main` o `dev`
- ✅ Push directo a `main` (red de seguridad)

**Ejecuta:**
- ✅ **Matrix Strategy**: Tests en paralelo de 7 módulos
  - common
  - logger
  - auth
  - middleware/gin
  - messaging/rabbit
  - database/postgres
  - database/mongodb

**Jobs por módulo:**
- ✅ Verificación de formato (gofmt)
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Build verification

**Jobs adicionales:**
- ✅ Linter (opcional, no bloquea CI)
- ✅ Compatibilidad con Go 1.23, 1.24, 1.25

**Cuándo se ejecuta:**
```bash
# Cuando creas un PR a dev o main
gh pr create --base dev --title "..."  # ← AQUÍ se ejecuta

# O cuando alguien hace push directo
git push origin main  # ← AQUÍ se ejecuta (red de seguridad)
```

**Duración estimada:** 2-3 minutos (gracias a matrix paralela)

---

### 2️⃣ **test.yml** - Tests con Cobertura

**Trigger:**
- ✅ Manual (workflow_dispatch desde GitHub UI)
- ✅ Pull Requests a `main` o `dev`

**Ejecuta:**
- ✅ **Matrix Strategy**: Cobertura en paralelo de 7 módulos
- ✅ Tests unitarios con cobertura detallada por módulo
- ✅ Upload de artifacts (retención 30 días)
- ✅ Integración con Codecov
- ✅ Resumen consolidado de cobertura

**Cuándo se ejecuta:**
```bash
# Manual desde GitHub UI:
# Actions → Tests with Coverage → Run workflow

# O automáticamente en PRs
gh pr create --base dev  # ← AQUÍ se ejecuta
```

**Meta de cobertura:**
- common: >85%
- logger: >80%
- auth: >90% (crítico para seguridad)
- middleware/gin: >85%
- messaging/rabbit: >75%
- database/*: >70%

**Duración estimada:** 2-3 minutos

---

### 3️⃣ **release.yml** - Release Automático

**Trigger:**
- ✅ Push de tags que empiecen con `v*` (ejemplo: `v2.0.5`)

**Ejecuta:**
- ✅ **Validación completa** de 7 módulos en paralelo
- ✅ Creación de **GitHub Release**
- ✅ Extracción de changelog desde CHANGELOG.md
- ✅ Instrucciones de instalación por módulo
- ✅ Upload de coverage a Codecov

**Cuándo se ejecuta:**
```bash
# Después de mergear PR de dev a main, crear tag manualmente
git checkout main
git pull
git tag v2.1.0
git push origin v2.1.0  # ← AQUÍ se ejecuta

# También funciona con tags de módulos específicos
git tag middleware/gin/v0.1.0
git push origin middleware/gin/v0.1.0  # ← AQUÍ se ejecuta
```

**Resultado:**
- ✅ GitHub Release creado con:
  - Notas de changelog
  - Instrucciones de instalación por módulo
  - Links a documentación (CHANGELOG, UPGRADE_GUIDE, README)

**Duración estimada:** 3-5 minutos

---

### 4️⃣ **sync-main-to-dev.yml** - Sincronización Automática

**Trigger:**
- ✅ Push a `main` (después de merge de PR)
- ✅ Push de tags `v*` (después de release)

**Ejecuta:**
- ✅ Detección de diferencias entre main y dev
- ✅ Creación automática de PR (main → dev)
- ✅ Intento de auto-merge si no hay conflictos

**Cuándo se ejecuta:**
```bash
# Automáticamente después de merge a main
gh pr merge  # ← Después de esto se ejecuta

# O después de crear tag de release
git push origin v2.1.0  # ← Después de esto se ejecuta
```

**Propósito:**
- Mantener `dev` sincronizado con los cambios de release en `main`
- Propagar tags de versión a `dev`
- Asegurar que el siguiente desarrollo parta del estado más reciente

**Resultado:**
- ✅ PR automático creado (si hay diferencias)
- ✅ Labels: `sync`, `automated`
- ✅ Auto-merge habilitado (si es posible)

**Duración estimada:** <1 minuto

---

## 🔄 Flujo Completo de Desarrollo

### Escenario 1: Nueva Feature

```bash
# 1. Crear feature branch desde dev
git checkout dev
git pull
git checkout -b feature/nueva-funcionalidad

# 2. Desarrollar y hacer commits
git add .
git commit -m "feat: agregar nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# 3. Crear PR a dev
gh pr create --base dev --head feature/nueva-funcionalidad
# ✅ ci.yml se ejecuta
# ✅ test.yml se ejecuta
# ✅ Copilot hace code review

# 4. Aprobar y mergear PR
gh pr merge --squash
# ✅ ci.yml se ejecuta en push a dev
```

### Escenario 2: Release de Nueva Versión

```bash
# 1. Crear PR de dev a main (cuando estés listo para release)
gh pr create --base main --head dev --title "Release v2.1.0"

# 2. En la descripción del PR, indicar tags que se crearán:
"""
## Release v2.1.0

### Cambios
- Nueva feature X
- Bugfix Y

### Tags a crear después del merge
- `v2.1.0` (tag global)
- `middleware/gin/v0.1.0` (nuevo módulo)

### Breaking Changes
- Ninguno (retrocompatible)
"""

# 3. Mergear PR a main
gh pr merge --squash

# 4. Crear tags manualmente
git checkout main
git pull
git tag v2.1.0
git push origin v2.1.0
# ✅ release.yml se ejecuta (crea GitHub Release)
# ✅ sync-main-to-dev.yml se ejecuta (sincroniza con dev)

# 5. Opcionalmente, crear tags por módulo
git tag middleware/gin/v0.1.0
git push origin middleware/gin/v0.1.0
```

### Escenario 3: Bugfix Urgente

```bash
# 1. Crear branch desde main
git checkout main
git pull
git checkout -b hotfix/bug-critico

# 2. Arreglar bug
git add .
git commit -m "fix: arreglar bug crítico"
git push origin hotfix/bug-critico

# 3. Crear PR a main (bypass dev por urgencia)
gh pr create --base main --head hotfix/bug-critico
# ✅ ci.yml se ejecuta
# ✅ test.yml se ejecuta

# 4. Mergear y crear tag patch
gh pr merge --squash
git checkout main
git pull
git tag v2.0.6
git push origin v2.0.6
# ✅ release.yml se ejecuta
# ✅ sync-main-to-dev.yml sincroniza fix con dev
```

---

## 🛠️ Comandos Make para CI Local

Antes de hacer push, valida localmente con estos comandos:

```bash
# Ejecutar CI completo localmente
make ci-all-modules

# Tests con race detection
make test-race-all-modules

# Cobertura de todos los módulos
make coverage-all-modules

# Lint de todos los módulos
make lint-all-modules

# Validación completa pre-PR
make check-all-modules

# Ver todos los comandos disponibles
make help
```

---

## 📊 Beneficios de la Matrix Strategy

edugo-shared usa **matrix strategy** para paralelizar tests:

```yaml
strategy:
  fail-fast: false
  matrix:
    module:
      - common
      - logger
      - auth
      - middleware/gin
      - messaging/rabbit
      - database/postgres
      - database/mongodb
```

**Ventajas:**
- ⚡ **Velocidad**: 7 módulos se testean en paralelo (vs secuencial)
- 🔍 **Aislamiento**: Errores de un módulo no ocultan errores de otros
- 📊 **Visibilidad**: Cada módulo tiene su propio job y reporte
- 💰 **Eficiencia**: Menor tiempo total = menos minutos de GitHub Actions

---

## 📝 Convenciones de Commits

Seguir **Conventional Commits** para facilitar generación de CHANGELOG:

```bash
# Features
feat: agregar middleware de rate limiting
feat(auth): implementar refresh tokens

# Bugfixes
fix: corregir validación de JWT expirado
fix(logger): arreglar formato de timestamps

# Breaking changes
feat!: cambiar firma de NewJWTManager
BREAKING CHANGE: ahora requiere parámetro expiration

# Otros
docs: actualizar README de módulo auth
test: agregar tests para middleware/gin
refactor: simplificar lógica de validación
chore: actualizar dependencias
```

---

## 🤖 GitHub Copilot Integration

Todos los workflows están configurados para trabajar con GitHub Copilot:

- ✅ Code reviews automáticos en español
- ✅ Sugerencias contextuales según el módulo
- ✅ Detección de anti-patrones en librerías
- ✅ Validación de breaking changes

Configuración en: `.github/copilot-instructions.md`

---

## ⚠️ Notas Importantes

### Versionado Manual

edugo-shared **requiere versionado manual** porque:
1. Es una librería, no una aplicación
2. Usa versionado híbrido (global + por módulo)
3. Breaking changes requieren coordinación con consumidores

### Semantic Versioning

Seguir estrictamente **SemVer 2.0.0**:

```
vMAJOR.MINOR.PATCH

MAJOR: Breaking changes
MINOR: Nuevas features (retrocompatibles)
PATCH: Bugfixes (retrocompatibles)
```

### Consumidores de edugo-shared

Estos proyectos dependen de edugo-shared:
- edugo-api-mobile
- edugo-api-administracion
- edugo-worker

Cualquier release nuevo debe comunicarse a estos equipos.

---

## 📚 Referencias

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Go Modules](https://go.dev/ref/mod)

---

**Última actualización**: 2025-11-01
**Versión de la librería**: v2.0.5
**Go Version**: 1.25.3
