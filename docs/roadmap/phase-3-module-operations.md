# Fase 3 - Operacion por modulo

Estado: implementada.

## Resultado cerrado

1. El `Makefile` raiz ya cubre los 17 modulos reales del repositorio, incluyendo `cache/redis` y `repository`.
2. Los `Makefile` locales quedaron homogeneos mediante `scripts/module-common.mk`.
3. CI, coverage y release leen el mismo inventario de modulos desde `scripts/module-manifest.tsv`.
4. El workflow `release.yml` soporta releases raiz y releases por modulo usando el `CHANGELOG.md` correcto.
5. Cada modulo ahora puede versionar su changelog y publicar su tag modular con `make changelog` y `make release`.

## Documentacion operativa

- [Overview de fase 3](../phase-3/README.md)
- [Contrato de validacion](../phase-3/validation-contract.md)
- [Flujo de release modular](../phase-3/release-workflow.md)

## Riesgos abiertos

- `audit/postgres` sigue sin tests propios y por eso mantiene umbral de coverage `0`.
- `repository` ya entro al pipeline central, pero necesita ampliar cobertura hacia CRUD e integracion.
- `common` sigue fuera de `coverage-validation` hasta cerrar definitivamente el comportamiento de covdata en CI.
