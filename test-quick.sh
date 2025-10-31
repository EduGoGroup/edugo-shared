#!/bin/bash

# Script RÁPIDO para verificación antes de commit
# Ejecuta solo los checks esenciales (no linter, no cobertura completa)

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo ""
echo -e "${BLUE}⚡ VERIFICACIÓN RÁPIDA PRE-COMMIT${NC}"
echo ""

echo "1/4 Formato..."
make fmt > /dev/null 2>&1

echo "2/4 Análisis estático..."
make vet > /dev/null 2>&1

echo "3/4 Tests rápidos..."
make test-short > /dev/null 2>&1

echo "4/4 Build..."
make build > /dev/null 2>&1

echo ""
echo -e "${GREEN}✓ Verificación rápida completada (ready to commit!)${NC}"
echo ""
