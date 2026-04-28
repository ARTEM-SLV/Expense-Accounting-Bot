# 💰 Expense Accounting Bot

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

---

🐳 Docker (optional)<br>
docker build -t expense-bot<br>
docker run -d expense-bot<br>

---

💬 Bot Commands<br>
Command	Description<br>
/start	Register user and show main menu<br>
/help	Show help information<br>

---

🧭 How It Works<br>
User starts the bot with /start<br>
Bot shows main menu:<br>
"Add Expense"<br>
"My Expenses"<br>
User selects a category<br>
User enters amount<br>
Expense is saved to SQLite<br>
User can view reports grouped by category and period<br>
📌 Roadmap<br>
 Custom user-defined categories<br>
 Export data (CSV / Excel)<br>
 Charts and analytics<br>
 Notifications and daily reports<br>
 Migration to PostgreSQL<br>
 REST API<br>
🤝 Contributing<br>

Contributions are welcome!
Feel free to open issues or submit pull requests.
