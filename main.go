package main

import "context"

func main() {
	InitLogger(&LoggerConfig{
		ServiceName:    "service_1",
		ServiceVersion: "1.0.0",
		FileName:       "/var/log/service_1/output.log",
		Compress:       true,
		MaxAge:         28,
		MaxBackup:      10,
		MaxSize:        1,
		Env:            EnvStgKey,
	})

	ctx := context.Background()

	ctx = WithFunctionName(ctx, "test")
	ctx = WithTraceID(ctx, "test")

	Info(ctx, "hello world")

}
