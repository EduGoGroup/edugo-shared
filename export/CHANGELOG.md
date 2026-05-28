# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/export` se registran aquí.

## [0.100.0] - 2026-04-02

### Added

- **Document**: Estructura estándar con Header, Sections y Footer para representar documentos exportables
- **Header**: Metadatos del documento (title, subtitle, generated_at, generated_by)
- **Section**: Organización de contenido con título, filas clave-valor y notas opcionales
- **Row**: Pares etiqueta-valor para datos estructurados
- **Footer**: Pie de página opcional para mensajes de cierre
- **Renderer**: Interfaz agnóstica para transformar documentos a múltiples formatos
- **MarkdownRenderer**: Implementación concreta para renderizar documentos a Markdown (text/markdown)
- Suite completa de tests unitarios sin dependencias externas
- Documentación técnica detallada en docs/README.md con componentes, flujos comunes y ejemplos de extensibilidad
- Makefile con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Agnóstico de formato**: La interfaz `Renderer` permite futuras implementaciones (PDF, HTML, Excel) sin modificar código existente
- **Composición simple**: El modelo `Document` es composición de tipos simples, serializable a JSON para almacenamiento
- **Extensibilidad clara**: Nuevos formatos se implementan solo extendiendo la interfaz `Renderer`

## [0.1.0] - 2026-03-17

### Added

- Modulo inicial con contrato `Document` y renderer Markdown (`MarkdownRenderer`).
- Interfaz `Renderer` para extensibilidad futura (PDF, HTML, etc.).
- Tests unitarios con cobertura de todos los campos opcionales.

