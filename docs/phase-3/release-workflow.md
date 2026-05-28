# Flujo de release modular

## Convencion de tags

- Release raiz: `vX.Y.Z`
- Release de modulo simple: `modulo/vX.Y.Z`
- Release de modulo anidado: `dominio/modulo/vX.Y.Z`

## Flujo recomendado

1. Preparar changelog del modulo.
2. Confirmar y publicar el cambio en Git.
3. Crear el tag modular desde el `Makefile` del modulo o desde la raiz.
4. Dejar que `release.yml` valide y cree el GitHub release.

## Ejemplos

Desde el modulo:

```bash
make changelog VERSION=v0.4.0
make release VERSION=v0.4.0
```

Desde la raiz:

```bash
make changelog-module MODULE=cache/redis VERSION=v0.4.0
make release-module MODULE=cache/redis VERSION=v0.4.0
```

## Notas operativas

- `make changelog` usa `scripts/update-module-changelog.sh` y versiona la seccion `Unreleased` del modulo.
- `make release` exige arbol Git limpio, crea el tag `<modulo>/vX.Y.Z` y lo empuja a `origin`.
- El GitHub release se genera en CI leyendo el `CHANGELOG.md` del modulo correspondiente.
- Si el tag es raiz, el release vuelve a usar `CHANGELOG.md` de la raiz y valida todos los modulos.
