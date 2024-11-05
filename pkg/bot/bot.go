package bot

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

var BtnTitlesList *BtnTitles
var MessagesList *Messages
var BtnCategoriesList = make(map[string]string, 8)
var BtnPeriodsList = make(map[string]string, 6)
var Categories = [10]string{"btn_groceries", "btn_beauty", "btn_health", "btn_restaurants", "btn_entertainment",
	"btn_growth", "btn_trips", "btn_transport", "btn_business", "btn_other"}
var Periods = [6]string{"period_day", "period_week", "period_month", "period_quarter", "period_halfyear", "period_year"}

//var Categories map[string]string

// Bot интерфейс для бота, поддерживающий различные мессенджеры
type Bot interface {
	Start()
	StartDailyReport(adminID int)
}

type BtnTitles struct {
	BtnMenu string `json:"btn_menu"`
	BtnBack string `json:"btn_back"`
	BtnHelp string `json:"btn_help"`

	BtnNewExpense string `json:"btn_new_expense"`
	BtnMyExpenses string `json:"btn_my_expenses"`
}

type Messages struct {
	Welcome        string `json:"welcome"`
	Help           string `json:"help"`
	SelectAction   string `json:"select_action"`
	SelectCategory string `json:"select_category"`
	EnterAmount    string `json:"enter_amount"`
	AddedExpense   string `json:"added_expense"`
	UnknownAction  string `json:"unknown_action"`
	NumberError    string `json:"number_error"`
	SelectPeriod   string `json:"select_period"`
	Category       string `json:"category"`
	Period         string `json:"period"`
	ErrorReg       string `json:"error_reg"`
	UserRegistered string `json:"user_registered"`
}

func InitStringValues() error {
	// Загружаем заголовки кнопок
	err := loadBtnTitles()
	if err != nil {
		return err
	}

	// Загружаем заголовки кнопок категории
	err = loadBtnCategories()
	if err != nil {
		return err
	}

	// Загружаем заголовки кнопок периоды
	err = loadBtnPeriods()
	if err != nil {
		return err
	}

	// Загружаем сообщения бота
	err = loadMessages()
	if err != nil {
		return err
	}

	return nil
}

func loadBtnTitles() error {
	filePath := "./config/string_values/button_titles.json"

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &BtnTitlesList)

	return nil
}

func loadBtnCategories() error {
	filePath := "./config/string_values/buttons_categories.json"

	// Открываем JSON файл
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Декодируем JSON в мапу
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&BtnCategoriesList); err != nil {
		return err
	}

	return nil
}

func loadBtnPeriods() error {
	filePath := "./config/string_values/buttons_periods.json"

	// Открываем JSON файл
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Декодируем JSON в мапу
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&BtnPeriodsList); err != nil {
		return err
	}

	return nil
}

func loadMessages() error {
	filePath := "./config/string_values/messages.json"

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &MessagesList)

	return nil
}

// Функция для расчета даты начала и конца периода
func GetPeriodDates(period string) (int64, int64) {
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
