package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Storage StorageConf `toml:"storage"`
	Server  ServerConf  `toml:"server"`
}

type LoggerConf struct {
	Level   string `toml:"level"`
	LogFile string `toml:"file"`
}

type StorageConf struct {
	InMemory bool   `toml:"in-memory"`
	Host     string `toml:"host"`
	DataBase string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type ServerConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
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
	err = toml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
