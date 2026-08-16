package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log  *logrus.Logger
	file *os.File
}

func New(level, logfile string) *Logger {
	l := &Logger{
		log: logrus.New(),
	}
	switch strings.ToLower(level) {
	case "error":
		l.log.SetLevel(logrus.ErrorLevel)
	case "warn":
		l.log.SetLevel(logrus.WarnLevel)
	case "info":
		l.log.SetLevel(logrus.InfoLevel)
	case "trace":
		l.log.SetLevel(logrus.TraceLevel)
	}
	if logfile != "" {
		correctFilePath := correctPath(logfile)
		dir := filepath.Dir(correctFilePath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			l.log.Errorf("couldn't create directory %s, %s", dir, err)
		}
		file, err := os.OpenFile(correctFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			l.log.Errorf("couldn't use log file %s, %s", logfile, err)
		} else {
			multiWriter := io.MultiWriter(os.Stdout, file)
			l.log.SetOutput(multiWriter)
			l.log.SetFormatter(&logrus.JSONFormatter{})
		}
	}
	return l
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func correctPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return strings.Replace(path, "~", home, 1)
	} else if !filepath.IsAbs(path) {
		absPath, _ := filepath.Abs(path)
		return absPath
	}
	return path
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

/*
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args...)
}*/
