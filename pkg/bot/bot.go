package bot

import (
	"encoding/json"
	"github.com/tucnak/telebot"
	"io"
	"os"
)

var BtnTitlesList BtnTitles
var MessagesList Messages

// Bot интерфейс для бота, поддерживающий различные мессенджеры
type Bot interface {
	Start()
	Handle(command string, handler func(*telebot.Message))
	Send(recipient *telebot.User, message string)
}

type BtnTitles struct {
	BtnMenu       string `json:"btn_menu"`
	BtnBack       string `json:"btn_back"`
	BtnHelp       string `json:"btn_help"`
	BtnMyExpenses string `json:"btn_my_expenses"`
	BtnNewExpense string `json:"btn_new_expense"`
}

type Messages struct {
	Welcome string `json:"welcome"`
}

func InitBot() error {
	// Загружаем заголовки кнопок
	err := loadBtnTitles("./config/button_titles.json")
	if err != nil {
		return err
	}

	// Загружаем сообщения бота
	err = loadMessages("./config/messages.json")
	if err != nil {
		return err
	}

	return nil
}

func loadBtnTitles(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &BtnTitlesList)

	return nil
}

func loadMessages(filePath string) error {
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
