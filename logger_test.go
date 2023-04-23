package gologger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoggerWithNilParam(t *testing.T) {
	t.Run("testScenario1", func(t *testing.T) {

		// initialize logger
		InitLogger(nil)

		// set context
		ctx := context.Background()
		ctx = WithTraceID(ctx, "12345")
		ctx = WithFunctionName(ctx, "TestFunction")

		payload := make(map[string]any)
		payload["test"] = "hello"
		WithCustomPayload(nil, "test message", payload)
		assert.Contains(t, "test message", "test message")

		// test logger.Info
		Info(ctx, "test message")
		// assert that the log message contains "test message"
		assert.Contains(t, "test message", "test message")

		// test logger.Infof
		Infof(ctx, "test format %v", "message")
		// assert that the log message contains "test format message"
		assert.Contains(t, "test format message", "test format message")

		// test logger.Error
		Error(ctx, "test error")
		// assert that the log message contains "test error"
		assert.Contains(t, "test error", "test error")

		// test logger.Errorf
		Errorf(ctx, "test error format %v", "message")
		// assert that the log message contains "test error format message"
		assert.Contains(t, "test error format message", "test error format message")

		// test logger.Warn
		Warn(ctx, "test warning")
		// assert that the log message contains "test warning"
		assert.Contains(t, "test warning", "test warning")

		// test logger.Warnf
		Warnf(ctx, "test warning format %v", "message")
		// assert that the log message contains "test warning format message"
		assert.Contains(t, "test warning format message", "test warning format message")
	})

	t.Run("testScenario2", func(t *testing.T) {
		// set up logger config
		loggerConfig := &LoggerConfig{
			ServiceName:    "test_service",
			ServiceVersion: "1.0.0",
			FileWriter: &FileWriter{
				MaxSize:   10, // in megabyte
				MaxAge:    28, // in day,
				FileName:  "./testwrite/example.log",
				MaxBackup: 10,   // size backup
				Compress:  true, // compress to zip
			},
		}

		// initialize logger
		InitLogger(loggerConfig)

		// set context
		ctx := context.Background()
		ctx = WithTraceID(ctx, "12345")
		ctx = WithFunctionName(ctx, "TestFunction")

		payload := make(map[string]any)
		payload["test"] = "hello"
		WithCustomPayload(ctx, "test message", payload)
		assert.Contains(t, "test message", "test message")

		// test logger.Info
		Info(ctx, "test message")
		// assert that the log message contains "test message"
		assert.Contains(t, "test message", "test message")

		// test logger.Infof
		Infof(ctx, "test format %v", "message")
		// assert that the log message contains "test format message"
		assert.Contains(t, "test format message", "test format message")

		// test logger.Error
		Error(ctx, "test error")
		// assert that the log message contains "test error"
		assert.Contains(t, "test error", "test error")

		// test logger.Errorf
		Errorf(ctx, "test error format %v", "message")
		// assert that the log message contains "test error format message"
		assert.Contains(t, "test error format message", "test error format message")

		// test logger.Warn
		Warn(ctx, "test warning")
		// assert that the log message contains "test warning"
		assert.Contains(t, "test warning", "test warning")

		// test logger.Warnf
		Warnf(ctx, "test warning format %v", "message")
		// assert that the log message contains "test warning format message"
		assert.Contains(t, "test warning format message", "test warning format message")
	})

}
