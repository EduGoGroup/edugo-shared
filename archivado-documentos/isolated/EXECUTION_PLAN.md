# EXECUTION PLAN - SHARED

## Información del Proyecto

**Proyecto:** EduGo Shared Library  
**Objetivo:** Librería Go reutilizable para todos los proyectos  
**Duración:** 4 Sprints (8 semanas)  
**Equipo:** 1-2 Backend engineers  
**Repositorio:** https://github.com/EduGoGroup/edugo-shared

---

## Fase 1: Logger + Database (Sprint 1)

### 1.1 Logger Module
- [ ] Crear interfaz Logger
- [ ] Implementar logger con logrus
- [ ] Agregar contexto (request_id, user_id)
- [ ] JSON output format
- [ ] Niveles: debug, info, warn, error
- [ ] Tests unitarios

### 1.2 Database Module - PostgreSQL
- [ ] Crear PostgreSQL connection manager
- [ ] Connection pooling con GORM
- [ ] Config vía Viper
- [ ] Health check
- [ ] Tests de integración

### 1.3 Database Module - MongoDB
- [ ] Crear MongoDB client wrapper
- [ ] Connection pooling
- [ ] Health check
- [ ] Tests de integración

### 1.4 Documentation
- [ ] README.md
- [ ] Examples de uso
- [ ] Swagger si aplica

**Entregables:**
- [ ] SHARED v1.0.0 tag
- [ ] Go module funcionando
- [ ] Documentación completa

---

## Fase 2: Auth + Messaging (Sprint 2)

### 2.1 Auth Module
- [ ] JWT validator
- [ ] Claims extraction
- [ ] Middleware Gin compatible
- [ ] Key rotation support
- [ ] Tests de validación

### 2.2 Messaging Module
- [ ] RabbitMQ connection manager
- [ ] Publisher interface
- [ ] Subscriber interface
- [ ] Auto-reconnect
- [ ] Dead letter queue support
- [ ] Tests con mock

### 2.3 Integration Tests
- [ ] Test end-to-end con todos servicios
- [ ] Documentación de integración

**Entregables:**
- [ ] SHARED v1.1.0 tag
- [ ] Auth module stable
- [ ] Messaging module stable

---

## Fase 3: Models + Context + Errors (Sprint 3)

### 3.1 Models Module
- [ ] User struct
- [ ] School struct
- [ ] AcademicUnit struct
- [ ] Validation tags
- [ ] JSON serialization

### 3.2 Context Module
- [ ] WithTimeout wrapper
- [ ] User context injection
- [ ] Request ID propagation
- [ ] Tests de timeout

### 3.3 Errors Module
- [ ] Error types (NotFound, BadRequest, etc)
- [ ] HTTP conversion
- [ ] Error wrapping
- [ ] Stack traces
- [ ] Tests de error handling

**Entregables:**
- [ ] SHARED v1.2.0 tag
- [ ] Todos los tipos compartidos
- [ ] Error handling estandarizado

---

## Fase 4: Health Checks + Optimizaciones (Sprint 4)

### 4.1 Health Module
- [ ] Health check interface
- [ ] Postgres checker
- [ ] MongoDB checker
- [ ] RabbitMQ checker
- [ ] Gin handler

### 4.2 Optimizaciones
- [ ] Benchmarks
- [ ] Memory profiling
- [ ] Connection pool tuning
- [ ] Dependency updates

### 4.3 Documentation
- [ ] API documentation
- [ ] Migration guide
- [ ] Best practices

**Entregables:**
- [ ] SHARED v1.3.0 tag (versión de producción)
- [ ] Documentación completa
- [ ] Health check endpoints

---

## Criterios de Aceptación

- [ ] Todos los módulos compilables
- [ ] Tests >= 80% cobertura
- [ ] Documentación 100%
- [ ] Ejemplos funcionales
- [ ] Compatible con todos proyectos
- [ ] Versionamiento semántico implementado

---

## Release Process

```bash
# 1. Cambios en rama dev
# 2. Tests locales pasen
# 3. Tag version
git tag v1.3.0

# 4. Push
git push origin v1.3.0

# 5. Otros proyectos actualizan:
go get -u github.com/EduGoGroup/edugo-shared@v1.3.0
```

---

## Riesgos

| Riesgo | Mitigation |
|--------|-----------|
| Incompatibilidad en otros proyectos | Tests de integración tempranos |
| Cambios de API frecuentes | Mantener v1.x stable, v2.x para breaking |
| Performance degradation | Benchmarks en cada sprint |

---

**Próxima revisión:** Después de Sprint 1  
**Última actualización:** 15 de Noviembre, 2025
