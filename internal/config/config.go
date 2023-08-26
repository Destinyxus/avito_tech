package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file in the path %s does not exist", path)
	}

	cfg := Config{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("can't read cfg: %s", err)
	}

	return &cfg

}
