package bot

import (
	"expense_accounting_bot/internal/utils/logger"
	"fmt"
	"github.com/tucnak/telebot"
)

// ExpenseBot структура для бота с телеграмом
type ExpenseBot struct {
	bot *telebot.Bot
}

// NewExpenseBot создает нового ExpenseBot
func NewExpenseBot(bot *telebot.Bot) *ExpenseBot {
	return &ExpenseBot{bot: bot}
}

// Start запускает обработку сообщений
func (e *ExpenseBot) Start() {
	// Создаем клавиатуру с кнопками
	menu := &telebot.ReplyMarkup{ResizeReplyKeyboard: true}

	// Обработчик команды /start
	e.bot.Handle("/start", func(m *telebot.Message) {
		logger.L.Info(fmt.Sprintf("Команда /start от пользователя %s", m.Sender.Username))
		msg := fmt.Sprintf(MessagesList.Welcome, m.Sender.LastName, m.Sender.FirstName)
		_, err := e.bot.Send(m.Sender, msg)
		if err != nil {
			logger.L.Error("Не удалось отправить приветственное сообщение:", err)
		}

		_, err = e.bot.Send(m.Sender, MessagesList.SelectAction, menu)
		if err != nil {
			logger.L.Error("Не удалось отправить приветственное сообщение:", err)
		}
	})

	createButtonsMainMenu(e, menu)

	// Запуск бота
	e.bot.Start()
}

func createButtonsMainMenu(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	// Создаем кнопки
	btnNewExpense := telebot.InlineButton{
		Unique: "btn_schedule",
		Text:   BtnTitlesList.BtnNewExpense,
	}
	btnMyExpenses := telebot.InlineButton{
		Unique: "btn_services",
		Text:   BtnTitlesList.BtnMyExpenses,
	}

	// Обработчики для кнопок
	e.bot.Handle(&btnNewExpense, btnNewExpenseFunc(e, menu))
	e.bot.Handle(&btnMyExpenses, btnMyExpensesFunc(e))

	menu.InlineKeyboard = [][]telebot.InlineButton{
		{btnNewExpense},
		{btnMyExpenses},
	}
}

// Обработчик нажатия кнопки "Добавить расход"
func btnNewExpenseFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(m *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnNewExpense, m.Sender.Username))

		menu.InlineKeyboard = nil

		row := make([]telebot.InlineButton, 0, 2)
		addBtnOfCategory(&row, "btn_groceries", BtnTitlesList.BtnGroceries)
		addBtnOfCategory(&row, "btn_beauty", BtnTitlesList.BtnBeauty)
		menu.InlineKeyboard = append(menu.InlineKeyboard, row)

		row = make([]telebot.InlineButton, 0, 2)
		addBtnOfCategory(&row, "btn_health", BtnTitlesList.BtnHealth)
		addBtnOfCategory(&row, "btn_restaurants", BtnTitlesList.BtnRestaurants)
		menu.InlineKeyboard = append(menu.InlineKeyboard, row)

		row = make([]telebot.InlineButton, 0, 2)
		addBtnOfCategory(&row, "btn_entertainment", BtnTitlesList.BtnEntertainment)
		addBtnOfCategory(&row, "btn_growth", BtnTitlesList.BtnGrowth)
		menu.InlineKeyboard = append(menu.InlineKeyboard, row)

		row = make([]telebot.InlineButton, 0, 2)
		addBtnOfCategory(&row, "btn_trips", BtnTitlesList.BtnTrips)
		addBtnOfCategory(&row, "btn_other", BtnTitlesList.BtnOther)
		menu.InlineKeyboard = append(menu.InlineKeyboard, row)

		btnBack := telebot.InlineButton{
			Unique: "btn_back",
			Text:   BtnTitlesList.BtnBack,
		}
		e.bot.Handle(&btnBack, btnBackFunc(e, menu))
		menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})

		_, err := e.bot.Edit(m.Message, MessagesList.Category, menu)
		if err != nil {
			logger.L.Error("Не удалось изменить сообщение сообщение:", err)
		}
	}
}

func addBtnOfCategory(row *[]telebot.InlineButton, unique string, text string) {
	btnGroceries := telebot.InlineButton{
		Unique: unique,
		Text:   text,
	}
	*row = append(*row, btnGroceries)
}

// Обработчик нажатия кнопки "Мои расходы"
func btnMyExpensesFunc(e *ExpenseBot) func(*telebot.Callback) {
	return func(m *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnMyExpenses, m.Sender.Username))
		e.viewExpenses(m)
	}
}

// viewExpenses обрабатывает просмотр расходов
func (e *ExpenseBot) viewExpenses(m *telebot.Callback) {
	logger.L.Info(fmt.Sprintf("Команда /expenses от пользователя %s", m.Sender.Username))
	// TODO: Загрузка данных из базы и форматирование вывода
	e.bot.Send(m.Sender, "Здесь будут отображены расходы (пока не реализовано)")
}

func btnBackFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(m *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnBack, m.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsMainMenu(e, menu)

		e.bot.Edit(m.Message, MessagesList.SelectAction, menu)
	}
}
