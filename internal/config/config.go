package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"main/internal/storage/postgres"
	"os"
	"time"
)

type Config struct {
	Env         string            `yaml:"env"`
	StoragePath string            `yaml:"storage_path"`
	TokenTTl    time.Duration     `yaml:"token_ttl"`
	Grpc        GrpcConfig        `yaml:"grpc"`
	Postgres    postgres.Postgres `yaml:"postgres"`
}

type GrpcConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	RetriesCount int           `yaml:"retries_count"`
	Timeout      time.Duration `yaml:"timeout"`
	Insecure     bool          `yaml:"insecure"`
}

func MustLoad() *Config {
	path := FetchConfigPath()
	if path == "" {
		panic("config path if empty")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path does not exists")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config")
	}
	return &cfg
}

func FetchConfigPath() (path string) {
	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
