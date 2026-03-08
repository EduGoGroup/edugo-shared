# Config

Carga configuracion YAML con Viper, la mezcla con variables de entorno y la valida con tags struct.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/config`
- Carpeta: `config`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/config
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir un `Loader` con path, nombre, tipo de archivo y prefijo de entorno.
2. Leer un archivo de configuracion si existe y tolerar ausencia del archivo en `Load`.
3. Aplicar override de variables de entorno mediante `AutomaticEnv` y `SetEnvKeyReplacer`.
4. Unmarshal la configuracion en structs como `BaseConfig`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El loader tolera la ausencia del archivo en `Load`, pero `LoadFromFile` exige que exista.
- El modulo documenta un shape de configuracion base, no una configuracion universal obligatoria para todos los servicios.
- Tiene tests unitarios sobre carga y validacion.
