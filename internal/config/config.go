package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	HeaderInc string = "client"
	HeaderDec string = "admin"
	HeaderKey string = "User-Role"
)

type Config struct {
	Env      string `env:"ENV"` // env-default:"local"`
	HTTP     HTTPConfig
	DBConfig DBConfig
}

type DBConfig struct {
	User     string `env:"DB_USER" env-default:"user"`
	Password string `env:"DB_PASSWORD" env-default:"password"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Dbname   string `env:"DB_DBNAME" env-default:"postgres"`
	Sslmode  string `env:"DB_SSLMODE" env-default:"disable"`
}

type HTTPConfig struct {
	URLPort string `env:"HTTP_PORT" env-default:"3000"`
}

// загрузка конфига из .env
func MustLoad() *Config {

	LoadEnv(".env")

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func LoadEnv(envFile string) {
	err := godotenv.Load(dir(envFile))
	if err != nil {
		panic("No .env file found " + err.Error())
	}
}

func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
