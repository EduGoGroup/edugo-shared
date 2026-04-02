// Package bootstrap coordina la inicializacion de recursos de infraestructura.
//
// # API Publica
//
// Tipos principales:
//   - [Resources] contiene los recursos inicializados (logger, DB, MQ, storage)
//   - [Factories] agrupa las factory interfaces para crear recursos
//   - [BootstrapOptions] configura el comportamiento del bootstrap
//
// Interfaces de factory:
//   - [LoggerFactory], [PostgreSQLFactory], [MongoDBFactory], [RabbitMQFactory], [S3Factory]
//
// Interfaces de recurso:
//   - [MessagePublisher], [StorageClient], [DatabaseClient], [HealthChecker]
//
// Constructores de factories por defecto:
//   - [NewDefaultLoggerFactory] — logrus
//   - [NewDefaultPostgreSQLFactory] — GORM
//   - [NewDefaultMongoDBFactory] — mongo-driver v2
//   - [NewDefaultRabbitMQFactory] — amqp091
//   - [NewDefaultS3Factory] — AWS SDK v2
//
// Opciones funcionales:
//   - [WithRequiredResources], [WithOptionalResources], [WithSkipHealthCheck],
//     [WithMockFactories], [WithStopOnFirstError]
//
// Funcion principal:
//   - [Bootstrap] orquesta la inicializacion completa
//
// Las implementaciones concretas de las factories son tipos no exportados.
// Los constructores retornan las interfaces correspondientes, permitiendo
// inyeccion de dependencias y testing sin acoplar a implementaciones especificas.
package bootstrap
