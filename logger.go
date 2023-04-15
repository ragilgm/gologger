package gologger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

const serviceName = "service_name"
const serviceVersion = "service_version"
const serviceUnknown = "unknown"

var (
	loggerCtxKey = "loggerCtx"
	traceIDKey   = "trace_id"
	functionKey  = "function"
)

var (
	EnvLocalKey = "local"
	EnvProdKey  = "prod"
	EnvStgKey   = "stg"
)

type LoggerConfig struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	lj             io.Writer
	FileName       string
	MaxSize        int
	MaxBackup      int
	MaxAge         int
	Compress       bool
}

var once sync.Once
var log *logrus.Logger

type RemoveMsgFieldHook struct{}

func InitLogger(param *LoggerConfig) {
	if param == nil {
		Error(context.Background(), "Invalid logger config")
		return
	}
	once.Do(func() {
		// Set up logrus logger
		log = logrus.New()

		fields := logrus.Fields{
			serviceName:    param.ServiceName,
			serviceVersion: param.ServiceVersion,
		}

		// Add service name and version to log fields
		log.WithFields(fields)

		// env for define output
		param.selectEnv(param.Env)

		// Set formatter to JSON format
		log.SetFormatter(&logrus.JSONFormatter{})

	})

}

func (l *LoggerConfig) initFileHook() {
	l.lj = &lumberjack.Logger{
		Filename:   l.FileName,  // path to log file
		MaxSize:    l.MaxSize,   // megabytes
		MaxBackups: l.MaxBackup, // maximum number of backups to keep
		MaxAge:     l.MaxAge,    // days
		Compress:   l.Compress,  // compress old log files
	}
	// Set up file hook for log rotation
	fileHook := &writer.Hook{
		Writer:    l.lj,
		LogLevels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	}

	// Add file hook to logger
	log.AddHook(fileHook)
}

func (l *LoggerConfig) selectEnv(env string) {
	switch env {
	case EnvLocalKey:
		logrus.SetOutput(os.Stdout)
		return
	case EnvStgKey:
		l.initFileHook()
		logrus.SetOutput(io.MultiWriter(os.Stdout, l.lj))
		return
	case EnvProdKey:
		l.initFileHook()
		logrus.SetOutput(l.lj)
		return
	default:
		l.initFileHook()
		logrus.SetOutput(io.MultiWriter(os.Stdout, l.lj))
		return
	}
}

func SetLoggerCtx(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetLoggerEntry(ctx context.Context) *logrus.Entry {

	if ctx == nil {
		ctx = context.Background()
	}

	logger := ctx.Value(loggerCtxKey)
	if logger != nil {
		return logger.(*logrus.Entry)
	}

	loggerModel := &LoggerConfig{
		ServiceName:    serviceUnknown,
		ServiceVersion: serviceUnknown,
		MaxSize:        10,
		MaxBackup:      10,
		MaxAge:         28,
		FileName:       "/var/log/default.log",
		Compress:       true,
	}

	InitLogger(loggerModel)

	fields := logrus.Fields{
		serviceName:    loggerModel.ServiceName,
		serviceVersion: loggerModel.ServiceVersion,
	}

	return log.WithFields(fields)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	logger := GetLoggerEntry(ctx)
	return SetLoggerCtx(ctx, logger.WithField(traceIDKey, traceID))
}

func WithFunctionName(ctx context.Context, functionName string) context.Context {
	logger := GetLoggerEntry(ctx)
	return SetLoggerCtx(ctx, logger.WithField(functionKey, functionName))
}

func Info(ctx context.Context, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Info(message)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Infof(format, message)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Errorf(format, message)
}

func Error(ctx context.Context, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Error(message)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Warnf(format, message)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger := GetLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Warn(message)
}
