# 💰 Expense Tracker Telegram Bot

A simple and lightweight Telegram bot for tracking personal expenses.  
It helps users quickly log expenses, categorize them, and view summaries over different time periods.

---

## 🚀 Features

- ➕ Add expenses via interactive buttons  
- 📂 Predefined expense categories:
  - Groceries
  - Beauty
  - Health
  - Restaurants
  - Entertainment
  - Education
  - Travel
  - Other
- 📊 View expenses by period:
  - Day
  - Week
  - Month
  - Quarter
  - Half-year
  - Year
- 🧾 Automatic grouping by category  
- 📈 Total expenses calculation  
- 👤 Multi-user support  
- 🔐 User registration (`/start`)  
- ❓ Help command (`/help`)  
- 💾 Data storage using SQLite  

---

## 🛠️ Tech Stack

- Go (Golang)  
- Telegram Bot API (telebot)  
- SQLite  

---

## ⚙️ Setup & Run

### 1. Clone repository

```bash
git clone https://github.com/your-repo/expense-bot.git
cd expense-bot
```

### 2. Configure environment variables

```bash
Create a .env file (for local development):
TELEGRAM_TOKEN=your_telegram_bot_token
ADMIN_ID=your_telegram_id
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the bot

```bash
go run ./cmd
```

🐳 Docker (optional)
docker build -t expense-bot .
docker run -d expense-bot

💬 Bot Commands
Command	Description
/start	Register user and show main menu
/help	Show help information

🧭 How It Works
User starts the bot with /start
Bot shows main menu:
"Add Expense"
"My Expenses"
User selects a category
User enters amount
Expense is saved to SQLite
User can view reports grouped by category and period
📌 Roadmap
 Custom user-defined categories
 Export data (CSV / Excel)
 Charts and analytics
 Notifications and daily reports
 Migration to PostgreSQL
 REST API
🤝 Contributing

Contributions are welcome!
Feel free to open issues or submit pull requests.

📄 License

This project is licensed under the MIT License.
