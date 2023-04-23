package main

import (
	"github.com/ragilgm/gologger"
)

type Payload struct {
	Message string
}

func main() {

	logConfig := &gologger.LoggerConfig{
		ServiceName:    "test_service",
		ServiceVersion: "1.0.0",
		FileWriter: &gologger.FileWriter{
			MaxSize:   10, // in megabyte
			MaxAge:    28, // in day,
			FileName:  "/",
			MaxBackup: 10,   // size backup
			Compress:  true, // compress to zip
		},
	}

	gologger.InitLogger(logConfig)

	gologger.Info(nil, "hello world")
}
