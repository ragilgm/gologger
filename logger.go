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

type FileWriter struct {
	FileName  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
}

type LoggerConfig struct {
	ServiceName    string
	ServiceVersion string
	lj             io.Writer
	FileWriter     *FileWriter
}

var once sync.Once
var log *logrus.Logger
var fields *logrus.Fields

func InitLogger(param *LoggerConfig) {
	if param == nil {
		Error(context.Background(), "Invalid logger config")
		return
	}
	once.Do(func() {
		// Set up logrus logger
		log = logrus.New()

		fields = &logrus.Fields{
			serviceName:    param.ServiceName,
			serviceVersion: param.ServiceVersion,
		}

		// env for define output
		param.initWriter()

		// format default json
		formatter := &logrus.JSONFormatter{}
		formatter.TimestampFormat = "2006-01-02T15:04:05.000000Z07:00"
		log.SetFormatter(formatter)

	})

	// Add service name and version to log fields
	log.WithFields(*fields)

}

func (l *LoggerConfig) initFileWriter() {
	w := l.FileWriter
	l.lj = &lumberjack.Logger{
		Filename:   w.FileName,  // path to log file
		MaxSize:    w.MaxSize,   // megabytes
		MaxBackups: w.MaxBackup, // maximum number of backups to keep
		MaxAge:     w.MaxAge,    // days
		Compress:   w.Compress,  // compress old log files
	}
	// Set up file hook for log rotation
	fileHook := &writer.Hook{
		Writer:    l.lj,
		LogLevels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	}

	// Add file hook to logger
	log.AddHook(fileHook)
}

func (l *LoggerConfig) initWriter() {
	if l.FileWriter != nil {
		l.consoleAndWriteFile()
		return
	}

	l.consolOnly()
}

func (l *LoggerConfig) consolOnly() {
	logrus.SetOutput(os.Stdout)
	return
}

func (l *LoggerConfig) consoleAndWriteFile() {
	l.initFileWriter()
	logrus.SetOutput(io.MultiWriter(os.Stdout, l.lj))
	return
}

func setLoggerContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func getLoggerEntry(ctx context.Context) *logrus.Entry {

	if ctx == nil {
		ctx = context.Background()
	}

	logger := ctx.Value(loggerCtxKey)
	if logger != nil {
		return logger.(*logrus.Entry)
	}

	if log != nil {
		entry := log.WithFields(*fields)
		setLoggerContext(ctx, entry)
		return entry
	}

	loggerModel := &LoggerConfig{
		ServiceName:    serviceUnknown,
		ServiceVersion: serviceUnknown,
	}

	InitLogger(loggerModel)

	fields = &logrus.Fields{
		serviceName:    loggerModel.ServiceName,
		serviceVersion: loggerModel.ServiceVersion,
	}

	return log.WithFields(*fields)
}

func WithCustomPayload(ctx context.Context, title string, payloads map[string]any) {
	logger := getLoggerEntry(ctx)
	field := make(logrus.Fields)
	for key, value := range payloads {
		field[key] = value
	}
	logger = logger.WithFields(field)
	logger.Info(title)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	logger := getLoggerEntry(ctx)
	return setLoggerContext(ctx, logger.WithField(traceIDKey, traceID))
}

func WithFunctionName(ctx context.Context, functionName string) context.Context {
	logger := getLoggerEntry(ctx)
	return setLoggerContext(ctx, logger.WithField(functionKey, functionName))
}

func Info(ctx context.Context, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Info(message)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Infof(format, message)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Errorf(format, message)
}

func Error(ctx context.Context, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Error(message)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Warnf(format, message)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger := getLoggerEntry(ctx)
	message := fmt.Sprint(args...)
	logger.Warn(message)
}
