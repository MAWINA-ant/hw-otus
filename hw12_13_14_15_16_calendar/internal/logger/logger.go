package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus.Logger
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
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Errorf("couldn't use log file %s", logfile)
		} else {
			multiWriter := io.MultiWriter(os.Stdout, file)
			log.SetOutput(multiWriter)
			log.SetFormatter(&logrus.JSONFormatter{})
		}
		defer file.Close()
	}
	return &Logger{}
}
