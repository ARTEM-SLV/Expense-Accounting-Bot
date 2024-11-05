package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/tucnak/telebot"

	"expense_accounting_bot/config"
	"expense_accounting_bot/internal/utils/logger"
	"expense_accounting_bot/pkg/bot"
	"expense_accounting_bot/pkg/bot/telegram"
	"expense_accounting_bot/pkg/repository"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Инициализируем логгер
	logPath := "bot.log" // Путь к файлу для логов
	err := logger.InitLogger(logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.L.Close()

	// Подключаемся к SQLite
	db, err := sql.Open("sqlite", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализируем репозиторий
	repo := repository.NewSQLiteExpenseRepository(db)
	if err = repo.InitSchema(); err != nil {
		log.Fatalf("Ошибка при инициализации схемы: %v", err)
	}

	// Инициализация бота
	err = bot.InitStringValues()
	if err != nil {
		log.Fatal("Ошибка при инициализации бота:", err)
	}

	// Инициализация бота
	b, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("Ошибка при создании бота: ", err)
	}

	// Создаем объект нашего бота с логгером
	expenseBot := telegram.NewExpenseBot(b, repo)

	// Запускаем бота
	logger.L.Info("Запуск бота...")
	expenseBot.Start()
}
