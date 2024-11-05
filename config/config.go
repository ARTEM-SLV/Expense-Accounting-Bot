package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации
type Config struct {
	TelegramToken string
	DatabasePath  string
	AdminID       string
}

// LoadConfig загружает конфигурацию из .env файла и переменных окружения
func LoadConfig() *Config {
	// Загружаем переменные из .env файла
	env := os.Getenv("RAILWAY_ENVIRONMENT")
	if env == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Не удалось загрузить .env файл, будут использованы переменные окружения", err)
		}
	}

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DatabasePath:  os.Getenv("DATABASE_PATH"),
		AdminID:       os.Getenv("ADMIN_ID"),
	}

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN не установлен")
	}
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = "expenses.db" // Устанавливаем значение по умолчанию
	}

	return cfg
}
