#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Análisis de Cobertura - edugo-shared"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

mkdir -p "$PROJECT_ROOT/docs/cicd/coverage-analysis"
OUTPUT_FILE="$PROJECT_ROOT/docs/cicd/coverage-analysis/coverage-report-$(date +%Y%m%d).md"
CURRENT_DATE=$(date '+%Y-%m-%d %H:%M')

cat > "$OUTPUT_FILE" << HEADER
# Reporte de Cobertura - edugo-shared

**Fecha:** $CURRENT_DATE  
**Generado por:** analyze-coverage.sh

---

## 📊 Resumen Ejecutivo

| Módulo | Coverage | Estado | Prioridad |
|--------|----------|--------|-----------|
HEADER

# Función para analizar un módulo
analyze_module() {
    local module_path="$PROJECT_ROOT/$1"
    local module_name=$(basename $1)
    
    echo "Analizando: $module_name..."
    
    if [ ! -d "$module_path" ]; then
        echo "⚠️  Módulo no encontrado: $module_path"
        return
    fi
    
    cd "$module_path"
    
    # Ejecutar tests con coverage
    if go test ./... -coverprofile=coverage.out -covermode=atomic > /dev/null 2>&1; then
        if [ -f coverage.out ]; then
            # Calcular coverage
            coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $NF}' | sed 's/%//')
            
            # Determinar estado
            if (( $(echo "$coverage >= 80" | bc -l) )); then
                status="✅ Excelente"
                priority="Baja"
            elif (( $(echo "$coverage >= 60" | bc -l) )); then
                status="🟢 Bueno"
                priority="Baja"
            elif (( $(echo "$coverage >= 40" | bc -l) )); then
                status="🟡 Aceptable"
                priority="Media"
            elif (( $(echo "$coverage >= 20" | bc -l) )); then
                status="🟠 Bajo"
                priority="Alta"
            else
                status="🔴 Crítico"
                priority="Crítica"
            fi
            
            # Agregar a reporte
            echo "| $module_name | ${coverage}% | $status | $priority |" >> "$OUTPUT_FILE"
            
            # Detalle por archivo
            echo "" >> "$OUTPUT_FILE"
            echo "### $module_name (${coverage}%)" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
            echo '```' >> "$OUTPUT_FILE"
            go tool cover -func=coverage.out >> "$OUTPUT_FILE"
            echo '```' >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
            
            rm coverage.out
        else
            echo "| $module_name | N/A | ⚠️ Sin coverage | - |" >> "$OUTPUT_FILE"
        fi
    else
        echo "| $module_name | ERROR | ❌ Tests fallan | - |" >> "$OUTPUT_FILE"
    fi
    
    cd "$PROJECT_ROOT"
}

# Módulos raíz
for dir in common logger auth bootstrap config lifecycle evaluation testing; do
    if [ -f "$PROJECT_ROOT/$dir/go.mod" ]; then
        analyze_module "$dir"
    fi
done

# Módulos en subdirectorios
for dir in database/mongodb database/postgres middleware/gin messaging/rabbit; do
    if [ -f "$PROJECT_ROOT/$dir/go.mod" ]; then
        analyze_module "$dir"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Análisis completado"
echo "📄 Reporte: $OUTPUT_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
