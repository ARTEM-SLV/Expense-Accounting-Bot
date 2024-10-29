package main

import (
	"expense_accounting_bot/pkg/bot"
	"log"
	"time"

	"expense_accounting_bot/config"
	"expense_accounting_bot/internal/utils/logger"
	"github.com/tucnak/telebot"
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
	expenseBot := bot.NewExpenseBot(b)

	// Запускаем бота
	logger.L.Info("Запуск бота...")
	expenseBot.Start()
}
