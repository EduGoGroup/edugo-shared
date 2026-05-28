# Fase 2

La fase 2 documenta `edugo-shared` como parte del ecosistema EduGo. A diferencia de la fase 1, aqui si se incorporan relaciones externas y dependencias entre repositorios.

## Fuentes usadas en esta fase

- `/Users/jhoanmedina/source/EduGo/Common/ecosistema.md`
- `/Users/jhoanmedina/source/EduGo/EduBack/go.work`
- `go.mod` de `edugo-api-iam-platform`, `edugo-api-admin-new`, `edugo-api-mobile-new` y `edugo-worker`
- `cmd/main.go`, `internal/container/container.go` y piezas de integracion relevantes de esas aplicaciones

## Regla de lectura

- La fase 1 sigue siendo la fuente para entender el interior de cada modulo.
- La fase 2 explica quien consume cada modulo, desde que servicio, y con que rol dentro del ecosistema.
- Cuando una afirmacion depende de imports o `go.mod`, se considera verificada contra codigo local.
- Cuando una afirmacion depende del mapa de topologia general, se considera verificada contra `ecosistema.md`.

## Documentos de esta fase

- [Overview del ecosistema](ecosystem-overview.md)
- [Matriz servicio-modulo](service-module-matrix.md)
- [Consumidores por modulo](module-consumers.md)
- [Flujos de integracion](integration-flows.md)
