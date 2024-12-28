package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyLogger = key("logger")

type Config struct {
	Service  Service
	Cache    Cache
	Logger   Logger
	Platform Platform
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

type Service struct {
	Port string `env:"SCHOOL_SERVICE_PORT"`
	Name string `env:"SCHOOL_SERVICE_NAME"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type Cache struct {
	Host string `env:"SCHOOL_SERVICE_REDIS_HOST"`
	Port string `env:"SCHOOL_SERVICE_REDIS_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
