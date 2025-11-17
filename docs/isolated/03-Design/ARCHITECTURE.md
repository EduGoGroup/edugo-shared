# Architecture - spec-04
## Módulos Compartidos
```
edugo-shared/
├── logger/ (logging estructurado)
├── database/ (connection helpers)
├── middleware/ (JWT, CORS)
└── testing/ (testcontainers - ya existe)
```
## Versionamiento
- Tags en GitHub (v0.7.0, v0.8.0)
- go.mod en proyectos apunta a versión específica
