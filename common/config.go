package common

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	StorageConfig `yaml:"storage"`
	//StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type StorageConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"4s"`
	Timeout           time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"60s"`
	//User        string        `yaml:"user" env-required:"true"`
	//Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func (sc *StorageConfig) URL() string {
	//url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", storage.Username, storage.Password, storage.Host, storage.Port, storage.Database)

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		sc.Driver,
		sc.Username,
		sc.Password,
		sc.Host,
		sc.Port,
		sc.Database,
	)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
