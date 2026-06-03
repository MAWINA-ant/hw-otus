package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrusLogger *logrus.Logger
}

func New(level, logfile string) *Logger {
	log := logrus.New()
	switch strings.ToLower(level) {
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	}
	return &Logger{logrusLogger: log}
}

func (l Logger) Error(msg string) {
	l.logrusLogger.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.logrusLogger.Warn(msg)
}

func (l Logger) Info(msg string) {
	l.logrusLogger.Info(msg)
}

func (l Logger) Debug(msg string) {
	l.logrusLogger.Debug(msg)
}
