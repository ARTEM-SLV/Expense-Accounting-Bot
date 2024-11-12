package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Logger структура для логирования в файл
type Logger struct {
	file *os.File
	mu   sync.Mutex
}

var L *Logger

// InitLogger инициализирует новый логгер с записью в файл
func InitLogger(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(file) // Перенаправляем стандартный вывод логов в файл

	L = &Logger{file: file}

	return nil
}

// Info записывает информационное сообщение
func (l *Logger) Info(msg string) {
	go func() {
		l.printLog("INFO: ", msg)
	}()
}

// Warning записывает сообщение об ошибке
func (l *Logger) Warning(msg string) {
	go func() {
		l.printLog("WARNING: ", msg)
	}()
}

// Error записывает сообщение об ошибке
func (l *Logger) Error(message string, err error) {
	msg := fmt.Sprintf("%s %v", message, err)
	go func() {
		l.printLog("ERROR: ", msg)
	}()
}

// Error ошибка отправки сообщение ботом
func (l *Logger) ErrorSendMessage(err error) {
	msg := fmt.Sprintf("Не удалось отправить сообщение %v", err)
	go func() {
		l.printLog("ERROR: ", msg)
	}()
}

// Error ошибка изменения сообщение ботом
func (l *Logger) ErrorEditMessage(err error) {
	msg := fmt.Sprintf("Не удалось изменить сообщение %v", err)
	go func() {
		l.printLog("ERROR: ", msg)
	}()
}

func (l *Logger) printLog(grade string, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Println(grade, msg)
}

// Close закрывает файл логов
func (l *Logger) Close() error {
	return l.file.Close()
}
