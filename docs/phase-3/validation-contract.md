# Contrato de validacion

## Fuente de verdad

- `scripts/module-manifest.tsv` define modulo, nivel de dependencia, necesidad de integracion y participacion en coverage validation.
- `scripts/list-modules.sh` expone ese inventario para `Makefile`, scripts y GitHub Actions.

## Sets operativos

- `level-0` a `level-3`: ordenan pruebas paralelas segun dependencias internas observadas.
- `integration`: modulos con tests que requieren Docker o infraestructura auxiliar.
- `coverage-validation`: modulos que entran al chequeo automatico de umbrales.
- `all`: catalogo completo de modulos Go del repositorio.

## Comandos raiz

- `make build-all`
- `make test-all`
- `make test-race-all`
- `make test-integration-all`
- `make lint-all`
- `make vet-all`
- `make check-all`
- `make changelog-module MODULE=<ruta> VERSION=vX.Y.Z`
- `make release-module MODULE=<ruta> VERSION=vX.Y.Z`

## Comandos por modulo

Cada modulo expone el mismo contrato:

- `make build`
- `make test`
- `make test-all`
- `make test-race`
- `make fmt`
- `make vet`
- `make lint`
- `make check`
- `make changelog VERSION=vX.Y.Z`
- `make release VERSION=vX.Y.Z`

## Integracion con CI

- `ci.yml` y `test.yml` calculan matrices desde `scripts/list-modules.sh`.
- `coverage-validation.yml` usa el mismo inventario y deja de depender de listas hardcodeadas.
- `release.yml` valida solo el modulo objetivo cuando el tag es modular; para tags raiz valida todo el repositorio.
