#!/bin/bash

# Script para calcular coverage de todos los módulos

cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

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

echo "================================"
echo "COVERAGE POR MÓDULO"
echo "================================"

total_coverage=0
count=0

for module in "${modules[@]}"; do
  echo ""
  echo "=== $module ==="
  cd "$module"

  # Ejecutar tests y capturar solo la línea de coverage
  result=$(go test -cover ./... 2>&1 | grep "coverage:")
  echo "$result"

  # Extraer porcentaje si está disponible
  if [[ $result =~ ([0-9]+\.[0-9]+)% ]]; then
    coverage="${BASH_REMATCH[1]}"
    total_coverage=$(echo "$total_coverage + $coverage" | bc)
    count=$((count + 1))
  fi

  cd ..
done

echo ""
echo "================================"
echo "COVERAGE GLOBAL PROMEDIO"
echo "================================"

if [ $count -gt 0 ]; then
  avg_coverage=$(echo "scale=2; $total_coverage / $count" | bc)
  echo "Promedio: ${avg_coverage}% (basado en $count módulos)"
else
  echo "No se pudo calcular coverage"
fi
