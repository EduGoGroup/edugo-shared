# Estrategia de Versionamiento - EduGo Shared

## 📦 Arquitectura Modular

Este repositorio es un **monorepo multi-módulo** con versionamiento independiente para cada módulo de Go.

### Estructura del Proyecto

```
edugo-shared/
├── auth/              → auth/vX.Y.Z
├── bootstrap/         → bootstrap/vX.Y.Z
├── common/            → common/vX.Y.Z
├── config/            → config/vX.Y.Z
├── database/
│   ├── mongodb/       → database/mongodb/vX.Y.Z
│   └── postgres/      → database/postgres/vX.Y.Z
├── evaluation/        → evaluation/vX.Y.Z
├── lifecycle/         → lifecycle/vX.Y.Z
├── logger/            → logger/vX.Y.Z
├── messaging/
│   └── rabbit/        → messaging/rabbit/vX.Y.Z
├── middleware/
│   └── gin/           → middleware/gin/vX.Y.Z
└── testing/           → testing/vX.Y.Z
```

Cada directorio con `go.mod` es un **módulo independiente** con su propio ciclo de versionamiento.

---

## 🏷️ Estrategia de Tags

### ✅ Correcto: Tags por Módulo

Cada módulo tiene su propia secuencia de versiones:

```bash
# Módulo Auth
auth/v0.3.0
auth/v0.4.0
auth/v0.5.0
auth/v0.7.0

# Módulo Bootstrap
bootstrap/v0.1.0
bootstrap/v0.4.0
bootstrap/v0.7.0
bootstrap/v0.8.0
bootstrap/v0.9.0  ← Versión más reciente

# Módulo Evaluation
evaluation/v0.7.0
evaluation/v0.8.0  ← Versión más reciente
```

### ❌ Incorrecto: Tags Globales

**NO SE UTILIZAN** tags globales del tipo `v0.X.Y` sin prefijo de módulo.

**Razón:** Este repositorio NO es un módulo único, sino una colección de módulos independientes. Un tag global no tiene significado en este contexto.

---

## 📋 Versionamiento por Módulo

### Principios

1. **Independencia Total**
   - Cada módulo evoluciona a su propio ritmo
   - Un cambio en `auth` no requiere nueva versión de `logger`
   - Las versiones pueden estar desincronizadas entre módulos

2. **Semantic Versioning**
   - **MAJOR** (`vX.0.0`): Breaking changes incompatibles
   - **MINOR** (`v0.X.0`): Nueva funcionalidad compatible
   - **PATCH** (`v0.0.X`): Corrección de bugs

3. **Formato de Tags**
   ```
   <module-path>/v<MAJOR>.<MINOR>.<PATCH>
   ```
   
   Ejemplos:
   ```
   auth/v1.2.3
   database/mongodb/v2.0.1
   middleware/gin/v0.5.0
   ```

---

## 🚀 Proceso de Release

### Release de un Solo Módulo

Cuando se modifica solo un módulo:

```bash
# 1. Realizar cambios en el módulo
cd bootstrap/
# ... hacer modificaciones ...

# 2. Actualizar tests
go test ./...

# 3. Commit de cambios
git add .
git commit -m "feat(bootstrap): nueva funcionalidad X"

# 4. Crear tag del módulo
git tag bootstrap/v0.10.0

# 5. Push de cambios y tag
git push origin dev
git push origin bootstrap/v0.10.0
```

**Resultado:** Solo `bootstrap` recibe nuevo tag. Los demás módulos no se afectan.

---

### Release Coordinado (Múltiples Módulos)

Cuando un cambio afecta múltiples módulos (ej: breaking change en `common` usado por varios módulos):

```bash
# 1. Realizar cambios en módulos afectados
# common/, auth/, logger/, etc.

# 2. Commit coordinado
git commit -m "feat!: actualizar interface Logger (BREAKING CHANGE)"

# 3. Crear tags para TODOS los módulos afectados
git tag common/v1.0.0       # Breaking change → MAJOR bump
git tag auth/v0.8.0         # Adaptación → MINOR bump
git tag logger/v0.8.0       # Adaptación → MINOR bump
git tag bootstrap/v0.10.0   # Adaptación → MINOR bump

# 4. Push todos los tags
git push origin --tags
```

**Criterio:** Solo crear tags para módulos que **realmente cambiaron**. No sincronizar versiones artificialmente.

---

## 🔍 Verificar Versiones Actuales

### Ver todas las versiones de un módulo

```bash
git tag -l 'auth/v*'
# Output:
# auth/v0.3.0
# auth/v0.4.0
# auth/v0.5.0
# auth/v0.7.0
```

### Ver versión más reciente de cada módulo

```bash
for module in auth bootstrap common config logger lifecycle testing; do
  echo "$module: $(git tag -l "$module/v*" | sort -V | tail -1)"
done

for module in database/mongodb database/postgres; do
  echo "$module: $(git tag -l "$module/v*" | sort -V | tail -1)"
done

for module in messaging/rabbit middleware/gin; do
  echo "$module: $(git tag -l "$module/v*" | sort -V | tail -1)"
done
```

### Estado actual (Enero 2025)

```
auth:              auth/v0.7.0
bootstrap:         bootstrap/v0.9.0  ← Más avanzado
common:            common/v0.7.0
config:            config/v0.7.0
database/mongodb:  database/mongodb/v0.7.0
database/postgres: database/postgres/v0.7.0
evaluation:        evaluation/v0.8.0  ← Más avanzado
lifecycle:         lifecycle/v0.7.0
logger:            logger/v0.7.0
messaging/rabbit:  messaging/rabbit/v0.7.0
middleware/gin:    middleware/gin/v0.7.0
testing:           testing/v0.7.0
```

**Esto es normal y esperado** en un monorepo modular.

---

## 📥 Consumo de Módulos (Go Modules)

### Instalar un Módulo Específico

```bash
# Versión específica
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0

# Última versión (recomendado)
go get github.com/EduGoGroup/edugo-shared/auth@latest

# Versión específica de otro módulo
go get github.com/EduGoGroup/edugo-shared/bootstrap@v0.9.0
```

### Actualizar un Módulo

```bash
# Actualizar solo auth
go get -u github.com/EduGoGroup/edugo-shared/auth

# Actualizar múltiples módulos
go get -u github.com/EduGoGroup/edugo-shared/auth
go get -u github.com/EduGoGroup/edugo-shared/logger
```

---

## 🎯 Casos de Uso Comunes

### 1. Bug Fix en un Solo Módulo

```bash
# Ejemplo: Fix en logger
cd logger/
# ... fix bug ...
git commit -m "fix(logger): corregir memory leak en JSON formatter"
git tag logger/v0.7.1  # PATCH bump
git push origin dev
git push origin logger/v0.7.1
```

**Otros módulos:** Sin cambios, mantienen sus versiones.

---

### 2. Nueva Feature en un Módulo

```bash
# Ejemplo: Nueva feature en auth
cd auth/
# ... implementar refresh tokens ...
git commit -m "feat(auth): agregar soporte para refresh tokens"
git tag auth/v0.8.0  # MINOR bump
git push origin dev
git push origin auth/v0.8.0
```

---

### 3. Breaking Change en Módulo Base

Cuando `common` cambia y rompe compatibilidad:

```bash
# 1. Actualizar common
cd common/
# ... breaking change en errors package ...
git commit -m "feat(common)!: rediseñar error handling (BREAKING)"

# 2. Actualizar módulos dependientes
cd ../auth && # adaptar código
cd ../logger && # adaptar código
cd ../bootstrap && # adaptar código

# 3. Tags coordinados
git tag common/v1.0.0       # MAJOR (breaking)
git tag auth/v0.8.0         # MINOR (adaptación)
git tag logger/v0.8.0       # MINOR (adaptación)
git tag bootstrap/v0.10.0   # MINOR (adaptación)

# 4. Push
git push origin --tags
```

---

### 4. Sincronización Voluntaria (Freeze de versión base)

En momentos clave (ej: MVP freeze), sincronizar todas las versiones:

```bash
# Commit: "Release v0.7.0 - FROZEN base for EduGo MVP"
git tag auth/v0.7.0
git tag bootstrap/v0.7.0
git tag common/v0.7.0
git tag config/v0.7.0
# ... etc para todos los módulos ...

git push origin --tags
```

**Nota:** Esto es una **decisión voluntaria** de sincronización, no un requisito del sistema de versionamiento.

---

## 📚 Referencias

- [Go Modules: Multi-Module Repositories](https://go.dev/wiki/Modules#faqs--multi-module-repositories)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## ❓ FAQ

### ¿Por qué no hay un tag `v0.X.Y` global?

Porque este repositorio **NO es un módulo único**. Es una colección de módulos independientes. Go Modules no reconoce tags globales en monorepos multi-módulo.

### ¿Puedo tener `auth/v0.5.0` y `logger/v0.9.0` al mismo tiempo?

**Sí, absolutamente.** Cada módulo evoluciona independientemente. Las versiones pueden estar desincronizadas.

### ¿Cuándo debo sincronizar versiones?

Solo cuando:
1. Un breaking change en módulo base afecta a varios
2. Milestone importante del proyecto (MVP, v1.0, etc.)
3. Decisión estratégica del equipo

**No es obligatorio** sincronizar versiones.

### ¿Cómo sé qué versión instalar?

```bash
# Ver última versión de un módulo en GitHub
git ls-remote --tags https://github.com/EduGoGroup/edugo-shared.git | grep auth/

# O simplemente usar @latest
go get github.com/EduGoGroup/edugo-shared/auth@latest
```

### ¿Puedo crear un tag global si quiero?

**No es recomendable.** Go Modules ignora tags globales en monorepos. Usa tags por módulo o etiquetas de release en GitHub.

---

## 🔗 Ver También

- [CHANGELOG.md](CHANGELOG.md) - Historial detallado de cambios
- [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md) - Guía de migración entre versiones
- [README.md](README.md) - Documentación general del proyecto
