module github.com/EduGoGroup/edugo-shared/common

go 1.27.4

// Módulo common: dependencias mínimas
// Solo google/uuid (librería estándar ligera)

require (
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
