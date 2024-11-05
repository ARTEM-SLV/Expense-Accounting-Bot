package telegram

import (
	"fmt"

	"github.com/tucnak/telebot"

	"expense_accounting_bot/internal/utils/logger"
	"expense_accounting_bot/pkg/bot"
	"expense_accounting_bot/pkg/repository"
)

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
			logger.L.Error(bot.MessagesList.ErrorReg, err)
			sendBotMessage(e, m, bot.MessagesList.ErrorReg)
			return
		}

		if isRegistered {
			deleteBotMessage(e, userID)

			sendBotMessage(e, m, fmt.Sprintf(bot.MessagesList.UserRegistered, userName, dateReg))
			sendBotMessageWithMenu(e, m, bot.MessagesList.SelectAction, menu)

			return
		}

		// Регистрируем нового пользователя
		err = e.repo.AddUser(userID, userName)
		if err != nil {
			logger.L.Error(bot.MessagesList.ErrorReg, err)
			sendBotMessage(e, m, bot.MessagesList.ErrorReg)
			return
		}

		msg := fmt.Sprintf(bot.MessagesList.Welcome, m.Sender.FirstName, m.Sender.LastName)
		sendBotMessage(e, m, msg)

		sendBotMessageWithMenu(e, m, bot.MessagesList.SelectAction, menu)
	})

	createButtonsMainMenu(e, menu)

	// Запуск бота
	e.bot.Start()
}

func createButtonsMainMenu(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	// Создаем кнопки
	btnNewExpense := telebot.InlineButton{
		Unique: "btn_schedule",
		Text:   bot.BtnTitlesList.BtnNewExpense,
	}
	btnMyExpenses := telebot.InlineButton{
		Unique: "btn_services",
		Text:   bot.BtnTitlesList.BtnMyExpenses,
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

func (e *ExpenseBot) handleOnText(menu *telebot.ReplyMarkup) func(*telebot.Message) {
	return func(m *telebot.Message) {
		userID := m.Sender.ID
		deleteBotMessage(e, userID)

		// Перехват сообщения от пользователя
		if m.Text == "/help" {
			sendBotMessage(e, m, bot.MessagesList.Help)
		} else {
			sendBotMessage(e, m, bot.MessagesList.UnknownAction)
		}

		createButtonsMainMenu(e, menu)
		sendBotMessageWithMenu(e, m, bot.MessagesList.SelectAction, menu)
	}
}

func btnBackFunc(e *ExpenseBot, menu *telebot.ReplyMarkup, backTo string) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", bot.BtnTitlesList.BtnBack, c.Sender.Username))

		menu.InlineKeyboard = nil
		msg := ""

		switch backTo {
		case "MainMenu":
			createButtonsMainMenu(e, menu)
			msg = bot.MessagesList.SelectAction
		case "SelectCategory":
			createButtonsOfCategories(e, menu)
			msg = bot.MessagesList.SelectCategory
		}

		editBotMessageWithMenu(e, c, msg, menu)
	}
}

func sendBotMessage(e *ExpenseBot, m *telebot.Message, msg string) {
	sentMessage, err := e.bot.Send(m.Sender, msg)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}

	err = e.repo.SetLastBotMsgID(m.Sender.ID, sentMessage.ID, m.Chat.ID)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
}

func sendBotMessageWithMenu(e *ExpenseBot, m *telebot.Message, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Send(m.Sender, msg, menu)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}

	err = e.repo.SetLastBotMsgID(m.Sender.ID, sentMessage.ID, m.Chat.ID)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
}

func editBotMessageWithMenu(e *ExpenseBot, c *telebot.Callback, msg string, menu *telebot.ReplyMarkup) {
	sentMessage, err := e.bot.Edit(c.Message, msg, menu)
	if err != nil {
		logger.L.ErrorEditMessage(err)
	}

	err = e.repo.SetLastBotMsgID(c.Sender.ID, sentMessage.ID, c.Message.Chat.ID)
	if err != nil {
		logger.L.ErrorSendMessage(err)
	}
}

func deleteBotMessage(e *ExpenseBot, userID int) {
	msg := "Ошибка при удалении сообщения:"

	lastBotMsg, err := getUserMessage(e, userID)
	if err != nil {
		logger.L.Error(msg, err)
		return
	}

	if lastBotMsg != nil {
		err = e.bot.Delete(lastBotMsg)
		if err != nil {
			logger.L.Error(msg, err)
		}
	}
}

// Получение сообщения по user_id
func getUserMessage(e *ExpenseBot, userID int) (*telebot.Message, error) {
	messageID, chatID, err := e.repo.GetLastBotMsgID(userID)
	if err != nil {
		return nil, err
	}

	if messageID == 0 || chatID == 0 {
		return nil, nil
	}

	message := &telebot.Message{
		ID:   messageID,
		Chat: &telebot.Chat{ID: chatID},
	}

	return message, nil
}
