package gologger

import (
	"github.com/sirupsen/logrus"
)

type CustomJsonFormatter struct {
	logrus.JSONFormatter
}

func (f *CustomJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b, err := f.JSONFormatter.Format(entry)
	switch entry.Level {
	case logrus.WarnLevel:
		return []byte("\033[33m" + string(b) + "\033[0m"), err
	case logrus.ErrorLevel:
		return []byte("\033[31m" + string(b) + "\033[0m"), err
	default:
		return b, err
	}

}
