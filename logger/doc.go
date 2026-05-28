// Package logger provee una interfaz unificada de logging con multiples backends.
//
// # API Publica
//
// Interfaz principal:
//   - [Logger] — interfaz comun con metodos Debug, Info, Warn, Error, Fatal, With, Sync
//   - [Fields] — alias de map[string]any para campos estructurados
//
// Constructores (retornan [Logger]):
//   - [NewZapLogger] — backend Uber Zap (alto rendimiento, JSON/console)
//   - [NewLogrusLogger] — backend Logrus
//   - [NewSlogAdapter] — wrapper de *slog.Logger para satisfacer [Logger]
//
// Slog nativo:
//   - [NewSlogProvider] — crea *slog.Logger configurado con [SlogConfig]
//   - [NewSlogProviderFromEnv] — crea *slog.Logger desde variables de entorno
//   - [SlogAdapter.SlogLogger] — accede al *slog.Logger subyacente
//
// Contexto:
//   - [NewContext] — embede *slog.Logger en context.Context
//   - [FromContext] / [L] — extrae *slog.Logger del contexto
//
// Campos estandarizados:
//   - Constantes Field* (FieldUserID, FieldRequestID, FieldError, etc.)
//   - Helpers With* (WithRequestID, WithError, WithDuration, etc.)
//
// Las implementaciones concretas (zapLogger, logrusLogger) son tipos no exportados.
// Los constructores retornan la interfaz [Logger], desacoplando a los consumidores
// de las implementaciones especificas.
package logger
