package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus.Logger
}

func New(level, logfile string) *Logger {
	log := &Logger{}
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
		correctFilePath := correctPath(logfile)
		dir := filepath.Dir(correctFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Errorf("couldn't create directory %s, %s", dir, err)
			}
		}
		file, err := os.OpenFile(correctFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Errorf("couldn't use log file %s, %s", logfile, err)
		} else {
			multiWriter := io.MultiWriter(os.Stdout, file)
			log.SetOutput(multiWriter)
			log.SetFormatter(&logrus.JSONFormatter{})
		}
		defer file.Close()
	}
	return log
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
