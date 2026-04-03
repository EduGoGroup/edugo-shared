// Package bootstrap define tipos compartidos para la inicializacion de
// recursos de infraestructura.
//
// # Arquitectura
//
// El modulo raiz contiene configs, opciones y tipos de error.
// Cada tecnologia tiene su propio sub-modulo con la implementacion:
//
//   - [bootstrap/postgres] — Factory PostgreSQL + GORM
//   - [bootstrap/mongodb]  — Factory MongoDB
//   - [bootstrap/rabbitmq] — Factory RabbitMQ
//   - [bootstrap/s3]       — Factory S3
//
// # Tipos principales
//
//   - [PostgreSQLConfig], [MongoDBConfig], [RabbitMQConfig], [S3Config] — configuracion
//   - [GORMOption], [WithGORMLogger], [WithSimpleProtocol] — opciones funcionales para GORM
//   - [LifecycleManager] — interface para registro de cleanup
//   - [ErrMissingFactory], [ErrConnectionFailed] — errores tipados
package bootstrap
