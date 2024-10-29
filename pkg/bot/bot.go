package bot

import (
	"encoding/json"
	"io"
	"os"

	"github.com/tucnak/telebot"
)

var BtnTitlesList *BtnTitles
var BtnCategoriesList *BtnCategories
var MessagesList *Messages

// Bot интерфейс для бота, поддерживающий различные мессенджеры
type Bot interface {
	Start()
	Handle(command string, handler func(*telebot.Message))
	Send(recipient *telebot.User, message string)
}

type BtnTitles struct {
	BtnMenu string `json:"btn_menu"`
	BtnBack string `json:"btn_back"`
	BtnHelp string `json:"btn_help"`

	BtnNewExpense string `json:"btn_new_expense"`
	BtnMyExpenses string `json:"btn_my_expenses"`
}

type BtnCategories struct {
	BtnGroceries     string `json:"btn_groceries"`
	BtnBeauty        string `json:"btn_beauty"`
	BtnHealth        string `json:"btn_health"`
	BtnRestaurants   string `json:"btn_restaurants"`
	BtnEntertainment string `json:"btn_entertainment"`
	BtnGrowth        string `json:"btn_growth"`
	BtnTrips         string `json:"btn_trips"`
	BtnOther         string `json:"btn_other"`
}

type Messages struct {
	Welcome      string `json:"welcome"`
	SelectAction string `json:"select_action"`
	Category     string `json:"category"`
}

func InitStringValues() error {
	// Загружаем заголовки кнопок
	err := loadBtnTitles()
	if err != nil {
		return err
	}

	// Загружаем сообщения бота
	err = loadBtnCategories()
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

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &BtnCategoriesList)

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
