# EduGo Shared Makefile
# Comandos básicos para desarrollo Go

# Variables
GO_VERSION = 1.25.3
MODULE_NAME = github.com/EduGoGroup/edugo-shared
BINARY_NAME = edugo-shared
BUILD_DIR = build
COVERAGE_DIR = coverage

# Colores para output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

.PHONY: help
help: ## Mostrar ayuda de comandos disponibles
	@echo "$(BLUE)EduGo Shared - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

.PHONY: setup
setup: ## Configurar entorno de desarrollo
	@echo "$(BLUE)Configurando entorno de desarrollo...$(NC)"
	@go version
	@go mod download
	@go mod verify
	@echo "$(GREEN)✓ Entorno configurado correctamente$(NC)"

.PHONY: deps
deps: ## Descargar e instalar dependencias
	@echo "$(BLUE)Descargando dependencias...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencias actualizadas$(NC)"

.PHONY: build
build: ## Verificar que el proyecto compila correctamente
	@echo "$(BLUE)Verificando compilación...$(NC)"
	@go build -v ./...
	@echo "$(GREEN)✓ Compilación exitosa$(NC)"

.PHONY: test
test: ## Ejecutar tests unitarios
	@echo "$(BLUE)Ejecutando tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests completados$(NC)"

.PHONY: test-race
test-race: ## Ejecutar tests con detección de race conditions
	@echo "$(BLUE)Ejecutando tests con race detection...$(NC)"
	@go test -race -v ./...
	@echo "$(GREEN)✓ Tests con race detection completados$(NC)"

.PHONY: test-coverage test-coverage-critical test-coverage-all
test-coverage: test-coverage-critical ## Ejecutar tests con cobertura (solo paquetes críticos)

test-coverage-critical: ## Ejecutar tests con cobertura SOLO en paquetes críticos
	@echo "$(BLUE)Ejecutando tests con cobertura (solo paquetes críticos)...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@echo "$(YELLOW)⭐ Paquetes críticos: auth, database, logger, messaging, validator$(NC)"
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/critical.out -covermode=atomic \
		./pkg/auth/... ./pkg/database/... ./pkg/logger/... ./pkg/messaging/... ./pkg/validator/...
	@go tool cover -html=$(COVERAGE_DIR)/critical.out -o $(COVERAGE_DIR)/coverage-critical.html
	@echo "$(BLUE)📊 Cobertura de paquetes críticos:$(NC)"
	@go tool cover -func=$(COVERAGE_DIR)/critical.out | tail -1
	@echo "$(GREEN)✓ Reporte crítico generado en $(COVERAGE_DIR)/coverage-critical.html$(NC)"

test-coverage-all: ## Ejecutar tests con cobertura en TODOS los paquetes
	@echo "$(BLUE)Ejecutando tests con cobertura completa...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/all.out -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_DIR)/all.out -o $(COVERAGE_DIR)/coverage-all.html
	@echo "$(BLUE)📊 Cobertura completa (incluye config/errors/types):$(NC)"
	@go tool cover -func=$(COVERAGE_DIR)/all.out | tail -1
	@echo "$(GREEN)✓ Reporte completo generado en $(COVERAGE_DIR)/coverage-all.html$(NC)"

.PHONY: test-short
test-short: ## Ejecutar tests cortos (skip tests largos)
	@echo "$(BLUE)Ejecutando tests cortos...$(NC)"
	@go test -short ./...
	@echo "$(GREEN)✓ Tests cortos completados$(NC)"

.PHONY: benchmark
benchmark: ## Ejecutar benchmarks
	@echo "$(BLUE)Ejecutando benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completados$(NC)"

.PHONY: lint
lint: ## Ejecutar linter (requiere golangci-lint)
	@echo "$(BLUE)Ejecutando linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✓ Linter completado$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint no está instalado. Instalando...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
		echo "$(GREEN)✓ Linter completado$(NC)"; \
	fi

.PHONY: fmt
fmt: ## Formatear código
	@echo "$(BLUE)Formateando código...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Código formateado$(NC)"

.PHONY: vet
vet: ## Ejecutar go vet (análisis estático)
	@echo "$(BLUE)Ejecutando análisis estático...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Análisis estático completado$(NC)"

.PHONY: mod-verify
mod-verify: ## Verificar integridad de módulos
	@echo "$(BLUE)Verificando integridad de módulos...$(NC)"
	@go mod verify
	@echo "$(GREEN)✓ Módulos verificados$(NC)"

.PHONY: mod-tidy
mod-tidy: ## Limpiar y organizar go.mod
	@echo "$(BLUE)Limpiando go.mod...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ go.mod limpio$(NC)"

.PHONY: mod-upgrade
mod-upgrade: ## Actualizar todas las dependencias
	@echo "$(BLUE)Actualizando dependencias...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependencias actualizadas$(NC)"

.PHONY: clean
clean: ## Limpiar archivos generados
	@echo "$(BLUE)Limpiando archivos generados...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean -cache
	@go clean -testcache
	@go clean -modcache
	@echo "$(GREEN)✓ Limpieza completada$(NC)"

.PHONY: coverage-info
coverage-info: ## Mostrar información sobre configuración de cobertura
	@echo "$(BLUE)📋 Configuración de Cobertura:$(NC)"
	@echo ""
	@echo "$(GREEN)✅ Paquetes CRÍTICOS (deben tener buena cobertura):$(NC)"
	@echo "  🔐 pkg/auth/     - Autenticación JWT"
	@echo "  🗄️  pkg/database/ - Conexiones y transacciones"  
	@echo "  📝 pkg/logger/   - Configuración de logging"
	@echo "  📨 pkg/messaging/- Publisher/Consumer"
	@echo "  ✅ pkg/validator/- Validaciones de entrada"
	@echo ""
	@echo "$(YELLOW)⚠️  Paquetes EXCLUIDOS (no afectan cobertura crítica):$(NC)"
	@echo "  ⚙️  pkg/config/   - Solo getters de env vars"
	@echo "  ❌ pkg/errors/   - Solo constructores de errores"
	@echo "  🏷️  pkg/types/enum/ - Solo constantes y métodos simples"
	@echo ""
	@echo "$(BLUE)📊 Comandos disponibles:$(NC)"
	@echo "  make test-coverage-critical  - Solo paquetes críticos"
	@echo "  make test-coverage-all      - Todos los paquetes (informativo)"
	@echo "  make test-coverage          - Alias para critical"
	@echo ""
	@echo "Ver configuración completa en: .testcoverage.yml"

.PHONY: docs
docs: ## Generar documentación
	@echo "$(BLUE)Generando documentación...$(NC)"
	@go doc -all ./...
	@echo "$(GREEN)✓ Documentación generada$(NC)"

.PHONY: docs-serve
docs-serve: ## Servir documentación en localhost:6060
	@echo "$(BLUE)Sirviendo documentación en http://localhost:6060$(NC)"
	@echo "$(YELLOW)Presiona Ctrl+C para detener$(NC)"
	@godoc -http=:6060

.PHONY: security
security: ## Ejecutar análisis de seguridad (requiere gosec)
	@echo "$(BLUE)Ejecutando análisis de seguridad...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
		echo "$(GREEN)✓ Análisis de seguridad completado$(NC)"; \
	else \
		echo "$(YELLOW)⚠ gosec no está instalado. Instalando...$(NC)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
		echo "$(GREEN)✓ Análisis de seguridad completado$(NC)"; \
	fi

.PHONY: check-all
check-all: fmt vet lint test security ## Ejecutar todas las verificaciones
	@echo "$(GREEN)✓ Todas las verificaciones completadas$(NC)"

.PHONY: ci
ci: deps fmt vet lint test-race test-coverage ## Pipeline CI completo
	@echo "$(GREEN)✓ Pipeline CI completado$(NC)"

.PHONY: pre-commit
pre-commit: fmt vet lint test-short ## Verificaciones rápidas antes de commit
	@echo "$(GREEN)✓ Verificaciones pre-commit completadas$(NC)"

.PHONY: install-tools
install-tools: ## Instalar herramientas de desarrollo
	@echo "$(BLUE)Instalando herramientas de desarrollo...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "$(GREEN)✓ Herramientas instaladas$(NC)"

.PHONY: version
version: ## Mostrar versiones de herramientas
	@echo "$(BLUE)Versiones de herramientas:$(NC)"
	@echo "Go: $(shell go version)"
	@echo "Module: $(MODULE_NAME)"
	@echo "Build Dir: $(BUILD_DIR)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint: $(shell golangci-lint version)"; \
	else \
		echo "golangci-lint: $(RED)no instalado$(NC)"; \
	fi
	@if command -v gosec >/dev/null 2>&1; then \
		echo "gosec: $(shell gosec -version 2>/dev/null | head -n1)"; \
	else \
		echo "gosec: $(RED)no instalado$(NC)"; \
	fi

# Comando por defecto
.DEFAULT_GOAL := help