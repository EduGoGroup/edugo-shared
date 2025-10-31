package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger es la implementación de Logger usando Zap
type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger crea un nuevo logger usando Zap
// level: "debug", "info", "warn", "error", "fatal"
// format: "json" o "console"
func NewZapLogger(level, format string) Logger {
	// Configurar nivel de log
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Configurar encoder (formato)
	var encoder zapcore.Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// Console format con colores
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Core de Zap
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Crear logger con opciones
	zapLog := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // Skip para mostrar el caller correcto
	)

	return &zapLogger{
		logger: zapLog.Sugar(),
	}
}

// Debug registra un mensaje de nivel debug
func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debugw(msg, fields...)
}

// Info registra un mensaje de nivel info
func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.logger.Infow(msg, fields...)
}

// Warn registra un mensaje de nivel warning
func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warnw(msg, fields...)
}

// Error registra un mensaje de nivel error
func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.logger.Errorw(msg, fields...)
}

// Fatal registra un mensaje de nivel fatal y termina la aplicación
func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatalw(msg, fields...)
}

// With agrega campos contextuales al logger
func (l *zapLogger) With(fields ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// Sync sincroniza el buffer del logger
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
