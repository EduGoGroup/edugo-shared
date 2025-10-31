# 🔄 Workflows de CI/CD - edugo-shared

## 📋 Workflows Configurados

### 1️⃣ **release.yml** - Release Completo (TAGS)

**Trigger:** Solo cuando creas un tag `v*` (ej: `v2.0.0`)

**Ejecuta:**
- ✅ Verificación de formato
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Cobertura de código
- ✅ Build verification
- ✅ Creación automática de GitHub Release

**Cuándo se ejecuta:**
```bash
git tag -a v2.0.0 -m "Release 2.0.0"
git push origin v2.0.0  # ← AQUÍ se ejecuta
```

**Duración estimada:** 3-5 minutos

---

### 2️⃣ **ci.yml** - Pipeline de Integración Continua

**Trigger:**
- ✅ Pull Requests a `main` o `develop`
- ✅ Push directo a `main` (red de seguridad)

**Ejecuta:**
- ✅ Verificación de formato
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Build verification
- ✅ Linter (opcional, no bloquea)
- ✅ Compatibilidad con Go 1.23, 1.24, 1.25

**Cuándo se ejecuta:**
```bash
# Cuando creas un PR
gh pr create --title "..." --body "..."  # ← AQUÍ se ejecuta

# O cuando alguien hace push directo a main (no recomendado)
git push origin main  # ← AQUÍ se ejecuta
```

**Duración estimada:** 2-3 minutos

---

### 3️⃣ **test.yml** - Tests con Cobertura (MANUAL/PR)

**Trigger:**
- ✅ Manual (workflow_dispatch desde GitHub UI)
- ✅ Pull Requests a `main` o `develop`

**Ejecuta:**
- ✅ Tests del módulo core con cobertura
- ✅ Tests del módulo PostgreSQL con cobertura
- ✅ Tests del módulo MongoDB con cobertura
- ✅ Upload de reportes a Codecov
- ✅ Artifacts con reportes de cobertura

**Cuándo se ejecuta:**
```bash
# Manual desde GitHub UI:
# Actions → Tests with Coverage → Run workflow

# O automáticamente en PRs
gh pr create  # ← AQUÍ se ejecuta
```

**Duración estimada:** 2-3 minutos

---

## 🎯 Estrategia de CI/CD Optimizada

### **Flujo Normal de Desarrollo:**

```
┌─────────────────────────────────────────────────────────────┐
│  1. Desarrollo Local                                        │
│     - Hacer cambios                                         │
│     - ./test-quick.sh (30s)                                 │
│     - git commit                                            │
│     ✅ NO GASTA MINUTOS DE GITHUB                           │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Crear Pull Request                                      │
│     - gh pr create                                          │
│     - CI automático (ci.yml + test.yml)                     │
│     - Revisar resultados                                    │
│     ✅ VALIDA ANTES DE MERGE                                │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Merge a Main                                            │
│     - gh pr merge                                           │
│     - CI de seguridad (ci.yml) si se hace push directo     │
│     ✅ CÓDIGO VALIDADO EN MAIN                              │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Crear Release                                           │
│     - ./test-ci-local.sh (validación final local)          │
│     - git tag -a v2.0.0 -m "..."                            │
│     - git push origin v2.0.0                                │
│     - Release automático (release.yml)                      │
│     ✅ RELEASE CON VALIDACIÓN COMPLETA                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 💰 Ahorro de Minutos de GitHub Actions

### **Antes (cada push ejecutaba 3 workflows):**
```
Push a main → 3 workflows × 5 min = 15 minutos
10 pushes al día = 150 minutos/día
Mes = 4,500 minutos (¡casi 100% del plan gratuito!)
```

### **Después (optimizado):**
```
Push a main → 1 workflow × 3 min = 3 minutos
PR → 2 workflows × 5 min = 10 minutos
Tag → 1 workflow × 5 min = 5 minutos

Mes típico:
- 5 PRs = 50 minutos
- 2 releases = 10 minutos
- 10 pushes = 30 minutos
Total = 90 minutos/mes (✅ Solo 3-4% del plan gratuito)
```

**Ahorro:** ~95% de minutos 🎉

---

## 🚀 Guía Rápida

### **Para desarrollo normal:**
```bash
# 1. Desarrollar localmente
vim pkg/auth/jwt.go

# 2. Probar localmente (NO usa GitHub)
./test-quick.sh

# 3. Commit
git commit -m "feat: nueva funcionalidad"

# 4. Push a tu rama
git push origin feature/nueva-funcionalidad

# 5. Crear PR (ejecuta CI automáticamente)
gh pr create --title "Nueva funcionalidad" --body "..."

# 6. Esperar aprobación y merge
```

### **Para crear una release:**
```bash
# 1. Probar todo localmente
./test-ci-local.sh

# 2. Actualizar CHANGELOG.md y versión en README.md

# 3. Commit y push
git add .
git commit -m "chore: preparar release v2.1.0"
git push origin main

# 4. Crear y push tag (ejecuta release.yml)
git tag -a v2.1.0 -m "Release v2.1.0"
git push origin v2.1.0

# 5. GitHub Actions creará la release automáticamente
```

---

## 📊 Comparación de Configuraciones

| Escenario | Antes | Después | Ahorro |
|-----------|-------|---------|--------|
| Push a main | 15 min (3 workflows) | 3 min (1 workflow) | 80% |
| Pull Request | No configurado | 10 min (2 workflows) | ✅ Nueva feature |
| Release/Tag | No automatizado | 5 min (1 workflow) | ✅ Nueva feature |
| Minutos mensuales | ~4,500 | ~90 | 98% |

---

## 🛡️ Branch Protection (Recomendado)

Para forzar el uso de PRs, configura protección de rama:

1. GitHub → Settings → Branches → Add rule
2. Branch name pattern: `main`
3. Configurar:
   - ✅ Require pull request before merging
   - ✅ Require status checks to pass before merging
   - ✅ Status checks: "Tests and Checks", "Tests with Coverage"
   - ✅ Require branches to be up to date before merging

Esto previene push directo a `main` y garantiza que todo pase por PR + CI.

---

## 🔍 Ver Estado de Workflows

```bash
# Ver últimos workflows ejecutados
gh run list --limit 10

# Ver detalles de un workflow específico
gh run view <run-id>

# Ver logs de un workflow
gh run view <run-id> --log

# Re-ejecutar un workflow fallido
gh run rerun <run-id>
```

---

## 📚 Recursos

- [Testing Local (CI-LOCAL.md)](../../CI-LOCAL.md)
- [Changelog](../../CHANGELOG.md)
- [Upgrade Guide](../../UPGRADE_GUIDE.md)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

---

**Última actualización:** 2025-10-31
**Mantenedor:** Equipo EduGo
