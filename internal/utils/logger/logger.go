package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger структура для логирования в файл
type Logger struct {
	file *os.File
}

// NewLogger инициализирует новый логгер с записью в файл
func NewLogger(filepath string) (*Logger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file) // Перенаправляем стандартный вывод логов в файл

	return &Logger{file: file}, nil
}

// Info записывает информационное сообщение
func (l *Logger) Info(message string) {
	log.Println("INFO: ", message)
}

// Error записывает сообщение об ошибке
func (l *Logger) Error(message string, err error) {
	msg := fmt.Sprintf("%s %s", message, err)
	log.Println("INFO: ", msg)
}

// Close закрывает файл логов
func (l *Logger) Close() error {
	return l.file.Close()
}
