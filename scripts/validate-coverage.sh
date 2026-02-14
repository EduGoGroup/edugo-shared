#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
THRESHOLDS_FILE="$PROJECT_ROOT/.coverage-thresholds.yml"

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Validación de Umbrales de Coverage"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ ! -f "$THRESHOLDS_FILE" ]; then
    echo -e "${RED}❌ Archivo de umbrales no encontrado: $THRESHOLDS_FILE${NC}"
    exit 1
fi

# Función para extraer threshold de un módulo del YAML
get_threshold() {
    local module=$1
    # Usar awk para extraer el threshold del YAML
    awk -v mod="$module:" -v RS="" '
        $0 ~ mod {
            for(i=1; i<=NF; i++) {
                if($i == "threshold:") {
                    print $(i+1)
                    exit
                }
            }
        }
    ' "$THRESHOLDS_FILE"
}

# Validar un módulo
validate_module() {
    local module_path="$PROJECT_ROOT/$1"
    local module_name=$(basename $1)
    local mod_download_log
    local test_log

    # Obtener threshold configurado
    local threshold=$(get_threshold "$module_name")

    if [ -z "$threshold" ]; then
        # Usar threshold por defecto si no está definido
        threshold=50
        echo -e "${YELLOW}⚠️  $module_name: Sin umbral definido, usando default $threshold%${NC}"
    fi

    cd "$module_path"

    # Descargar dependencias explícitamente para evitar fallos intermitentes de resolución de módulos
    mod_download_log=$(mktemp)
    if ! go mod download > "$mod_download_log" 2>&1; then
        echo -e "${RED}❌ $module_name: Falló go mod download${NC}"
        echo "   Detalle:"
        sed -n '1,80p' "$mod_download_log" | sed 's/^/   /'
        rm -f "$mod_download_log"
        cd "$PROJECT_ROOT"
        return 1
    fi
    rm -f "$mod_download_log"

    # Ejecutar tests con coverage
    test_log=$(mktemp)
    if ! go test -short ./... -coverprofile=coverage.out -covermode=atomic > "$test_log" 2>&1; then
        echo -e "${RED}❌ $module_name: Tests fallan${NC}"
        echo "   Detalle:"
        sed -n '1,120p' "$test_log" | sed 's/^/   /'
        rm -f "$test_log"
        cd "$PROJECT_ROOT"
        return 1
    fi
    rm -f "$test_log"

    if [ ! -f coverage.out ]; then
        echo -e "${YELLOW}⚠️  $module_name: Sin archivo de coverage${NC}"
        cd "$PROJECT_ROOT"
        return 0
    fi

    # Calcular coverage actual
    local coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $NF}' | sed 's/%//')

    # Limpiar
    rm coverage.out
    cd "$PROJECT_ROOT"

    # Comparar con threshold
    local meets=$(echo "$coverage >= $threshold" | bc -l)

    if [ "$meets" -eq 1 ]; then
        local diff=$(echo "$coverage - $threshold" | bc -l | awk '{printf "%.1f", $0}')
        echo -e "${GREEN}✅ $module_name: ${coverage}% (umbral: ${threshold}%, +${diff}%)${NC}"
        return 0
    else
        local diff=$(echo "$threshold - $coverage" | bc -l | awk '{printf "%.1f", $0}')
        echo -e "${RED}❌ $module_name: ${coverage}% (umbral: ${threshold}%, -${diff}%)${NC}"
        return 1
    fi
}

# Contadores
total=0
passed=0
failed=0

# Módulos a validar
# Nota: common está excluido debido a issue técnico con covdata en Go 1.25
modules=(
    "logger"
    "auth"
    "bootstrap"
    "config"
    "lifecycle"
    "evaluation"
    "testing"
    "database/mongodb"
    "database/postgres"
    "middleware/gin"
    "messaging/rabbit"
)

# Validar cada módulo
for module in "${modules[@]}"; do
    if [ -f "$PROJECT_ROOT/$module/go.mod" ]; then
        total=$((total + 1))
        if validate_module "$module"; then
            passed=$((passed + 1))
        else
            failed=$((failed + 1))
        fi
    fi
done

# Resumen
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Resumen"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "Total de módulos: $total"
echo -e "${GREEN}Pasaron: $passed${NC}"
echo -e "${RED}Fallaron: $failed${NC}"
echo ""

if [ $failed -gt 0 ]; then
    echo -e "${RED}❌ Validación falló: $failed módulo(s) por debajo del umbral${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Todos los módulos cumplen con los umbrales${NC}"
    exit 0
fi
