package logger

// Constantes para campos de log estandarizados
// Estas constantes facilitan el uso consistente de campos en toda la aplicación
const (
	// FieldUserID es el ID del usuario que realiza la acción
	FieldUserID = "user_id"

	// FieldRequestID es el ID único de la petición HTTP
	FieldRequestID = "request_id"

	// FieldTraceID es el ID de tracing distribuido
	FieldTraceID = "trace_id"

	// FieldSpanID es el ID del span en tracing
	FieldSpanID = "span_id"

	// FieldMethod es el método HTTP de la petición
	FieldMethod = "method"

	// FieldPath es la ruta HTTP de la petición
	FieldPath = "path"

	// FieldStatusCode es el código de estado HTTP de la respuesta
	FieldStatusCode = "status_code"

	// FieldDuration es la duración de la operación en milisegundos
	FieldDuration = "duration_ms"

	// FieldError es el mensaje de error
	FieldError = "error"

	// FieldErrorType es el tipo de error
	FieldErrorType = "error_type"

	// FieldErrorStack es el stack trace del error
	FieldErrorStack = "error_stack"

	// FieldComponent es el componente/módulo que genera el log
	FieldComponent = "component"

	// FieldFunction es la función que genera el log
	FieldFunction = "function"

	// FieldHost es el hostname del servidor
	FieldHost = "host"

	// FieldIP es la dirección IP del cliente
	FieldIP = "ip"

	// FieldEnvironment es el entorno de ejecución (dev, qa, prod)
	FieldEnvironment = "env"

	// FieldVersion es la versión de la aplicación
	FieldVersion = "version"

	// FieldDatabase es el nombre de la base de datos
	FieldDatabase = "database"

	// FieldTable es el nombre de la tabla
	FieldTable = "table"

	// FieldQuery es la consulta SQL
	FieldQuery = "query"

	// FieldQueryDuration es la duración de la consulta en milisegundos
	FieldQueryDuration = "query_duration_ms"

	// FieldEvent es el tipo de evento
	FieldEvent = "event"

	// FieldEventID es el ID del evento
	FieldEventID = "event_id"

	// FieldQueueName es el nombre de la cola de mensajes
	FieldQueueName = "queue_name"

	// FieldMessageID es el ID del mensaje
	FieldMessageID = "message_id"

	// FieldCorrelationID es el ID de correlación entre servicios
	FieldCorrelationID = "correlation_id"

	// FieldService es el nombre del servicio
	FieldService = "service"

	// FieldAction es la acción que se está ejecutando
	FieldAction = "action"

	// FieldResource es el recurso afectado
	FieldResource = "resource"

	// FieldResourceID es el ID del recurso
	FieldResourceID = "resource_id"

	// FieldSessionID es el ID de la sesión
	FieldSessionID = "session_id"

	// FieldTenantID es el ID del tenant en sistemas multi-tenant
	FieldTenantID = "tenant_id"

	// FieldBucket es el nombre del bucket de almacenamiento
	FieldBucket = "bucket"

	// FieldKey es la clave del objeto en almacenamiento
	FieldKey = "key"

	// FieldFileSize es el tamaño del archivo en bytes
	FieldFileSize = "file_size"

	// FieldContentType es el tipo de contenido MIME
	FieldContentType = "content_type"
)
