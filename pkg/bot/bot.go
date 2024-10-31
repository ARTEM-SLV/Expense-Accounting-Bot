package bot

import (
	"encoding/json"
	"io"
	"os"
)

var BtnTitlesList *BtnTitles
var MessagesList *Messages
var BtnCategoriesList = make(map[string]string, 8)
var BtnPeriodsList = make(map[string]string, 6)
var Categories = [8]string{"btn_groceries", "btn_beauty", "btn_health", "btn_restaurants", "btn_entertainment", "btn_growth", "btn_trips", "btn_other"}
var Periods = [6]string{"period_day", "period_week", "period_month", "period_quarter", "period_halfyear", "period_year"}

//var Categories map[string]string

// Bot интерфейс для бота, поддерживающий различные мессенджеры
type Bot interface {
	Start()
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
	SelectAction   string `json:"select_action"`
	SelectCategory string `json:"select_category"`
	EnterAmount    string `json:"enter_amount"`
	AddedExpense   string `json:"added_expense"`
	UnknownAction  string `json:"unknown_action"`
	NumberError    string `json:"number_error"`
	SelectPeriod   string `json:"select_period"`
	Category       string `json:"category"`
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
	filePath := "./config/button_titles.json"

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
	filePath := "./config/buttons_categories.json"

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
	filePath := "./config/buttons_periods.json"

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
	filePath := "./config/messages.json"

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &MessagesList)

	return nil
}

//func loadStringValues(filePath string, structValues any) error {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	byteValue, _ := io.ReadAll(file)
//
//	json.Unmarshal(byteValue, &structValues)
//
//	return nil
//}
