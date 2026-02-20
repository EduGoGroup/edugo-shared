# Makefile para config
MODULE_NAME = github.com/EduGoGroup/edugo-shared/config

RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

.PHONY: help
help: ## Mostrar ayuda
	@echo "$(BLUE)$(MODULE_NAME) - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

.PHONY: build
build: ## Verificar que el modulo compila (libreria, sin binario)
	@echo "$(BLUE)Compilando $(MODULE_NAME)...$(NC)"
	@go build ./...
	@echo "$(GREEN)Compilacion exitosa$(NC)"

.PHONY: test
test: ## Ejecutar tests unitarios
	@echo "$(BLUE)Ejecutando tests...$(NC)"
	@go test -short -v ./...
	@echo "$(GREEN)Tests completados$(NC)"

.PHONY: test-all
test-all: ## Ejecutar todos los tests (incluye integracion)
	@echo "$(BLUE)Ejecutando todos los tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)Tests completados$(NC)"

.PHONY: test-race
test-race: ## Tests con race detection
	@echo "$(BLUE)Ejecutando tests con race detection...$(NC)"
	@go test -race -short -v ./...
	@echo "$(GREEN)Tests con race detection completados$(NC)"

.PHONY: lint
lint: ## Ejecutar linter
	@echo "$(BLUE)Ejecutando linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)Linter completado$(NC)"

.PHONY: fmt
fmt: ## Formatear codigo
	@echo "$(BLUE)Formateando codigo...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Codigo formateado$(NC)"

.PHONY: vet
vet: ## Analisis estatico
	@echo "$(BLUE)Ejecutando go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Analisis completado$(NC)"

.PHONY: tidy
tidy: ## Limpiar go.mod
	@echo "$(BLUE)Ejecutando go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)go.mod limpio$(NC)"

.PHONY: deps
deps: ## Actualizar dependencias
	@echo "$(BLUE)Actualizando dependencias...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Dependencias actualizadas$(NC)"

.PHONY: check
check: fmt vet lint test ## Validacion completa

.PHONY: clean
clean: ## Limpiar archivos generados
	@go clean -testcache
	@echo "$(GREEN)Limpieza completada$(NC)"

.DEFAULT_GOAL := help
