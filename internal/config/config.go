package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Port   int         `env:"HTTP_SERVER_PORT" env-required:"true"`
	Alias  AliasConfig `yaml:"alias" env-required:"true"`
	DBType string      `yaml:"db_type" env-required:"true"`
}

type AliasConfig struct {
	Length            int    `yaml:"length" env-required:"true"`
	GenerationRetries int    `yaml:"generation_retries" env-required:"true"`
	Charset           string `yaml:"charset" env-required:"true"`
}

type PostgresConfig struct {
	MigrationsPath string `env:"MIGRATIONS_PATH" env-required:"true"`
	Host           string `env:"POSTGRES_HOST" env-required:"true"`
	Port           string `env:"POSTGRES_PORT" env-required:"true"`
	Name           string `env:"POSTGRES_DB" env-required:"true"`
	Username       string `env:"POSTGRES_USER" env-required:"true"`
	Password       string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

func LoadConfig(cfg interface{}) {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	err = cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
}
