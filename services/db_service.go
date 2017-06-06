package services

import (
	"github.com/tradesim/model"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"encoding/json"
	"time"
	"os"
	"strings"
	"strconv"
)

var db *sql.DB
var InitialBalance = map[string]interface{} {
	"USD": 50000.00,
	"BTC": 0.0,
}
var UserCommentLimits map[string]int64

func Initialize() {
	tempdb, err := sql.Open("postgres", "postgres://financedb:financedb@" + os.Getenv("POSTGRES_HOST") + "/financedb")
	if err != nil {
		panic(err)
	}
	db = tempdb
	InitRates()
	InitializeSessions()

	UserCommentLimits = make(map[string]int64)
	go ManageCommentLimits()
}

// User functions

func GetAllUsers() []model.User {
	var users []model.User
	rows, _ := db.Query("SELECT username FROM users")
	for rows.Next() {
		var temp model.User
		err := rows.Scan(&temp.Username)
		if err != nil {
			panic(err)
		}
		users = append(users, temp)
	}
	return users
}

func GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT user_id, username, email FROM users WHERE username = $1", strings.ToLower(username)).Scan(&user.Id, &user.Username, &user.Email)
	return user, err
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
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", strings.ToLower(username)).Scan(&count)
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
	result, err := db.Exec("UPDATE users SET password_hash = $1 WHERE username = $2 AND password_hash = $3", newHash, strings.ToLower(username), oldHash)
	if result != nil {
		updated, _ := result.RowsAffected()
		if updated != 1 {
			return fmt.Errorf("Failed to verify old password for provided username")
		}
	}
	return err
}

func UpdatePasswordForce(newPass string, username string) error {
	newHash := HashString(newPass)
	result, err := db.Exec("UPDATE users SET password_hash = $1 WHERE username = $2", newHash, strings.ToLower(username))
	if result != nil {
		updated, _ := result.RowsAffected()
		if updated != 1 {
			return fmt.Errorf("Provided username not found")
		}
	}
	return err
}

func CheckEmailVerified(username string) bool {
	rows, _ := db.Query("SELECT last_login FROM users WHERE username = $1", strings.ToLower(username))
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
	pwHash := HashString(login.Password)
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND password_hash = $2)", login.Username, pwHash).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if !exists {
		return "", fmt.Errorf("Failed to validate username and password combination")
	}
	fmt.Println("Checking if email is verified")
	if !CheckEmailVerified(login.Username) {
		return "", fmt.Errorf("This account's email has not been verified. Please check your email for a verification link.")
	}
	err = UpdateUserLastLogin(login.Username)
	if err != nil {
		panic(err)
	}
	return CreateSession(login.Username), nil
}

func UpdateUserLastLogin(username string) error {
	result, err := db.Exec("UPDATE users SET last_login = current_timestamp WHERE username = $1", strings.ToLower(username))
	if result != nil {
		updated, err2 := result.RowsAffected()
		if updated != 1 || err != nil || err2 != nil {
			return fmt.Errorf("User not found")
		}
		return nil
	}
	panic(err)
}

// Trade simulator functions

func GetBalance(user string) map[string]interface{} {
	var balByte []byte
	user = strings.ToLower(user)
	err := db.QueryRow("SELECT balances FROM balances WHERE user_id = (SELECT user_id FROM users WHERE username = $1) AND balance_datetime = (SELECT MAX(balance_datetime) FROM balances WHERE user_id = (SELECT user_id FROM users WHERE username = $2))", user, user).Scan(&balByte)
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
	for symbol, quantity := range newBalances {
		if quantity.(float64) == 0 {
			delete(newBalances, symbol)
		}
	}
	balString, _ := json.Marshal(newBalances)
	_, err := db.Exec("INSERT INTO balances (user_id, balances) VALUES((SELECT user_id FROM users WHERE username = $1), $2)", strings.ToLower(user), balString)
	return err
}

func AddTransaction(username string, symbol string, tradetype string, quantity float64, rate float64, fee float64) error {
	_, err := db.Exec("INSERT INTO transactions (user_id, symbol, tradetype, quantity, rate, fee_amount) VALUES((SELECT user_id FROM users WHERE username = $1), $2, $3, $4, $5, $6)", strings.ToLower(username), symbol, tradetype, quantity, rate, fee)
	return err
}

func ResetBalance(username string) error {
	err := UpdateBalance(strings.ToLower(username), InitialBalance)
	return err
}

func GetTransactions(user string) []model.TransactionDetail {
	transactions := make([]model.TransactionDetail, 0)
	rows, err := db.Query("SELECT occurred_at, tradetype, quantity, symbol, rate, fee_amount FROM transactions WHERE user_id = (select user_id FROM users WHERE username = $1);", strings.ToLower(user))
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

func GetArticle(articleId string) model.Article {
	var article model.Article
	fmt.Println("Finding article with id " + articleId)
	err := db.QueryRow("SELECT author, title, creation_time, content FROM articles WHERE article_id = $1", articleId).Scan(&article.Author, &article.Title, &article.CreatedDate, &article.Content)
	if err != nil {
		fmt.Println("Warning: Error selecting row from DB - " + err.Error())
	}
	fmt.Println("Article title: " + article.Title)
	return article
}

func GetArticles(limit int) []model.Article {
	var articles []model.Article
	sql := "SELECT article_id, thumbnail_url, title, creation_time FROM articles ORDER BY creation_time DESC"
	if limit > 0 {
		sql +=  " LIMIT " + strconv.Itoa(limit)
	}
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("Warning: Error selecting rows from DB - " + err.Error())
	} else {
		for rows.Next() {
			var tempArticle model.Article
			rows.Scan(&tempArticle.Id, &tempArticle.ThumbnailUrl, &tempArticle.Title, &tempArticle.CreatedDate)
			articles = append(articles, tempArticle)
		}
	}
	
	return articles
}

// COMMENTS

func GetCommentsForArticle(articleid string) []model.Comment {
	var comments []model.Comment
	sql := "SELECT comment_time, username, comment_id, content FROM comments JOIN users ON comments.user_id = users.user_id WHERE comments.article_id = $1"
	rows, err := db.Query(sql, articleid)
	if err != nil {
		fmt.Println("Warning: Error selecting rows from DB - " + err.Error())
	} else {
		for rows.Next() {
			var tempComment model.Comment
			rows.Scan(&tempComment.Time, &tempComment.Username, &tempComment.Id, &tempComment.Content)
			comments = append(comments, tempComment)
		}
	}

	return comments
}

func AddComment(articleid string, username string, content string) {
	if _, ok := UserCommentLimits[username]; ok {
		if (UserCommentLimits[username] > 0) {
			UserCommentLimits[username] = max(UserCommentLimits[username] - 1, 0)
			_, err := db.Exec("INSERT INTO comments (user_id, article_id, content) values((select user_id from users where username = $1), $2, $3)", username, articleid, content)
			if err != nil {
				fmt.Println("AddComment - Warning: Error inserting row into DB - " + err.Error())
			}
		} else {
			fmt.Println("User " + username + " has reached comment limit")
		}
	}
}

func DeleteComment(commentid string, username string) {
	_, err := db.Exec("DELETE FROM comments WHERE user_id = (select user_id from users where username = $1) AND comment_id = $2", username, commentid)
	if err != nil {
		fmt.Println("DeleteCOmment - Warning: Error deleting row from DB - " + err.Error())
	}
}

func DeleteCommentForce(commentid string) {
	_, err := db.Exec("DELETE FROM comments WHERE comment_id = $1", commentid)
	if err != nil {
		fmt.Println("DeleteCOmment - Warning: Error deleting row from DB - " + err.Error())
	}
}










