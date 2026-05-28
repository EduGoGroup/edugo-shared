#!/bin/bash

# Script para calcular coverage de todos los módulos

# Detectar módulos dinámicamente buscando archivos go.mod
# Excluir el directorio actual (.) si tiene go.mod pero es el root (aunque en este repo no tiene)
modules=$(find . -name "go.mod" -not -path "./go.mod" | xargs -n 1 dirname | sort)

echo "================================"
echo "COVERAGE POR MÓDULO"
echo "================================"

total_coverage=0
count=0

for module in $modules; do
  # Limpiar el path (quitar ./ al inicio)
  module=${module#./}

  echo ""
  echo "=== $module ==="

  # Ir al directorio del módulo
  pushd "$module" > /dev/null || continue

  # Ejecutar tests y capturar solo la línea de coverage
  # Usamos grep para buscar "coverage:" y awk para extraer el porcentaje
  # go test -cover ./... puede generar múltiples líneas de coverage si hay subpaquetes.
  # Tomaremos el promedio del módulo si hay múltiples paquetes, o el total si es uno solo.

  # Ejecutamos los tests y guardamos la salida
  output=$(go test -cover ./... 2>&1)

  # Mostramos la salida de coverage para información
  echo "$output" | grep "coverage:"

  # Extraer porcentajes
  # Formato típico 1: "coverage: 87.7% of statements"
  # Formato típico 2: "ok package_name 0.00s coverage: 87.7% of statements"
  percentages=$(echo "$output" | grep "coverage:" | grep -o "[0-9]\+\.[0-9]\+%" | sed 's/%//')

  module_total=0
  module_count=0

  for p in $percentages; do
    module_total=$(awk "BEGIN {print $module_total + $p}")
    module_count=$((module_count + 1))
  done

  if [ $module_count -gt 0 ]; then
    module_avg=$(awk "BEGIN {print $module_total / $module_count}")
    echo "  -> Promedio del módulo: ${module_avg}%"

    total_coverage=$(awk "BEGIN {print $total_coverage + $module_avg}")
    count=$((count + 1))
  else
    echo "  -> No coverage data found"
  fi

  popd > /dev/null
done

echo ""
echo "================================"
echo "COVERAGE GLOBAL PROMEDIO"
echo "================================"

if [ $count -gt 0 ]; then
  avg_coverage=$(awk "BEGIN {print $total_coverage / $count}")
  echo "Promedio Global: ${avg_coverage}% (basado en $count módulos)"
else
  echo "No se pudo calcular coverage"
fi
