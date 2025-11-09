package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

type LogContext struct {
	ProjectID  *uuid.UUID
	DatabaseID *uuid.UUID
	TableID    *uuid.UUID
	WorkflowID *uuid.UUID
	UserID     *uuid.UUID
	RequestID  string
	Operation  string
}

type Logger struct {
	logger     zerolog.Logger
	storage    LogStorage
	bufferSize int
	verbosity  LogVerbosity
}

func NewLogger(output io.Writer, storage LogStorage, verbosity LogVerbosity) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{
		logger:     logger,
		storage:    storage,
		bufferSize: 100,
		verbosity:  verbosity,
	}
}

func NewProductionLogger(storage LogStorage, verbosity LogVerbosity) *Logger {
	return NewLogger(os.Stdout, storage, verbosity)
}

func NewDevelopmentLogger(storage LogStorage, verbosity LogVerbosity) *Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			if level, ok := i.(string); ok {
				switch level {
				case "debug":
					return "\033[35m●\033[0m"
				case "info":
					return "\033[36m●\033[0m"
				case "warn":
					return "\033[33m●\033[0m"
				case "error":
					return "\033[31m●\033[0m"
				case "fatal":
					return "\033[31;1m●\033[0m"
				default:
					return level
				}
			}
			return ""
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("\033[1m%s\033[0m", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("\033[2m%s\033[0m", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%v", i)
		},
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
	}
	return NewLogger(consoleWriter, storage, verbosity)
}

func (l *Logger) WithContext(ctx LogContext) *ContextLogger {
	event := l.logger.With()

	if ctx.ProjectID != nil {
		event = event.Str("project_id", ctx.ProjectID.String())
	}
	if ctx.DatabaseID != nil {
		event = event.Str("database_id", ctx.DatabaseID.String())
	}
	if ctx.TableID != nil {
		event = event.Str("table_id", ctx.TableID.String())
	}
	if ctx.WorkflowID != nil {
		event = event.Str("workflow_id", ctx.WorkflowID.String())
	}
	if ctx.UserID != nil {
		event = event.Str("user_id", ctx.UserID.String())
	}
	if ctx.RequestID != "" {
		event = event.Str("request_id", ctx.RequestID)
	}
	if ctx.Operation != "" {
		event = event.Str("operation", ctx.Operation)
	}

	return &ContextLogger{
		logger:    event.Logger(),
		context:   ctx,
		storage:   l.storage,
		verbosity: l.verbosity,
	}
}

func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	l.addFields(event, fields...)
	event.Msg(msg)

	l.storeAndStream(LevelInfo, msg, nil, fields...)
}

func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	l.addFields(event, fields...)
	event.Msg(msg)

	l.storeAndStream(LevelWarn, msg, nil, fields...)
}

func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	l.addFields(event, fields...)
	event.Msg(msg)

	l.storeAndStream(LevelError, msg, nil, fields...)
}

func (l *Logger) ErrorWithStack(msg string, err error, stackTrace string, fields ...map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	if stackTrace != "" {
		event = event.Str("stack", stackTrace)
	}
	l.addFields(event, fields...)
	event.Msg(msg)

	allFields := make(map[string]interface{})
	if stackTrace != "" {
		allFields["stack"] = stackTrace
	}
	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			allFields[key] = value
		}
	}
	l.storeAndStream(LevelError, msg, nil, allFields)
}

func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *Logger) addFields(event *zerolog.Event, fields ...map[string]interface{}) {
	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			event.Interface(key, value)
		}
	}
}

func (l *Logger) storeAndStream(level LogLevel, msg string, ctx *LogContext, fields ...map[string]interface{}) {
	entry := &LogEntry{
		ID:        uuid.New(),
		Level:     level,
		Message:   msg,
		Timestamp: time.Now(),
		Fields:    make(map[string]interface{}),
	}

	if ctx != nil {
		entry.ProjectID = ctx.ProjectID
		entry.DatabaseID = ctx.DatabaseID
		entry.TableID = ctx.TableID
		entry.WorkflowID = ctx.WorkflowID
		entry.UserID = ctx.UserID
		entry.RequestID = ctx.RequestID
		entry.Operation = ctx.Operation
	}

	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			entry.Fields[key] = value
		}
	}

	if l.storage != nil {
		go func() {
			if err := l.storage.Store(context.Background(), entry); err != nil {
				fmt.Printf("Failed to store log entry: %v\n", err)
			}
		}()
	}
}

type ContextLogger struct {
	logger    zerolog.Logger
	context   LogContext
	storage   LogStorage
	verbosity LogVerbosity
}

const LoggerKey = "logger"

func GetLogger(ctx context.Context) *ContextLogger {
	if logger, ok := ctx.Value(LoggerKey).(*ContextLogger); ok {
		return logger
	}
	return nil
}

func (cl *ContextLogger) GetVerbosity() LogVerbosity {
	return cl.verbosity
}

func (cl *ContextLogger) Debug(msg string, fields ...map[string]interface{}) {
	event := cl.logger.Debug()
	cl.addFields(event, fields...)
	event.Msg(msg)
}

func (cl *ContextLogger) Info(msg string, fields ...map[string]interface{}) {
	event := cl.logger.Info()
	cl.addFields(event, fields...)
	event.Msg(msg)

	cl.storeAndStream(LevelInfo, msg, fields...)
}

func (cl *ContextLogger) Warn(msg string, fields ...map[string]interface{}) {
	event := cl.logger.Warn()
	cl.addFields(event, fields...)
	event.Msg(msg)

	cl.storeAndStream(LevelWarn, msg, fields...)
}

func (cl *ContextLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := cl.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	cl.addFields(event, fields...)
	event.Msg(msg)

	cl.storeAndStream(LevelError, msg, fields...)
}

func (cl *ContextLogger) ErrorWithStack(msg string, err error, stackTrace string, fields ...map[string]interface{}) {
	event := cl.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	if stackTrace != "" {
		event = event.Str("stack", stackTrace)
	}
	cl.addFields(event, fields...)
	event.Msg(msg)

	allFields := make(map[string]interface{})
	if stackTrace != "" {
		allFields["stack"] = stackTrace
	}
	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			allFields[key] = value
		}
	}
	cl.storeAndStream(LevelError, msg, allFields)
}

func (cl *ContextLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := cl.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	cl.addFields(event, fields...)
	event.Msg(msg)
}

func (cl *ContextLogger) addFields(event *zerolog.Event, fields ...map[string]interface{}) {
	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			event.Interface(key, value)
		}
	}
}

func (cl *ContextLogger) storeAndStream(level LogLevel, msg string, fields ...map[string]interface{}) {
	entry := &LogEntry{
		ID:         uuid.New(),
		Level:      level,
		Message:    msg,
		Timestamp:  time.Now(),
		ProjectID:  cl.context.ProjectID,
		DatabaseID: cl.context.DatabaseID,
		TableID:    cl.context.TableID,
		WorkflowID: cl.context.WorkflowID,
		UserID:     cl.context.UserID,
		RequestID:  cl.context.RequestID,
		Operation:  cl.context.Operation,
		Fields:     make(map[string]interface{}),
	}

	for _, fieldMap := range fields {
		for key, value := range fieldMap {
			entry.Fields[key] = value
		}
	}

	if cl.storage != nil {
		go func() {
			if err := cl.storage.Store(context.Background(), entry); err != nil {
				fmt.Printf("Failed to store log entry: %v\n", err)
			}
		}()
	}
}
