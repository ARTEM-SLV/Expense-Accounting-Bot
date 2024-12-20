package repository

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteExpenseRepository реализация ExpenseRepository для SQLite
type SQLiteExpenseRepository struct {
	db *sql.DB
}

// NewSQLiteExpenseRepository создает новый SQLite репозиторий
func NewSQLiteExpenseRepository(db *sql.DB) *SQLiteExpenseRepository {
	return &SQLiteExpenseRepository{db: db}
}

// InitSchema инициализирует таблицу расходов
func (r *SQLiteExpenseRepository) InitSchema() error {
	query := `
    CREATE TABLE IF NOT EXISTS expenses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        date TEXT,
        date_ms INTEGER,
        category TEXT,
        amount REAL
    );
	CREATE INDEX IF NOT EXISTS idx_user_date ON expenses (user_id, date_ms);`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	query = `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		user_name TEXT,
		registered TEXT,
		last_bot_msg_id INTEGER DEFAULT 0,
		chat_id INTEGER DEFAULT 0
	);`
	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}

	query = `
    CREATE TABLE IF NOT EXISTS user_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,        
        category TEXT NOT NULL,
        UNIQUE(user_id, category),
        FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );
    CREATE INDEX IF NOT EXISTS idx_user ON user_categories (user_id);`
	_, err = r.db.Exec(query)
	return err

	return nil
}

func (r *SQLiteExpenseRepository) AddUser(userID int, userName string) error {
	date := time.Now()

	_, err := r.db.Exec(`
        INSERT INTO users (user_id, user_name, registered) VALUES (?, ?, ?)
    `, userID, userName, date.Format("2006-01-02 15:04:05"))

	return err
}

func (r *SQLiteExpenseRepository) GetUserCount() (int, error) {
	query := `
		SELECT count(*) 
		FROM users 
	`
	// Выполнение запроса
	row := r.db.QueryRow(query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLiteExpenseRepository) SetLastBotMsgID(userID int, msgID int, chatID int64) error {
	query := `
        UPDATE users 
        SET last_bot_msg_id = ?, chat_id = ?
        WHERE user_id = ?;
    `

	// Выполнение запроса
	_, err := r.db.Exec(query, msgID, chatID, userID)

	return err
}

func (r *SQLiteExpenseRepository) GetLastBotMsgID(userID int) (int, int64, error) {
	query := `
		SELECT last_bot_msg_id, chat_id 
		FROM users 
		WHERE user_id = ?
	`

	row := r.db.QueryRow(query, userID)

	// Выполнение запроса
	var msgID int
	var chatID int64
	if err := row.Scan(&msgID, &chatID); err != nil {
		return 0, 0, err
	}

	return msgID, chatID, nil
}

func (r *SQLiteExpenseRepository) IsUserRegistered(userID int) (bool, string, error) {
	var isReg bool
	var registered string

	rows, err := r.db.Query(`
        SELECT registered FROM users WHERE user_id = ?
    `, userID)
	if err != nil {
		return false, registered, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&registered); err != nil {
			return false, registered, err
		}
		isReg = true
	}

	return isReg, registered, nil
}

// AddExpense добавляет новый расход в таблицу
func (r *SQLiteExpenseRepository) AddExpense(expense Expense) error {
	date := expense.Date
	dateMs := date.UnixMilli()

	_, err := r.db.Exec(`
        INSERT INTO expenses (user_id, date, date_ms, category, amount) VALUES (?, ?, ?, ?, ?)
    `, expense.UserID, date.Format("2006-01-02 15:04:05"), dateMs, expense.Category, expense.Amount)
	return err
}

// Функция запроса расходов за определенный период из базы данных
func (r *SQLiteExpenseRepository) GetExpensesByPeriod(userID int, startDate, endDate time.Time) (map[string]float64, error) {
	rows, err := r.db.Query(`
        SELECT category, SUM(amount) as total
        FROM expenses
        WHERE user_id = ? date >= ? AND date <= ?
        GROUP BY category
    `, userID, startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make(map[string]float64)
	var category string
	var total float64

	for rows.Next() {
		if err = rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		expenses[category] = total
	}

	return expenses, nil
}

// Метод для получения расходов за период на основе Unix меток времени
func (r *SQLiteExpenseRepository) GetExpensesByPeriodUnix(userID int, startUnixMilli, endUnixMilli int64) (map[string]float64, error) {
	rows, err := r.db.Query(`
        SELECT category, SUM(amount) as total
        FROM expenses
        WHERE user_id = ? AND date_ms >= ? AND date_ms <= ?
        GROUP BY category
    `, userID, startUnixMilli, endUnixMilli)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make(map[string]float64)
	var category string
	var total float64

	for rows.Next() {
		if err = rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		expenses[category] = total
	}

	return expenses, nil
}

func (r *SQLiteExpenseRepository) AddUserCategory(userID int, category string) error {
	_, err := r.db.Exec(`
        INSERT INTO user_categories (user_id, category) VALUES (?, ?)
    `, userID, category)

	return err
}

func (r *SQLiteExpenseRepository) GetUserCategories(userID int) ([]string, error) {
	var categories []string

	rows, err := r.db.Query(`
        SELECT registered FROM users WHERE user_id = ?
    `, userID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&categories); err != nil {
			return categories, err
		}
	}

	return categories, nil
}
