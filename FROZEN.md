# 🔒 REPOSITORIO CONGELADO

**Fecha de congelamiento:** 2025-11-15
**Versión congelada:** v0.7.0
**Status:** 🔒 FROZEN - NO NEW FEATURES

---

## ⚠️ Política de Congelamiento

Este repositorio está **CONGELADO** para nuevas features hasta después del MVP de EduGo.

### ✅ Permitido

#### 🐛 Bug Fixes Críticos
- **Versiones:** v0.7.1, v0.7.2, v0.7.3, etc. (PATCH bumps solamente)
- **Criterio:** Solo bugs que bloquean producción o security issues
- **Proceso:** Requiere aprobación explícita del tech lead
- **Branch:** `hotfix/v0.7.x-descripcion`

#### 📝 Documentación
- Mejoras de README
- Guías de uso
- Comentarios en código
- **No afecta:** Versiones de módulos

### ❌ NO Permitido

- ✨ **Nuevas features** - Cualquier funcionalidad nueva
- 🔄 **Refactoring** - Cambios estructurales sin bug fix
- ⬆️ **Dependency upgrades** - Excepto security patches críticos
- 🏗️ **Breaking changes** - Cambios incompatibles en APIs públicas
- 🧪 **Experimental features** - POCs o features en desarrollo

---

## 📦 Versión Congelada: v0.7.0

### Módulos Incluidos (12)

Todos los módulos están en versión **v0.7.0**:

| Módulo | Versión | Coverage | Descripción |
|--------|---------|----------|-------------|
| auth | v0.7.0 | 87.3% | JWT Authentication |
| logger | v0.7.0 | 95.8% | Logging con Zap |
| common | v0.7.0 | >94% | Errors, Types, Validator |
| config | v0.7.0 | 82.9% | Configuration loader |
| bootstrap | v0.7.0 | 31.9% | Dependency injection |
| lifecycle | v0.7.0 | 91.8% | Application lifecycle |
| middleware/gin | v0.7.0 | 98.5% | Gin middleware |
| messaging/rabbit | v0.7.0 | 3.2% | RabbitMQ + DLQ |
| database/postgres | v0.7.0 | 58.8% | PostgreSQL utilities |
| database/mongodb | v0.7.0 | 54.5% | MongoDB utilities |
| testing | v0.7.0 | 59.0% | Testing utilities |
| evaluation | v0.7.0 | 100% | Assessment models |

### Características Clave

- ✅ Go version: **1.24.10** (todos los módulos)
- ✅ Coverage global: ~**75%**
- ✅ Dead Letter Queue (DLQ) en messaging/rabbit
- ✅ Módulo evaluation completo (100% coverage)
- ✅ Tests comprehensivos en módulos core

### Instalación

```bash
# Instalar módulo específico
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/logger@v0.7.0
# ... (resto de módulos)

# O actualizar todos
go get github.com/EduGoGroup/edugo-shared/...@v0.7.0
go mod tidy
```

---

## 🚀 Proyectos Consumidores

Los siguientes proyectos deben usar **exclusivamente v0.7.0**:

- **edugo-api-mobile** - API móvil de estudiantes
- **edugo-api-administracion** - API de administración
- **edugo-worker** - Background worker

⚠️ **NO actualizar** a versiones posteriores sin aprobación del equipo.

---

## 📋 Proceso de Bug Fix

Si encuentras un bug crítico:

### 1. Verificar Criticidad
- ¿Bloquea producción?
- ¿Es un security issue?
- ¿Afecta funcionalidad core?

Si NO es crítico → **Esperar a descongelamiento**

### 2. Crear Issue
- Título: `[CRITICAL] Descripción del bug`
- Labels: `bug`, `critical`, `v0.7.x`
- Describir impacto y reproducción

### 3. Obtener Aprobación
- Tag al tech lead en el issue
- Esperar aprobación explícita

### 4. Crear Hotfix Branch
```bash
git checkout main
git pull origin main
git checkout -b hotfix/v0.7.x-nombre-descriptivo
```

### 5. Implementar Fix Mínimo
- ✅ Solo el fix del bug
- ❌ No refactoring
- ❌ No mejoras adicionales
- ✅ Tests que reproduzcan el bug + fix

### 6. Crear PR
- Base: `main`
- Labels: `hotfix`, `critical`
- Reviewers: Tech lead + 1

### 7. Después del Merge
```bash
# Bump PATCH version
git tag auth/v0.7.1  # (o el módulo afectado)
git push origin auth/v0.7.1

# Merge main → dev
git checkout dev
git merge main
git push origin dev
```

### 8. Update CHANGELOG
```markdown
## [0.7.1] - YYYY-MM-DD

### Fixed
- auth: Fixed critical bug in token validation
```

---

## 🔓 Descongelamiento

El repositorio se descongelará después de:

1. ✅ MVP lanzado a producción
2. ✅ Período de estabilización (2-4 semanas)
3. ✅ 0 bugs críticos pendientes
4. ✅ Decisión explícita del equipo

**Próxima versión después de descongelar:** v0.8.0 (MINOR bump para nuevas features)

---

## 📊 Estadísticas del Release

**Fecha de congelamiento:** 2025-11-15
**Commits desde v0.6.x:** ~15 commits
**Sprints completados:** 3 (Sprint 0, 1, 2, 3)
**Tests agregados:** ~100+ nuevos tests
**Coverage mejorado:** +15 puntos (~60% → ~75%)
**Módulos nuevos:** 1 (evaluation)
**Features nuevas:** 1 (DLQ)

---

## 🏆 Equipo

**Mantenedores:** Equipo EduGo
**Tech Lead:** Jhoan Medina
**Asistencia técnica:** Claude Code

---

## 📞 Soporte

Para issues críticos:
- **GitHub Issues:** https://github.com/EduGoGroup/edugo-shared/issues
- **Tag:** @medinatello (tech lead)
- **Label:** `critical`, `v0.7.x`

---

**Última actualización:** 2025-11-15
**Próxima revisión:** Post-MVP (TBD)

🔒 **FROZEN UNTIL POST-MVP** 🔒
