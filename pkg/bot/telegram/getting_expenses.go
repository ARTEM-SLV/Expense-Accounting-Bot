package telegram

import (
	"fmt"
	"strings"

	"github.com/tucnak/telebot"

	"expense_accounting_bot/internal/utils/logger"
	"expense_accounting_bot/pkg/bot"
)

// Обработчик нажатия кнопки "Мои расходы"
func btnMyExpensesFunc(e *ExpenseBot, menu *telebot.ReplyMarkup) func(*telebot.Callback) {
	return func(c *telebot.Callback) {
		logger.L.Info(fmt.Sprintf("Нажата кнопка '%s' пользователем %s", bot.BtnTitlesList.BtnMyExpenses, c.Sender.Username))

		menu.InlineKeyboard = nil

		createButtonsOfPeriods(e, menu)
		editBotMessageWithMenu(e, c, bot.MessagesList.SelectPeriod, menu) // e.bot.Edit(c.Message, MessagesList.SelectPeriod, menu)
	}
}

func createButtonsOfPeriods(e *ExpenseBot, menu *telebot.ReplyMarkup) {
	for _, key := range bot.Periods {
		value := bot.BtnPeriodsList[key]
		newBtn := telebot.InlineButton{
			Unique: key,
			Text:   value,
		}
		e.bot.Handle(&newBtn, btnPeriodFunc(e, key, value, menu))

		menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{newBtn})
	}

	btnBack := telebot.InlineButton{
		Unique: "btn_back",
		Text:   bot.BtnTitlesList.BtnBack,
	}
	e.bot.Handle(&btnBack, btnBackFunc(e, menu, "MainMenu"))
	menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{btnBack})
}

func btnPeriodFunc(e *ExpenseBot, period_key string, period string, menu *telebot.ReplyMarkup) func(callback *telebot.Callback) {
	return func(c *telebot.Callback) {
		e.bot.Respond(c, &telebot.CallbackResponse{Text: fmt.Sprintf(bot.MessagesList.Period, period)})

		menu.InlineKeyboard = nil

		userID := c.Sender.ID
		report := getExpensesByPeriod(e, userID, period_key, period)
		editBotMessageWithMenu(e, c, report, menu)

		createButtonsMainMenu(e, menu)
		e.bot.Send(c.Sender, bot.MessagesList.SelectAction, menu)
		//e.sendBotMessageWithMenu(m, MessagesList.SelectAction, menu)
	}
}

// Функция для обработки запроса по расходам в зависимости от периода
func getExpensesByPeriod(e *ExpenseBot, userID int, period_key string, period string) string {

	// Получаем дату начала и конца периода
	startDate, endDate := bot.GetPeriodDates(period_key)

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
