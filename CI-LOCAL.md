# 🏠 Guía de CI/CD Local

Esta guía te ayuda a **probar localmente** antes de hacer push, ahorrando tiempo y evitando fallos en GitHub Actions.

---

## 🎯 ¿Cuándo se Ejecuta el CI/CD en GitHub?

Los workflows se ejecutan automáticamente en:

| Workflow | Trigger | Frecuencia |
|----------|---------|------------|
| **ci.yml** | Push a `main`, `develop` | ❌ CADA PUSH |
| **test.yml** | Push a `main`, PR | ❌ CADA PUSH |
| **simple-test.yml** | Push a `main` | ❌ CADA PUSH |

**Problema:** Cada push ejecuta 3 workflows = gasto de minutos de GitHub Actions

**Solución:** Probar localmente ANTES de hacer push

---

## 🚀 Scripts Disponibles

### 1️⃣ **Verificación Completa** (simula GitHub Actions)

```bash
./test-ci-local.sh
```

**Ejecuta:**
- ✅ Descargar dependencias
- ✅ Formatear código
- ✅ Análisis estático (go vet)
- ✅ Linter (golangci-lint)
- ✅ Tests unitarios
- ✅ Tests con race detection
- ✅ Cobertura de código
- ✅ Verificar build

**Duración:** ~2-3 minutos

**Cuándo usarlo:** Antes de hacer push a `main`

---

### 2️⃣ **Verificación Rápida** (pre-commit)

```bash
./test-quick.sh
```

**Ejecuta:**
- ✅ Formatear código
- ✅ Análisis estático
- ✅ Tests rápidos
- ✅ Build

**Duración:** ~30 segundos

**Cuándo usarlo:** Antes de cada commit

---

### 3️⃣ **Comandos Make Individuales**

Puedes ejecutar pasos individuales:

```bash
# Ver todos los comandos
make help

# Verificación rápida pre-commit
make pre-commit

# Pipeline CI completo (como GitHub)
make ci

# Solo tests
make test

# Solo tests con cobertura
make test-coverage

# Solo linter
make lint

# Solo formatear
make fmt
```

---

## 🔥 Workflow Recomendado

### **Antes de cada commit:**
```bash
./test-quick.sh
git add .
git commit -m "tu mensaje"
```

### **Antes de hacer push a main:**
```bash
./test-ci-local.sh
git push origin main
```

### **Antes de crear un tag de versión:**
```bash
./test-ci-local.sh
git tag -a v2.0.0 -m "mensaje"
git push origin main
git push origin v2.0.0
```

---

## 📊 Comparación de Opciones

| Método | Tiempo | Cobertura | Cuándo usar |
|--------|--------|-----------|-------------|
| `./test-quick.sh` | 30s | Básica | Antes de cada commit |
| `./test-ci-local.sh` | 2-3min | Completa | Antes de push a main |
| `make ci` | 2-3min | Completa | Equivalente a CI local |
| GitHub Actions | 3-5min | Completa | Automático (en cada push) |

---

## 🎯 Ventajas de Probar Localmente

✅ **Ahorras tiempo** - No esperas 5 minutos por GitHub
✅ **Ahorras minutos de GitHub Actions** - No gastas cuota
✅ **Iteras más rápido** - Fix inmediato de errores
✅ **Commits más limpios** - Solo subes código que ya pasó CI
✅ **Menos spam en PRs** - No múltiples commits de "fix CI"

---

## 🛠️ Instalación de Herramientas (una sola vez)

Si no tienes las herramientas instaladas:

```bash
make install-tools
```

Esto instala:
- `golangci-lint` - Linter de Go
- `gosec` - Security scanner
- Otras herramientas de desarrollo

---

## ❓ FAQ

### **¿Necesito instalar algo especial?**
No, solo Go 1.23+ y las herramientas con `make install-tools`

### **¿Puedo saltarme el linter?**
Sí, el linter no es crítico. Si falla, puedes hacer push igual.

### **¿Qué hago si `./test-ci-local.sh` falla?**
1. Lee el error
2. Corrige el problema
3. Vuelve a ejecutar
4. No hagas push hasta que pase

### **¿Puedo usar esto en otros proyectos?**
¡Sí! Copia los scripts y el Makefile a tus otros proyectos Go.

---

## 🎉 Ejemplo de Uso Completo

```bash
# 1. Hacer cambios en el código
vim pkg/auth/jwt.go

# 2. Verificación rápida
./test-quick.sh

# 3. Commit
git add .
git commit -m "feat: agregar nueva funcionalidad JWT"

# 4. Antes de push, verificación completa
./test-ci-local.sh

# 5. Si todo pasó, hacer push
git push origin main

# 6. (Opcional) GitHub Actions confirmará que todo está OK
```

---

**¿Resultado?** Solo haces push cuando **sabes que va a pasar** el CI de GitHub. 🎯
