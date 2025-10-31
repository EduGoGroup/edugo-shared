#!/bin/bash

# Script para simular CI/CD localmente antes de hacer push
# Esto te permite probar todo sin gastar tiempo en GitHub Actions

set -e  # Salir si cualquier comando falla

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo ""
echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                       ║${NC}"
echo -e "${BLUE}║      🔍 SIMULACIÓN DE CI/CD LOCAL - edugo-shared     ║${NC}"
echo -e "${BLUE}║                                                       ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Función para imprimir paso
print_step() {
    echo ""
    echo -e "${BLUE}▶ $1${NC}"
    echo "─────────────────────────────────────────────────────"
}

# Función para imprimir éxito
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Función para imprimir error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Función para imprimir advertencia
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# 1. Descargar dependencias
print_step "1/7 Descargando dependencias..."
make deps
print_success "Dependencias descargadas"

# 2. Verificar formato
print_step "2/7 Verificando formato del código..."
make fmt
if git diff --exit-code > /dev/null 2>&1; then
    print_success "Código formateado correctamente"
else
    print_error "Código no está formateado. Se aplicaron cambios automáticamente."
    echo "Por favor revisa los cambios con 'git diff'"
    exit 1
fi

# 3. Análisis estático (go vet)
print_step "3/7 Ejecutando análisis estático (go vet)..."
make vet
print_success "Análisis estático pasó"

# 4. Linter (no crítico, solo warning)
print_step "4/7 Ejecutando linter (golangci-lint)..."
if make lint 2>/dev/null; then
    print_success "Linter pasó sin problemas"
else
    print_warning "Linter encontró advertencias (no crítico, puedes continuar)"
fi

# 5. Tests básicos
print_step "5/7 Ejecutando tests básicos..."
make test
print_success "Tests básicos pasaron"

# 6. Tests con race detection
print_step "6/7 Ejecutando tests con race detection..."
make test-race
print_success "Tests con race detection pasaron"

# 7. Cobertura de tests
print_step "7/7 Generando reporte de cobertura..."
make test-coverage
print_success "Cobertura generada"

# 8. Build verification
print_step "BONUS: Verificando que compila correctamente..."
make build
print_success "Proyecto compila correctamente"

# Resumen final
echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                       ║${NC}"
echo -e "${GREEN}║           ✓ CI/CD LOCAL COMPLETADO CON ÉXITO          ║${NC}"
echo -e "${GREEN}║                                                       ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}📊 Resumen:${NC}"
echo "  ✓ Dependencias descargadas"
echo "  ✓ Código formateado"
echo "  ✓ Análisis estático pasó"
echo "  ✓ Linter ejecutado"
echo "  ✓ Tests unitarios pasaron"
echo "  ✓ Race detection pasó"
echo "  ✓ Cobertura generada"
echo "  ✓ Build verificado"
echo ""
echo -e "${GREEN}🚀 Puedes hacer push con confianza!${NC}"
echo ""
echo -e "${YELLOW}Siguiente paso:${NC}"
echo "  git push origin main"
echo "  git push origin v2.0.0"
echo ""
