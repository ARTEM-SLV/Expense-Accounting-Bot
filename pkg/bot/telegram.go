package bot

import (
	"fmt"
	"github.com/tucnak/telebot"

	"expense_accounting_bot/internal/utils/logger"
)

// ExpenseBot структура для бота с телеграмом
type ExpenseBot struct {
	bot    *telebot.Bot
	logger *logger.Logger
}

// NewExpenseBot создает нового ExpenseBot
func NewExpenseBot(bot *telebot.Bot, logger *logger.Logger) *ExpenseBot {
	return &ExpenseBot{bot: bot, logger: logger}
}

// Start запускает обработку сообщений
func (e *ExpenseBot) Start() {
	// Создаем клавиатуру с кнопками
	menu := &telebot.ReplyMarkup{ResizeReplyKeyboard: true}

	// Обработчик команды /start
	e.bot.Handle("/start", func(m *telebot.Message) {
		e.logger.Info(fmt.Sprintf("Команда /start от пользователя %s", m.Sender.Username))
		msg := fmt.Sprintf(MessagesList.Welcome, m.Sender.LastName, m.Sender.FirstName)
		_, err := e.bot.Send(m.Sender, msg, menu)
		if err != nil {
			e.logger.Error("Не удалось отправить приветственное сообщение:", err)
		}
	})

	// Создаем кнопки
	btnNewExpense := telebot.InlineButton{
		Unique: "btn_schedule",
		Text:   BtnTitlesList.BtnNewExpense,
	}
	btnMyExpenses := telebot.InlineButton{
		Unique: "btn_services",
		Text:   BtnTitlesList.BtnMyExpenses,
	}

	menu.InlineKeyboard = [][]telebot.InlineButton{
		{btnNewExpense},
		{btnMyExpenses},
	}

	// Обработчики для кнопок
	e.bot.Handle(&btnNewExpense, btnNewExpenseFunc)
	e.bot.Handle(&btnMyExpenses, btnSettingsFunc)

	// Запуск бота
	e.bot.Start()
}

// Обработчик нажатия кнопки "Добавить расход"
func btnNewExpenseFunc(e *ExpenseBot) func(*telebot.Message) {
	return func(m *telebot.Message) {
		e.logger.Info(fmt.Sprintf("Нажата кнопка 'Добавить расход' пользователем %s", m.Sender.Username))
		e.bot.Send(m.Sender, "Пожалуйста, отправьте сумму и категорию расхода в формате:\n<сумма> <категория>")
	}
}

// Обработчик нажатия кнопки "Мои расходы"
func btnSettingsFunc(e *ExpenseBot) func(*telebot.Message) {
	return func(m *telebot.Message) {
		e.logger.Info(fmt.Sprintf("Нажата кнопка 'Мои расходы' пользователем %s", m.Sender.Username))
		e.viewExpenses(m)
	}
}

// viewExpenses обрабатывает просмотр расходов
func (e *ExpenseBot) viewExpenses(m *telebot.Message) {
	e.logger.Info(fmt.Sprintf("Команда /expenses от пользователя %s", m.Sender.Username))
	// TODO: Загрузка данных из базы и форматирование вывода
	e.bot.Send(m.Sender, "Здесь будут отображены расходы (пока не реализовано)")
}
