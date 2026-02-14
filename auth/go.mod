module github.com/EduGoGroup/edugo-shared/auth

go 1.25

require (
	github.com/EduGoGroup/edugo-shared/common v0.3.3
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.46.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/EduGoGroup/edugo-shared/common => ../common
