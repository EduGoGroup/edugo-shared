#!/bin/bash
#
# Script para configurar pre-commit hooks
# Uso: ./scripts/setup-hooks.sh
#

set -e

echo "🔧 Configurando pre-commit hooks para edugo-shared..."
echo ""

# 1. Configurar Git hooks path
git config core.hooksPath .githooks

# 2. Hacer ejecutables todos los hooks
chmod +x .githooks/*

# 3. Verificar golangci-lint
if ! command -v golangci-lint &> /dev/null; then
  echo "⚠️  golangci-lint no está instalado"
  echo ""
  echo "Instalación recomendada:"
  echo "  macOS: brew install golangci-lint"
  echo "  Linux: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
  echo ""
  echo "Los hooks funcionarán sin él, pero algunos checks serán saltados."
  echo ""
else
  echo "✅ golangci-lint instalado: $(golangci-lint --version | head -1)"
fi

# 4. Verificar gofmt
if ! command -v gofmt &> /dev/null; then
  echo "❌ gofmt no encontrado (debería estar incluido con Go)"
  exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Hooks configurados exitosamente"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Los siguientes checks se ejecutarán antes de cada commit:"
echo "  • gofmt (formato)"
echo "  • go vet (análisis estático)"
echo "  • golangci-lint (linter avanzado)"
echo "  • go test -short (tests rápidos)"
echo "  • Detección de sensitive data"
echo ""
echo "Para saltear hooks en un commit específico:"
echo "  git commit --no-verify -m \"mensaje\""
echo ""
