# Sprint 3 - Día 2: Validación Completa

## Objetivo

Validar que:
1. ✅ Todos los tests pasan (0 failing)
2. ✅ Coverage global >85%
3. ✅ api-mobile, api-admin, worker compilan con las nuevas versiones

## Requisitos

- Go 1.24.10
- Docker instalado
- Acceso a los repos: api-mobile, api-admin, worker

---

## Paso 1: Ejecutar Suite Completa de Tests

### 1.1 Tests por módulo

Ejecuta tests en cada módulo para verificar que todos pasen:

```bash
cd /home/user/edugo-shared

# auth
cd auth && go test -v -cover ./... && cd ..

# logger
cd logger && go test -v -cover ./... && cd ..

# common
cd common && go test -v -cover ./... && cd ..

# config
cd config && go test -v -cover ./... && cd ..

# bootstrap
cd bootstrap && go test -v -cover ./... && cd ..

# lifecycle
cd lifecycle && go test -v -cover ./... && cd ..

# middleware/gin
cd middleware/gin && go test -v -cover ./... && cd ..

# messaging/rabbit
cd messaging/rabbit && go test -v -cover ./... && cd ..

# database/postgres
cd database/postgres && go test -v -cover ./... && cd ..

# database/mongodb
cd database/mongodb && go test -v -cover ./... && cd ..

# testing
cd testing && go test -v -cover ./... && cd ..

# evaluation
cd evaluation && go test -v -cover ./... && cd ..
```

### 1.2 Verificar resultados

**Todos los módulos deben mostrar:**
```
PASS
coverage: XX.X% of statements
ok      github.com/EduGoGroup/edugo-shared/[module]
```

**SI HAY FAILURES:**
- Anota el módulo y el test que falla
- Revisa el error
- Corrige antes de continuar

---

## Paso 2: Calcular Coverage Global

### 2.1 Coverage por módulo

Genera reportes de coverage para cada módulo:

```bash
#!/bin/bash
cd /home/user/edugo-shared

# Crear directorio para reportes
mkdir -p coverage-reports

# Generar coverage para cada módulo
modules=(
  "auth"
  "logger"
  "common"
  "config"
  "bootstrap"
  "lifecycle"
  "middleware/gin"
  "messaging/rabbit"
  "database/postgres"
  "database/mongodb"
  "testing"
  "evaluation"
)

for module in "${modules[@]}"; do
  echo "=== Coverage para $module ==="
  cd "$module"
  go test -coverprofile=../coverage-reports/${module//\//-}.out ./...
  coverage=$(go tool cover -func=../coverage-reports/${module//\//-}.out | grep total | awk '{print $3}')
  echo "$module: $coverage"
  cd - > /dev/null
done
```

### 2.2 Resultados esperados por módulo

| Módulo | Coverage Target | Status |
|--------|----------------|--------|
| auth | >80% | ✅ Ya tiene tests |
| logger | >80% | ✅ Sprint 2 |
| common/errors | >80% | ✅ Sprint 2 |
| common/validator | >80% | ✅ Sprint 2 |
| common/types | >80% | ✅ Sprint 2 |
| config | >80% | ✅ Sprint 3 Day 1 |
| bootstrap | >80% | ✅ Sprint 3 Day 1 |
| lifecycle | ? | ⚠️ Verificar |
| middleware/gin | >80% | ✅ Ya tiene tests |
| messaging/rabbit | >70% | ✅ Sprint 1 |
| database/postgres | >80% | ✅ Sprint 1 |
| database/mongodb | >80% | ⚠️ Verificar |
| testing | >80% | ✅ Ya tiene tests |
| evaluation | 100% | ✅ Sprint 1 |

### 2.3 Calcular coverage global

```bash
# Combinar todos los reportes
cd /home/user/edugo-shared/coverage-reports
cat *.out > combined.out

# Calcular coverage total
go tool cover -func=combined.out | grep total

# Debe mostrar: total: (statements) >85%
```

**Objetivo: >85% coverage global**

---

## Paso 3: Validar Compilación de Proyectos Consumidores

### 3.1 api-mobile

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Actualizar a las últimas versiones de edugo-shared
go get github.com/EduGoGroup/edugo-shared/auth@latest
go get github.com/EduGoGroup/edugo-shared/logger@latest
go get github.com/EduGoGroup/edugo-shared/common@latest
go get github.com/EduGoGroup/edugo-shared/config@latest
go get github.com/EduGoGroup/edugo-shared/bootstrap@latest
go get github.com/EduGoGroup/edugo-shared/lifecycle@latest
go get github.com/EduGoGroup/edugo-shared/middleware/gin@latest
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@latest
go get github.com/EduGoGroup/edugo-shared/database/postgres@latest
go get github.com/EduGoGroup/edugo-shared/database/mongodb@latest
go get github.com/EduGoGroup/edugo-shared/testing@latest
go get github.com/EduGoGroup/edugo-shared/evaluation@latest

# Tidy dependencies
go mod tidy

# Intentar compilar
go build ./cmd/api-mobile

# Debe compilar sin errores
echo $?  # Debe ser 0
```

**Resultado esperado:**
```
✅ Compilación exitosa (exit code 0)
```

**Si hay errores:**
- Anota el error completo
- Verifica si es breaking change en edugo-shared
- Ajusta código de api-mobile si es necesario

### 3.2 api-admin

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Actualizar dependencias (igual que api-mobile)
go get github.com/EduGoGroup/edugo-shared/auth@latest
go get github.com/EduGoGroup/edugo-shared/logger@latest
# ... (resto de módulos)

go mod tidy
go build ./cmd/api-admin

echo $?  # Debe ser 0
```

### 3.3 worker

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker

# Actualizar dependencias
go get github.com/EduGoGroup/edugo-shared/evaluation@latest
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@latest
# ... (otros que use worker)

go mod tidy
go build ./cmd/worker

echo $?  # Debe ser 0
```

---

## Checklist de Validación

Marca cada item cuando esté completado:

### Tests
- [ ] auth: tests PASS
- [ ] logger: tests PASS
- [ ] common: tests PASS
- [ ] config: tests PASS
- [ ] bootstrap: tests PASS
- [ ] lifecycle: tests PASS
- [ ] middleware/gin: tests PASS
- [ ] messaging/rabbit: tests PASS
- [ ] database/postgres: tests PASS
- [ ] database/mongodb: tests PASS
- [ ] testing: tests PASS
- [ ] evaluation: tests PASS

### Coverage
- [ ] Coverage global calculado
- [ ] Coverage global >85%
- [ ] Todos los módulos P0 >80%

### Compilación
- [ ] api-mobile compila sin errores
- [ ] api-admin compila sin errores
- [ ] worker compila sin errores

---

## Criterios de Éxito Día 2

Para pasar al Día 3 (Release v0.7.0), **TODOS** estos deben cumplirse:

✅ 0 tests failing
✅ Coverage global >85%
✅ Todos los consumidores compilan
✅ No hay breaking changes no documentados

---

## Si Algo Falla

### Tests fallan

1. Identifica el módulo y test específico
2. Revisa el error
3. Corrige el código o test
4. Re-ejecuta: `go test -v ./...`
5. Commit fix: `git commit -m "fix(module): description"`
6. Push: `git push`

### Coverage <85%

1. Identifica módulos con coverage bajo
2. Agrega tests adicionales
3. Re-calcula coverage
4. Commit tests: `git commit -m "test(module): increase coverage"`

### Consumidores no compilan

1. Identifica el error de compilación
2. Determina si es breaking change
3. **Opción A:** Ajusta código en edugo-shared para mantener compatibilidad
4. **Opción B:** Documenta breaking change y ajusta consumidores
5. Re-compila para verificar

---

## Próximos Pasos

Una vez que TODOS los criterios de éxito se cumplan:

➡️ **Continuar a Sprint 3 - Día 3: Release v0.7.0**

El Día 3 incluirá:
- Crear rama release/v0.7.0
- Actualizar CHANGELOG.md
- Mergear a main
- Crear todos los tags (12 módulos)
- Push de tags
- Crear GitHub Release
- Mergear main → dev
- **Congelar edugo-shared** 🔒
