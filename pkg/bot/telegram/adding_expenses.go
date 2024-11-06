package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tucnak/telebot"

	"expense_accounting_bot/internal/utils/logger"
	"expense_accounting_bot/pkg/bot"
	"expense_accounting_bot/pkg/repository"
)

// Обработчик нажатия кнопки "Добавить расход"
func btnNewExpenseFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", bot.BtnTitlesList.BtnNewExpense, c.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsOfCategories(e, menu)

		editBotMessageWithMenu(e, c, bot.MessagesList.SelectCategory, menu)
	}
}

func createButtonsOfCategories(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	row := make([]telebot.InlineButton, 0, 2)
	for _, key := range bot.Categories {
		value := bot.BtnCategoriesList[key]
		addBtnOfCategory(&row, e, key, value, menu)

		if len(row) == 2 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, row)
			row = make([]telebot.InlineButton, 0, 2)
		}
	}

	btnBack := telebot.InlineButton{
		Unique: "btn_back",
		Text:   bot.BtnTitlesList.BtnBack,
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

		e.bot.Respond(c, &telebot.CallbackResponse{Text: fmt.Sprintf(bot.MessagesList.Category, category)})

		//userID := c.Sender.ID
		//userStates[userID] = category

		btnBack := telebot.InlineButton{
			Unique: "btn_back",
			Text:   bot.BtnTitlesList.BtnBack,
		}
		e.bot.Handle(&btnBack, btnBackFunc(e, menu, "SelectCategory"))
		menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})

		editBotMessageWithMenu(e, c, bot.MessagesList.EnterAmount, menu)

		e.bot.Handle(telebot.OnText, addExpense(e, menu, c, category))
	}
}

func addExpense(e *ExpenseBot, menu *telebot.ReplyMarkup, c *telebot.Callback, category string) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//category, ok := userStates[m.Sender.ID]
		//if !ok || category == "" {
		//	return
		//}

		amount, err := strconv.ParseFloat(m.Text, 64)
		if err != nil {
			sendBotMessage(e, m, bot.MessagesList.NumberError)
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
			msg := fmt.Sprintf(bot.MessagesList.AddedExpense, expense.Date.Format("2006-01-02 15:04:05"), expense.Category, expense.Amount)
			sendBotMessage(e, m, msg)
		}

		menu.InlineKeyboard = nil
		editBotMessageWithMenu(e, c, bot.MessagesList.EnterAmount, menu)

		createButtonsMainMenu(e, menu)
		sendBotMessageWithMenu(e, m, bot.MessagesList.SelectAction, menu)
	}
}
