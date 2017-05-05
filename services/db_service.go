package services

import (
	"github.com/tradesim/model"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"encoding/json"
	"time"
)

var db *sql.DB
var InitialBalance = map[string]interface{} {
	"USD": 10000.00,
	"BTC": 0.0,
}

func Initialize() {
	tempdb, err := sql.Open("postgres", "postgres://financedb:financedb@10.32.0.4/financedb")
	if err != nil {
		panic(err)
	}
	db = tempdb
	InitializeSessions()
}

// User functions

func GetUserByUsername(username string) model.User {
	var user model.User
	err := db.QueryRow("SELECT user_id, username, email FROM users WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		panic(err)
	}
	return user
}

func GetUserByEmail(email string) model.User {
	var user model.User
	err := db.QueryRow("SELECT user_id, username, email FROM users WHERE email = $1", email).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		panic(err)
	}
	return user
}

func UserExists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func EmailExists(email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func AddUser(email string, username string, password string) error {

	hash := HashString(password)
	_, err := db.Exec("INSERT INTO users (username, email, password_hash) values($1, $2, $3)", username, email, hash)
	if err != nil {
		fmt.Println(err)
	} else {
		err = UpdateBalance(username, InitialBalance)
	}
	return err
}

func UpdatePassword(update model.UpdatePassword, username string) error {
	newHash := HashString(update.NewPassword)
	oldHash := HashString(update.OldPassword)
	result, err := db.Exec("UPDATE users SET password_hash = $1 WHERE username = $2 AND password_hash = $3", newHash, username, oldHash)
	if result != nil {
		updated, _ := result.RowsAffected()
		if updated != 1 {
			return fmt.Errorf("Failed to verify old password for provided username")
		}
	}
	return err
}

func CheckEmailVerified(username string) bool {
	rows, _ := db.Query("SELECT last_login FROM users WHERE username = $1", username)
	var last_login time.Time
	count := 0
	for rows.Next() {
		err := rows.Scan(&last_login)
		if err != nil {
			return false
		}
		count++
	}
	if count != 1 {
		return false
	}
	fmt.Println("LAST_LOGIN: " + last_login.String())
	if last_login.String() == "" {
		return false
	}
	return true
}



func Login(login model.Login) (string, error) {
	fmt.Println("Checking if email is verified")
	if !CheckEmailVerified(login.Username) {
		return "", fmt.Errorf("This account's email has not been verified. Please check your email for a verification link.")
	}
	pwHash := HashString(login.Password)
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND password_hash = $2)", login.Username, pwHash).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if !exists {
		return "", fmt.Errorf("Failed to validate username and password combination")
	}
	err = UpdateUserLastLogin(login.Username)
	if err != nil {
		panic(err)
	}
	return CreateSession(login.Username), nil
}

func UpdateUserLastLogin(username string) error {
	result, err := db.Exec("UPDATE users SET last_login = current_timestamp WHERE username = $1", username)
	if result != nil {
		updated, err2 := result.RowsAffected()
		if updated != 1 || err != nil || err2 != nil {
			return fmt.Errorf("User not found")
		}
	}
}

// Trade simulator functions

func GetBalance(user string) map[string]interface{} {
	var balByte []byte
	err := db.QueryRow("SELECT balances FROM balances WHERE user_id = (SELECT user_id FROM users WHERE username = $1) AND balance_datetime = (SELECT MAX(balance_datetime) FROM balances)", user).Scan(&balByte)
	if err != nil {
		panic(err)
	}
	var bal map[string]interface{}
	err = json.Unmarshal(balByte, &bal)
	if err != nil {
		panic(err)
	}
	return bal
}

func UpdateBalance(user string, newBalances map[string]interface{}) error {
	balString, _ := json.Marshal(newBalances)
	_, err := db.Exec("INSERT INTO balances (user_id, balances) VALUES((SELECT user_id FROM users WHERE username = $1), $2)", user, balString)
	return err
}

func AddTransaction(username string, symbol string, tradetype string, quantity float64, rate float64, fee float64) error {
	_, err := db.Exec("INSERT INTO transactions (user_id, symbol, tradetype, quantity, rate, fee_amount) VALUES((SELECT user_id FROM users WHERE username = $1), $2, $3, $4, $5, $6)", username, symbol, tradetype, quantity, rate, fee)
	return err
}

func ResetBalance(username string) error {
	err := UpdateBalance(username, InitialBalance)
	return err
}

func GetTransactions(user string) []model.TransactionDetail {
	transactions := make([]model.TransactionDetail, 0)
	rows, err := db.Query("SELECT occurred_at, tradetype, quantity, symbol, rate, fee_amount from transactions where user_id = (select user_id from users where username = $1);", user)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var temp model.TransactionDetail
		rows.Scan(&temp.Time, &temp.Action, &temp.Quantity, &temp.Symbol, &temp.Rate, &temp.Fee)
		transactions = append(transactions, temp)
	}
	return transactions
}










