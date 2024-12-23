package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service Service
	Cache   Cache
}

type Service struct {
	Port string `env:"SCHOOL_SERVICE_PORT"`
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
