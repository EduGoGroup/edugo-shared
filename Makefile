# EduGo Shared - Makefile raiz (multi-modulo)
# Orquesta operaciones en todos los modulos independientes

# Variables
MODULE_NAME = github.com/EduGoGroup/edugo-shared

# Todos los modulos en orden de dependencias
# Nivel 0: sin dependencias internas
# Nivel 1: dependen de nivel 0
# Nivel 2: dependen de nivel 0-1
# Nivel 3: dependen de multiples niveles
MODULES_L0 = common logger config testing messaging/events screenconfig
MODULES_L1 = auth lifecycle
MODULES_L2 = middleware/gin database/postgres database/mongodb messaging/rabbit
MODULES_L3 = bootstrap
MODULES = $(MODULES_L0) $(MODULES_L1) $(MODULES_L2) $(MODULES_L3)

# Colores para output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

.PHONY: help
help: ## Mostrar ayuda de comandos disponibles
	@echo "$(BLUE)EduGo Shared - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Modulos: $(MODULES)$(NC)"

# ============================================================================
# Comandos multi-modulo
# ============================================================================

.PHONY: build-all
build-all: ## Compilar todos los modulos
	@echo "$(BLUE)Compilando todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Building $$module...$(NC)"; \
		(cd $$module && mkdir -p build && go build -o build/ ./...) || exit 1; \
		echo "$(GREEN)  $$module compilado$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos compilados$(NC)"

.PHONY: test-all
test-all: ## Ejecutar tests unitarios en todos los modulos
	@echo "$(BLUE)Ejecutando tests en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Testing $$module...$(NC)"; \
		(cd $$module && go test -short -v ./...) || exit 1; \
		echo "$(GREEN)  $$module tests passed$(NC)"; \
		echo ""; \
	done
	@echo "$(GREEN)Todos los modulos pasaron los tests$(NC)"

.PHONY: test-race-all
test-race-all: ## Ejecutar tests con race detection en todos los modulos
	@echo "$(BLUE)Ejecutando tests con race detection...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Testing $$module with race detection...$(NC)"; \
		(cd $$module && go test -race -short -v ./...) || exit 1; \
		echo "$(GREEN)  $$module passed$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos pasaron los race tests$(NC)"

.PHONY: lint-all
lint-all: ## Ejecutar linter en todos los modulos
	@echo "$(BLUE)Ejecutando linter en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Linting $$module...$(NC)"; \
		(cd $$module && golangci-lint run ./...) || exit 1; \
		echo "$(GREEN)  $$module linted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos pasaron el linter$(NC)"

.PHONY: fmt-all
fmt-all: ## Formatear codigo en todos los modulos
	@echo "$(BLUE)Formateando codigo en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Formatting $$module...$(NC)"; \
		(cd $$module && go fmt ./...); \
		echo "$(GREEN)  $$module formatted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos formateados$(NC)"

.PHONY: vet-all
vet-all: ## Ejecutar go vet en todos los modulos
	@echo "$(BLUE)Ejecutando go vet en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Vetting $$module...$(NC)"; \
		(cd $$module && go vet ./...) || exit 1; \
		echo "$(GREEN)  $$module vetted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos pasaron go vet$(NC)"

.PHONY: tidy-all
tidy-all: ## Ejecutar go mod tidy en todos los modulos
	@echo "$(BLUE)Ejecutando go mod tidy en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Tidying $$module...$(NC)"; \
		(cd $$module && go mod tidy); \
		echo "$(GREEN)  $$module tidied$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos tidied$(NC)"

.PHONY: deps-all
deps-all: ## Actualizar dependencias en todos los modulos
	@echo "$(BLUE)Actualizando dependencias en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Updating $$module...$(NC)"; \
		(cd $$module && go get -u ./... && go mod tidy); \
		echo "$(GREEN)  $$module updated$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos actualizados$(NC)"

.PHONY: check-all
check-all: fmt-all vet-all lint-all test-all ## Validacion completa de todos los modulos
	@echo "$(GREEN)Validacion completa exitosa$(NC)"

.PHONY: ci
ci: fmt-all vet-all test-race-all ## Pipeline CI completo
	@echo "$(GREEN)Pipeline CI completado$(NC)"

.PHONY: clean-all
clean-all: ## Limpiar archivos generados en todos los modulos
	@echo "$(BLUE)Limpiando todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		(cd $$module && rm -rf build && go clean -testcache); \
	done
	@rm -rf build coverage
	@echo "$(GREEN)Limpieza completada$(NC)"

# ============================================================================
# Herramientas
# ============================================================================

.PHONY: install-tools
install-tools: ## Instalar herramientas de desarrollo
	@echo "$(BLUE)Instalando herramientas de desarrollo...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Herramientas instaladas$(NC)"

.PHONY: version
version: ## Mostrar versiones de herramientas
	@echo "$(BLUE)Versiones:$(NC)"
	@echo "Go: $$(go version)"
	@echo "Module: $(MODULE_NAME)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint: $$(golangci-lint version 2>&1 | head -1)"; \
	else \
		echo "golangci-lint: $(RED)no instalado$(NC)"; \
	fi

# ============================================================================
# Pre-commit Hooks
# ============================================================================

.PHONY: setup-hooks
setup-hooks: ## Configurar pre-commit hooks
	@./scripts/setup-hooks.sh

.PHONY: test-hooks
test-hooks: ## Probar pre-commit hooks manualmente
	@.githooks/pre-commit

.DEFAULT_GOAL := help
