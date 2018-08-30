package logger

import "github.com/sirupsen/logrus"

func DebugLogger() *logrus.Logger {
	var logger = logrus.StandardLogger()
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

func StdLogger() *logrus.Logger {
	var logger = logrus.StandardLogger()
	logger.SetLevel(logrus.ErrorLevel)
	return logger
}

type Logger interface {
	logrus.FieldLogger
}
