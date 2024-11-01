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
var lastBotMessage = make(map[int]*telebot.Message)

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
			sendBotMessage(e, m, MessagesList.ErrorReg)
			return
		}

		if isRegistered {
			if lastBotMessage[userID] != nil {
				err = e.bot.Delete(lastBotMessage[userID])
			}
			if err != nil {
				logger.L.Error("Ошибка при удалении сообщения:", err)
			}
			sendBotMessage(e, m, fmt.Sprintf(MessagesList.UserRegistered, userName, dateReg))
			sendBotMessageWithMenu(e, m, MessagesList.SelectAction, menu)
			return
		}

		// Регистрируем нового пользователя
		err = e.repo.AddUser(userID, userName)
		if err != nil {
			logger.L.Error(MessagesList.ErrorReg, err)
			sendBotMessage(e, m, MessagesList.ErrorReg)
			return
		}

		msg := fmt.Sprintf(MessagesList.Welcome, m.Sender.FirstName, m.Sender.LastName)
		sendBotMessage(e, m, msg)

		sendBotMessageWithMenu(e, m, MessagesList.SelectAction, menu)
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

		editBotMessageWithMenu(e, c, MessagesList.SelectCategory, menu) // e.bot.Edit(c.Message, MessagesList.SelectCategory, menu)
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

		editBotMessageWithMenu(e, c, MessagesList.EnterAmount, menu)

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
			sendBotMessage(e, m, MessagesList.NumberError)
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
			sendBotMessage(e, m, msg) // e.bot.Send(m.Sender, )
		}

		menu.InlineKeyboard = nil
		editBotMessageWithMenu(e, c, MessagesList.EnterAmount, menu) // e.bot.Edit(c.Message, MessagesList.EnterAmount)

		createButtonsMainMenu(e, menu)
		sendBotMessageWithMenu(e, m, MessagesList.SelectAction, menu) // e.bot.Send(m.Sender, MessagesList.SelectAction, menu)

		delete(userStates, m.Sender.ID)
	}
}

// Обработчик нажатия кнопки "Мои расходы"
func btnMyExpensesFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", BtnTitlesList.BtnMyExpenses, c.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsOfPeriods(e, menu)
		editBotMessageWithMenu(e, c, MessagesList.SelectPeriod, menu) // e.bot.Edit(c.Message, MessagesList.SelectPeriod, menu)
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
		e.bot.Respond(c, &telebot.CallbackResponse{Text: fmt.Sprintf(MessagesList.Period, period)})

		menu.InlineKeyboard = nil

		userID := c.Sender.ID
		report := getExpensesByPeriod(e, userID, period_key, period)
		editBotMessageWithMenu(e, c, report, menu)

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

	currYear := now.Year()
	currMonth := now.Month()
	currDay := now.Day()

	switch period {
	case "period_day":
		// Начало текущего дня
		startDate = time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, now.Location())
		// Конец текущего дня
		endDate = time.Date(currYear, currMonth, currDay, 23, 59, 59, int(time.Nanosecond*999999999), now.Location())
	case "period_week":
		// Получаем день недели (в Go воскресенье - 0, понедельник - 1, ... суббота - 6)
		weekday := int(now.Weekday())
		// Смещение для начала недели (если неделя начинается с понедельника)
		offset := (weekday + 6) % 7 // Количество дней, на которые нужно сдвинуть назад, чтобы получить понедельник
		// Начало недели (понедельник текущей недели)
		startDate = now.AddDate(0, 0, -offset).Truncate(24 * time.Hour)
		// Конец недели (воскресенье текущей недели)
		endDate = startDate.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Nanosecond*999999999)
	case "period_month":
		// Начало текущего месяца
		startDate = time.Date(currYear, currMonth, 1, 0, 0, 0, 0, now.Location())
		// Конец текущего месяца
		endDate = startDate.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Nanosecond*999999999)
	case "period_quarter":
		// Определяем начало квартала
		switch now.Month() {
		case time.January, time.February, time.March:
			startDate = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		case time.April, time.May, time.June:
			startDate = time.Date(now.Year(), time.April, 1, 0, 0, 0, 0, now.Location())
		case time.July, time.August, time.September:
			startDate = time.Date(now.Year(), time.July, 1, 0, 0, 0, 0, now.Location())
		case time.October, time.November, time.December:
			startDate = time.Date(now.Year(), time.October, 1, 0, 0, 0, 0, now.Location())
		}
		// Конец квартала — добавляем 3 месяца к началу квартала и вычитаем 1 день
		endDate = startDate.AddDate(0, 3, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Nanosecond*999999999)
	case "period_halfyear":
		// Определяем начало полугодия
		if now.Month() <= time.June { // Первое полугодие: Январь - Июнь
			startDate = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		} else { // Второе полугодие: Июль - Декабрь
			startDate = time.Date(now.Year(), time.July, 1, 0, 0, 0, 0, now.Location())
		}
		// Конец полугодия: добавляем 6 месяцев к началу полугодия и вычитаем 1 день
		endDate = startDate.AddDate(0, 6, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Nanosecond*999999999)
	case "period_year":
		// Начало года: 1 января текущего года
		startDate = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		// Конец года: 31 декабря текущего года, конец дня
		endDate = time.Date(now.Year(), time.December, 31, 23, 59, 59, int(time.Nanosecond*999999999), now.Location())
	}

	return startDate.UnixMilli(), endDate.UnixMilli()
}

func (e *ExpenseBot) handleOnText(menu *telebot.ReplyMarkup) func(*telebot.Message) {
	return func(m *telebot.Message) {
		userID := m.Sender.ID
		if lastBotMessage[userID] != nil {
			err := e.bot.Delete(lastBotMessage[userID])
			if err != nil {
				logger.L.Error("Ошибка при удалении сообщения:", err)
			}
		}

		// Перехват сообщения от пользователя
		if m.Text == "/help" {
			sendBotMessage(e, m, MessagesList.Help)
		} else {
			sendBotMessage(e, m, MessagesList.UnknownAction)
		}

		createButtonsMainMenu(e, menu)
		sendBotMessageWithMenu(e, m, MessagesList.SelectAction, menu)
	}
}

func sendBotMessage(e *ExpenseBot, m *telebot.Message, msg string) {
	sentMessage, err := e.bot.Send(m.Sender, msg)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
	lastBotMessage[m.Sender.ID] = sentMessage
}

func sendBotMessageWithMenu(e *ExpenseBot, m *telebot.Message, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Send(m.Sender, msg, menu)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
	lastBotMessage[m.Sender.ID] = sentMessage
}

func editBotMessageWithMenu(e *ExpenseBot, c *telebot.Callback, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Edit(c.Message, msg, menu)
	if err != nil {
		logger.L.ErrorEditMessage(err)
	}
	lastBotMessage[c.Sender.ID] = sentMessage
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

		editBotMessageWithMenu(e, c, msg, menu)
	}
}
