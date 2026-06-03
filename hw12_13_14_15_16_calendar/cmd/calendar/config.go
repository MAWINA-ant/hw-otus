package main

import (
	"os"

	"go.yaml.in/yaml/v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	Server  ServerConf  `yaml:"server"`
}

type LoggerConf struct {
	Level   string `yaml:"level"`
	LogFile string `yaml:"file"`
}

type StorageConf struct {
	InMemory bool   `yaml:"in-memory"`
	Host     string `yaml:"host"`
	DataBase string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) ParseConfigFromFile(path string) error {
	var data []byte
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
