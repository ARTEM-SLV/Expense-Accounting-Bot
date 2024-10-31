package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tucnak/telebot"

	"expense_accounting_bot/internal/utils/logger"
	"expense_accounting_bot/pkg/repository"
)

// Словарь для хранения временного состояния пользователя
var userStates = make(map[int]string)

// Поле для хранения последнего сообщения от бота
var lastBotMessage *telebot.Message

// ExpenseBot структура для бота с телеграмом
type ExpenseBot struct {
	bot  *telebot.Bot
	repo repository.ExpenseRepository
}

// NewExpenseBot создает нового ExpenseBot
func NewExpenseBot(bot *telebot.Bot, repo repository.ExpenseRepository) *ExpenseBot {
	return &ExpenseBot{bot: bot, repo: repo}
}

// Start запускает обработку сообщений
func (e *ExpenseBot) Start() {
	// Создаем клавиатуру с кнопками
	menu := &telebot.ReplyMarkup{ResizeReplyKeyboard: true}

	// Обработчик команды /start
	e.bot.Handle("/start", func(m *telebot.Message) {
		logger.L.Info(fmt.Sprintf("Команда /start от пользователя %s", m.Sender.Username))

		userID := m.Sender.ID
		userName := m.Sender.Username

		// Проверяем, зарегистрирован ли пользователь
		isRegistered, dateReg, err := e.repo.IsUserRegistered(userID)
		if err != nil {
			logger.L.Error(MessagesList.ErrorReg, err)
			e.sendBotMessage(m, MessagesList.ErrorReg)
			return
		}

		if isRegistered {
			err = e.bot.Delete(lastBotMessage)
			if err != nil {
				logger.L.Error("Ошибка при удалении сообщения:", err)
			}
			e.sendBotMessage(m, fmt.Sprintf(MessagesList.UserRegistered, userName, dateReg))
			e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu)
			return
		}

		// Регистрируем нового пользователя
		err = e.repo.AddUser(userID, userName)
		if err != nil {
			logger.L.Error(MessagesList.ErrorReg, err)
			e.sendBotMessage(m, MessagesList.ErrorReg)
			return
		}

		msg := fmt.Sprintf(MessagesList.Welcome, m.Sender.FirstName, m.Sender.LastName)
		e.sendBotMessage(m, msg)

		e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu)
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
	e.bot.Handle(&btnMyExpenses, btnMyExpensesFunc(e, menu))

	menu.InlineKeyboard = [][]telebot.InlineButton{
		{btnNewExpense},
		{btnMyExpenses},
	}

	e.bot.Handle(telebot.OnText, e.handleOnText(menu))
}

// Обработчик нажатия кнопки "Добавить расход"
func btnNewExpenseFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnNewExpense, c.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsOfCategories(e, menu)

		e.editBotMessageWithMenu(c, MessagesList.SelectCategory, menu) // e.bot.Edit(c.Message, MessagesList.SelectCategory, menu)
	}
}

func createButtonsOfCategories(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	row := make([]telebot.InlineButton, 0, 2)
	for _, key := range Categories {
		value := BtnCategoriesList[key]
		addBtnOfCategory(&row, e, key, value, menu)

		if len(row) == 2 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, row)
			row = make([]telebot.InlineButton, 0, 2)
		}
	}

	btnBack := telebot.InlineButton{
		Unique: "btn_back",
		Text:   BtnTitlesList.BtnBack,
	}
	e.bot.Handle(&btnBack, btnBackFunc(e, menu, "MainMenu"))
	menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})
}

func addBtnOfCategory(row *[]telebot.InlineButton, e *ExpenseBot, unique string, category string, menu *telebot.ReplyMarkup) {
	newBtn := telebot.InlineButton{
		Unique: unique,
		Text:   category,
	}

	e.bot.Handle(&newBtn, btnCategoryFunc(e, category, menu))

	*row = append(*row, newBtn)
}

func btnCategoryFunc(e *ExpenseBot, category string, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		menu.InlineKeyboard = nil

		e.bot.Respond(c, &telebot.CallbackResponse{Text: fmt.Sprintf(MessagesList.Category, category)})

		userID := c.Sender.ID
		userStates[userID] = category

		btnBack := telebot.InlineButton{
			Unique: "btn_back",
			Text:   BtnTitlesList.BtnBack,
		}
		e.bot.Handle(&btnBack, btnBackFunc(e, menu, "SelectCategory"))
		menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})

		e.editBotMessageWithMenu(c, MessagesList.EnterAmount, menu)

		e.bot.Handle(telebot.OnText, addExpense(e, menu, c))
	}
}

func addExpense(e *ExpenseBot, menu *telebot.ReplyMarkup, c *telebot.Callback) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		category, ok := userStates[m.Sender.ID]
		if !ok || category == "" {
			return
		}

		amount, err := strconv.ParseFloat(m.Text, 64)
		if err != nil {
			e.sendBotMessage(m, MessagesList.NumberError)
			return
		}

		expense := repository.Expense{
			Date:     time.Now(),
			UserID:   c.Sender.ID,
			Category: category,
			Amount:   amount,
		}

		if err = e.repo.AddExpense(expense); err != nil {
			logger.L.Error("Ошибка при добавлении расхода:", err)
		} else {
			msg := fmt.Sprintf(MessagesList.AddedExpense, expense.Date.Format("2006-01-02 15:04:05"), expense.Category, expense.Amount)
			e.sendBotMessage(m, msg) // e.bot.Send(m.Sender, )
		}

		menu.InlineKeyboard = nil
		e.editBotMessageWithMenu(c, MessagesList.EnterAmount, menu) // e.bot.Edit(c.Message, MessagesList.EnterAmount)

		createButtonsMainMenu(e, menu)
		e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu) // e.bot.Send(m.Sender, MessagesList.SelectAction, menu)

		delete(userStates, m.Sender.ID)
	}
}

// Обработчик нажатия кнопки "Мои расходы"
func btnMyExpensesFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnMyExpenses, c.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsOfPeriods(e, menu)
		e.editBotMessageWithMenu(c, MessagesList.SelectPeriod, menu) // e.bot.Edit(c.Message, MessagesList.SelectPeriod, menu)
	}
}

func createButtonsOfPeriods(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	for _, key := range Periods {
		value := BtnPeriodsList[key]
		newBtn := telebot.InlineButton{
			Unique: key,
			Text:   value,
		}
		e.bot.Handle(&newBtn, btnPeriodFunc(e, key, value, menu))

		menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{newBtn})
	}

	btnBack := telebot.InlineButton{
		Unique: "btn_back",
		Text:   BtnTitlesList.BtnBack,
	}
	e.bot.Handle(&btnBack, btnBackFunc(e, menu, "MainMenu"))
	menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})
}

func btnPeriodFunc(e *ExpenseBot, period_key string, period string, menu *telebot.ReplyMarkup) func(callback *telebot.Callback) {
	return func(c *telebot.Callback) {
		menu.InlineKeyboard = nil

		userID := c.Sender.ID
		report := getExpensesByPeriod(e, userID, period_key, period)
		e.editBotMessageWithMenu(c, report, menu)

		createButtonsMainMenu(e, menu)
		e.bot.Send(c.Sender, MessagesList.SelectAction, menu)
		//e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu)
	}
}

// Функция для обработки запроса по расходам в зависимости от периода
func getExpensesByPeriod(e *ExpenseBot, userID int, period_key string, period string) string {

	// Получаем дату начала и конца периода
	startDate, endDate := getPeriodDates(period_key)

	// Получаем данные о расходах из базы данных
	expenses, err := e.repo.GetExpensesByPeriodUnix(userID, startDate, endDate)
	if err != nil {
		logger.L.Error("Ошибка при получении данных.", err)
		return ""
	}

	// Формируем сообщение с результатами
	report := formatExpensesReport(expenses, period)

	return report
}

// Форматирование отчета о расходах
func formatExpensesReport(expenses map[string]float64, period string) string {
	var report strings.Builder
	var totalSum float64

	report.WriteString(fmt.Sprintf("Расходы по категориям за %s:\n", period))
	for category, sum := range expenses {
		report.WriteString(fmt.Sprintf("%s: %.2f\n", category, sum))
		totalSum += sum
	}

	report.WriteString(fmt.Sprintf("\nИтоговая сумма: %.2f", totalSum))
	return report.String()
}

// Функция для расчета даты начала и конца периода
func getPeriodDates(period string) (int64, int64) {
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "period_day":
		startDate = now.Truncate(24 * time.Hour)
		endDate = now
	case "period_week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "period_month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "period_quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "period_halfyear":
		startDate = now.AddDate(0, -6, 0)
		endDate = now
	case "period_year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	}

	return startDate.UnixMilli(), endDate.UnixMilli()
}

func (e *ExpenseBot) handleOnText(menu *telebot.ReplyMarkup) func(*telebot.Message) {
	return func(m *telebot.Message) {
		if lastBotMessage != nil {
			err := e.bot.Delete(lastBotMessage)
			if err != nil {
				logger.L.Error("Ошибка при удалении сообщения:", err)
			}
		}

		// Перехват сообщения от пользователя
		//userMessage := m.Text
		e.sendBotMessage(m, MessagesList.UnknownAction) // e.bot.Send(m.Sender, MessagesList.UnknownAction)

		createButtonsMainMenu(e, menu)
		e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu) // e.bot.Send(m.Sender, MessagesList.SelectAction, menu)
	}
}

func (e *ExpenseBot) sendBotMessage(m *telebot.Message, msg string) {
	sentMessage, err := e.bot.Send(m.Sender, msg)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
	lastBotMessage = sentMessage
}

func (e *ExpenseBot) sendBotMessageWithMenu(m *telebot.Message, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Send(m.Sender, msg, menu)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
	lastBotMessage = sentMessage
}

func (e *ExpenseBot) editBotMessageWithMenu(c *telebot.Callback, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Edit(c.Message, msg, menu)
	if err != nil {
		logger.L.ErrorEditMessage(err)
	}
	lastBotMessage = sentMessage
}

func btnBackFunc(e *ExpenseBot, menu *telebot.ReplyMarkup, backTo string) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnBack, c.Sender.Username))

		menu.InlineKeyboard = nil
		msg := ""

		switch backTo {
		case "MainMenu":
			createButtonsMainMenu(e, menu)
			msg = MessagesList.SelectAction
		case "SelectCategory":
			createButtonsOfCategories(e, menu)
			msg = MessagesList.SelectCategory
		}

		e.editBotMessageWithMenu(c, msg, menu)
	}
}
