package repository

import (
	"time"
)

// Expense структура для хранения данных о расходах
type Expense struct {
	UserID   int
	Date     time.Time
	Category string
	Amount   float64
}

// ExpenseRepository интерфейс для работы с расходами
type ExpenseRepository interface {
	InitSchema() error
	AddUser(userID int, userName string) error
	GetUserCount() (int, error)
	SetLastBotMsgID(userID int, msgID int, chatID int64) error
	GetLastBotMsgID(userID int) (int, int64, error)
	IsUserRegistered(userID int) (bool, string, error)
	AddExpense(expense Expense) error
	GetExpensesByPeriod(userID int, startDate, endDate time.Time) (map[string]float64, error)
	GetExpensesByPeriodUnix(userID int, tartUnixMilli, endUnixMilli int64) (map[string]float64, error)
}
