# Reporte de Cobertura - edugo-shared

**Fecha:** $(date '+%Y-%m-%d %H:%M')  
**Generado por:** analyze-coverage.sh

---

## 📊 Resumen Ejecutivo

| Módulo | Coverage | Estado | Prioridad |
|--------|----------|--------|-----------|
| common | ERROR | ❌ Tests fallan | - |
| logger | 95.8% | ✅ Excelente | Baja |

### logger (95.8%)

```
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:18:	NewZapLogger	100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:80:	Debug		100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:85:	Info		100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:90:	Warn		100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:95:	Error		100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:100:	Fatal		0.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:105:	With		100.0%
github.com/EduGoGroup/edugo-shared/logger/zap_logger.go:112:	Sync		100.0%
total:								(statements)	95.8%
```

| auth | 85.0% | ✅ Excelente | Baja |

### auth (85.0%)

```
github.com/EduGoGroup/edugo-shared/auth/jwt.go:32:		NewJWTManager		100.0%
github.com/EduGoGroup/edugo-shared/auth/jwt.go:40:		GenerateToken		87.5%
github.com/EduGoGroup/edugo-shared/auth/jwt.go:69:		ValidateToken		76.5%
github.com/EduGoGroup/edugo-shared/auth/jwt.go:106:		RefreshToken		100.0%
github.com/EduGoGroup/edugo-shared/auth/jwt.go:117:		ExtractUserID		85.7%
github.com/EduGoGroup/edugo-shared/auth/jwt.go:133:		ExtractRole		85.7%
github.com/EduGoGroup/edugo-shared/auth/password.go:19:		HashPassword		83.3%
github.com/EduGoGroup/edugo-shared/auth/password.go:33:		VerifyPassword		100.0%
github.com/EduGoGroup/edugo-shared/auth/refresh_token.go:22:	GenerateRefreshToken	85.7%
github.com/EduGoGroup/edugo-shared/auth/refresh_token.go:45:	HashToken		100.0%
total:								(statements)		85.0%
```

| bootstrap | 29.5% | 🟠 Bajo | Alta |

### bootstrap (29.5%)

```
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:21:			Bootstrap			85.3%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:110:			initLogger			87.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:138:			initPostgreSQL			25.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:190:			initMongoDB			23.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:241:			initRabbitMQ			20.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:300:			initS3				28.6%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:351:			performHealthChecks		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:391:			mergeFactories			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:416:			isRequired			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:425:			logWarning			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:433:			extractEnvAndVersion		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:439:			extractPostgreSQLConfig		38.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:468:			extractMongoDBConfig		38.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:497:			extractRabbitMQConfig		38.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:526:			extractS3Config			38.5%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:555:			registerPostgreSQLCleanup	0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:576:			registerMongoDBCleanup		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/bootstrap.go:599:			registerRabbitMQCleanup		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_logger.go:18:		NewDefaultLoggerFactory		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_logger.go:23:		CreateLogger			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_mongodb.go:23:		NewDefaultMongoDBFactory	100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_mongodb.go:30:		CreateConnection		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_mongodb.go:60:		GetDatabase			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_mongodb.go:65:		Ping				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_mongodb.go:79:		Close				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:24:		NewDefaultPostgreSQLFactory	100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:34:		CreateConnection		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:68:		CreateRawConnection		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:92:		Ping				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:106:		Close				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_postgresql.go:120:		buildDSN			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_rabbitmq.go:21:		NewDefaultRabbitMQFactory	100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_rabbitmq.go:28:		CreateConnection		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_rabbitmq.go:55:		CreateChannel			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_rabbitmq.go:75:		DeclareQueue			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_rabbitmq.go:98:		Close				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_s3.go:21:			NewDefaultS3Factory		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_s3.go:26:			CreateClient			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_s3.go:55:			CreatePresignClient		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/factory_s3.go:60:			ValidateBucket			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/interfaces.go:190:			Validate			53.8%
github.com/EduGoGroup/edugo-shared/bootstrap/interfaces.go:227:			Error				100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:42:			WithRequiredResources		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:49:			WithOptionalResources		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:56:			WithSkipHealthCheck		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:63:			WithMockFactories		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:70:			WithStopOnFirstError		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:81:			DefaultBootstrapOptions		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/options.go:92:			ApplyOptions			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:25:	Publish				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:30:	PublishWithPriority		0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:59:	Close				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:78:	Upload				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:95:	Download			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:115:	Delete				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:128:	GetPresignedURL			0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resource_implementations.go:134:	Exists				0.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resources.go:35:			HasLogger			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resources.go:40:			HasPostgreSQL			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resources.go:45:			HasMongoDB			100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resources.go:50:			HasMessagePublisher		100.0%
github.com/EduGoGroup/edugo-shared/bootstrap/resources.go:55:			HasStorageClient		100.0%
total:										(statements)			29.5%
```

| config | 82.9% | ✅ Excelente | Baja |

### config (82.9%)

```
github.com/EduGoGroup/edugo-shared/config/base.go:69:		ConnectionString	100.0%
github.com/EduGoGroup/edugo-shared/config/base.go:74:		ConnectionStringWithDB	100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:22:		WithConfigPath		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:29:		WithConfigName		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:36:		WithConfigType		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:43:		WithEnvPrefix		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:50:		NewLoader		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:66:		Load			84.6%
github.com/EduGoGroup/edugo-shared/config/loader.go:96:		LoadFromFile		88.9%
github.com/EduGoGroup/edugo-shared/config/loader.go:114:	Get			100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:119:	GetString		100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:124:	GetInt			100.0%
github.com/EduGoGroup/edugo-shared/config/loader.go:129:	GetBool			100.0%
github.com/EduGoGroup/edugo-shared/config/validator.go:15:	NewValidator		100.0%
github.com/EduGoGroup/edugo-shared/config/validator.go:22:	Validate		80.0%
github.com/EduGoGroup/edugo-shared/config/validator.go:33:	ValidateField		0.0%
github.com/EduGoGroup/edugo-shared/config/validator.go:54:	NewValidationError	100.0%
github.com/EduGoGroup/edugo-shared/config/validator.go:72:	Error			83.3%
github.com/EduGoGroup/edugo-shared/config/validator.go:86:	buildErrorMessage	63.6%
total:								(statements)		82.9%
```

| lifecycle | 91.8% | ✅ Excelente | Baja |

### lifecycle (91.8%)

```
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:29:	NewManager	100.0%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:39:	Register	100.0%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:56:	RegisterSimple	100.0%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:61:	Startup		86.7%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:103:	Cleanup		90.5%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:157:	Count		100.0%
github.com/EduGoGroup/edugo-shared/lifecycle/manager.go:165:	Clear		100.0%
total:								(statements)	91.8%
```

| evaluation | 100.0% | ✅ Excelente | Baja |

### evaluation (100.0%)

```
github.com/EduGoGroup/edugo-shared/evaluation/assessment.go:30:	Validate		100.0%
github.com/EduGoGroup/edugo-shared/evaluation/assessment.go:47:	IsPublished		100.0%
github.com/EduGoGroup/edugo-shared/evaluation/attempt.go:35:	CalculatePercentage	100.0%
github.com/EduGoGroup/edugo-shared/evaluation/attempt.go:44:	CheckPassed		100.0%
github.com/EduGoGroup/edugo-shared/evaluation/attempt.go:49:	IsSubmitted		100.0%
github.com/EduGoGroup/edugo-shared/evaluation/question.go:39:	Validate		100.0%
github.com/EduGoGroup/edugo-shared/evaluation/question.go:53:	GetCorrectOptions	100.0%
total:								(statements)		100.0%
```

| testing | 59.0% | 🟡 Aceptable | Media |

### testing (59.0%)

```
github.com/EduGoGroup/edugo-shared/testing/containers/helpers.go:20:	ExecSQLFile		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/helpers.go:46:	WaitForHealthy		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/helpers.go:72:	RetryOperation		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:35:	GetManager		48.7%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:102:	PostgreSQL		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:108:	MongoDB			0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:114:	RabbitMQ		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:120:	Cleanup			100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:125:	cleanup			42.3%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:174:	CleanPostgreSQL		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:183:	CleanMongoDB		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/manager.go:192:	PurgeRabbitMQ		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:26:	createMongoDB		59.1%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:87:	ConnectionString	100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:92:	Client			100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:97:	Database		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:102:	DropAllCollections	75.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:122:	DropCollections		80.0%
github.com/EduGoGroup/edugo-shared/testing/containers/mongodb.go:135:	Terminate		80.0%
github.com/EduGoGroup/edugo-shared/testing/containers/options.go:57:	NewConfig		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/options.go:65:	WithPostgreSQL		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/options.go:92:	WithMongoDB		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/options.go:111:	WithRabbitMQ		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/options.go:131:	Build			100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:27:	createPostgres		55.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:82:	ConnectionString	100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:87:	DB			100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:92:	ExecScript		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:97:	Truncate		70.6%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:134:	Terminate		80.0%
github.com/EduGoGroup/edugo-shared/testing/containers/postgres.go:145:	connectWithRetry	62.5%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:25:	createRabbitMQ		70.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:75:	ConnectionString	100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:80:	Connection		100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:85:	Channel			100.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:91:	PurgeAll		0.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:107:	PurgeQueue		75.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:123:	DeleteQueue		75.0%
github.com/EduGoGroup/edugo-shared/testing/containers/rabbitmq.go:139:	Terminate		80.0%
total:									(statements)		59.0%
```

| mongodb | 54.5% | 🟡 Aceptable | Media |

### mongodb (54.5%)

```
github.com/EduGoGroup/edugo-shared/database/mongodb/config.go:36:	DefaultConfig	100.0%
github.com/EduGoGroup/edugo-shared/database/mongodb/connection.go:21:	Connect		70.0%
github.com/EduGoGroup/edugo-shared/database/mongodb/connection.go:49:	GetDatabase	0.0%
github.com/EduGoGroup/edugo-shared/database/mongodb/connection.go:54:	HealthCheck	80.0%
github.com/EduGoGroup/edugo-shared/database/mongodb/connection.go:66:	Close		0.0%
total:									(statements)	54.5%
```

| postgres | 58.8% | 🟡 Aceptable | Media |

### postgres (58.8%)

```
github.com/EduGoGroup/edugo-shared/database/postgres/config.go:54:	DefaultConfig			100.0%
github.com/EduGoGroup/edugo-shared/database/postgres/connection.go:18:	Connect				0.0%
github.com/EduGoGroup/edugo-shared/database/postgres/connection.go:55:	HealthCheck			80.0%
github.com/EduGoGroup/edugo-shared/database/postgres/connection.go:67:	GetStats			100.0%
github.com/EduGoGroup/edugo-shared/database/postgres/connection.go:72:	Close				66.7%
github.com/EduGoGroup/edugo-shared/database/postgres/transaction.go:14:	WithTransaction			78.6%
github.com/EduGoGroup/edugo-shared/database/postgres/transaction.go:47:	WithTransactionIsolation	78.6%
total:									(statements)			58.8%
```

| gin | 98.5% | ✅ Excelente | Baja |

### gin (98.5%)

```
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:21:	GetUserID		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:37:	MustGetUserID		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:47:	GetEmail		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:63:	MustGetEmail		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:73:	GetRole			85.7%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:89:	MustGetRole		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:99:	GetClaims		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/context.go:115:	MustGetClaims		100.0%
github.com/EduGoGroup/edugo-shared/middleware/gin/jwt_auth.go:22:	JWTAuthMiddleware	100.0%
total:									(statements)		98.5%
```

| rabbit | 2.9% | 🔴 Crítico | Crítica |

### rabbit (2.9%)

```
github.com/EduGoGroup/edugo-shared/messaging/rabbit/config.go:57:		DefaultConfig		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:18:		Connect			0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:38:		GetChannel		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:43:		GetConnection		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:48:		Close			0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:63:		IsClosed		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:68:		DeclareExchange		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:81:		DeclareQueue		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:93:		BindQueue		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:104:		SetPrefetchCount	0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/connection.go:113:		HealthCheck		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer.go:26:		NewConsumer		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer.go:34:		Consume			0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer.go:81:		Close			0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer.go:87:		UnmarshalMessage	0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer.go:95:		HandleWithUnmarshal	0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer_dlq.go:12:		ConsumeWithDLQ		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer_dlq.go:124:	setupDLQ		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer_dlq.go:166:	sendToDLQ		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer_dlq.go:192:	getRetryCount		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/consumer_dlq.go:211:	cloneHeaders		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/dlq.go:18:			DefaultDLQConfig	100.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/dlq.go:30:			CalculateBackoff	100.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/publisher.go:25:		NewPublisher		0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/publisher.go:32:		Publish			0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/publisher.go:37:		PublishWithPriority	0.0%
github.com/EduGoGroup/edugo-shared/messaging/rabbit/publisher.go:71:		Close			0.0%
total:										(statements)		2.9%
```

