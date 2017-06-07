package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/tradesim/model"
	"github.com/tradesim/services"
	"net/http"
	"time"
	"strings"
	"strconv"
	"fmt"
	"net/url"
)

func PingHandler(c *gin.Context) {
	c.Writer.Write([]byte("Pong\n"))
}

func CreateUser(c *gin.Context) {
	var user model.NewUser
	if c.BindJSON(&user) == nil {
		user.Username = strings.ToLower(user.Username)
		user.Email = strings.ToLower(user.Email)
		if services.UserExists(user.Username) {
			c.JSON(http.StatusConflict, gin.H{"error":true,"message":"Username already used"})
		} else if services.EmailExists(user.Email) {
			c.JSON(http.StatusConflict, gin.H{"error":true,"message":"Email already used"})
		} else {
			err := services.AddUser(user.Email, user.Username, user.Password)
			if err != nil {
				panic(err)
			} else {
				services.ResetBalance(user.Username);
				services.SendAccountVerificationEmail(user.Username, user.Email)
				c.JSON(http.StatusOK, gin.H{"error":false})
			}
		}
	}
}

func GetMe(c *gin.Context) {
	var user model.User
	user, err := services.GetUserByUsername(GetUsernameFromContext(c))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"Username":user.Username,"Email":user.Email});
}

func UpdatePassword(c *gin.Context) {
	var update model.UpdatePassword
	if c.BindJSON(&update) == nil {
		err := services.UpdatePassword(update, GetUsernameFromContext(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":true,"message":err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"error":false})
		}
	}
}

func Login(c *gin.Context) {
	var login model.Login
	if c.BindJSON(&login) == nil {
		login.Username = strings.ToLower(login.Username)
		message, err := services.Login(login)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error":true,"message":err.Error()})
		} else {
			c.Writer.Header().Add("Set-Cookie", "tradesim=" + message + ";path=/api")
			c.JSON(http.StatusOK, gin.H{"error":false})
		}
	}
}

func Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie("tradesim")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":"Cannot perform logout because no login session was found"})
	} else {
		services.DeleteSessionByKey(cookie.Value)
		c.JSON(http.StatusOK, gin.H{"error":false,"message":"You have been logged out successfully"})
	}
}

func VerifyCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("tradesim")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":"No session cookie provided"})
	} else {
		sesh, err := services.GetSession(cookie.Value)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error":true,"message":"Provided cookie not found. It may have expired."})
		} else {
			now := int(time.Now().UnixNano() / int64(time.Millisecond))
			c.JSON(http.StatusOK, gin.H{"error":false,"message":"Session found for user " + sesh.Username + " expiring in " + strconv.Itoa((sesh.Expiration - now) / 1000) + " seconds"})
		}
	}
}

func GetBalance(c *gin.Context) {
	var bal map[string]interface{}
	bal = services.GetBalance(GetUsernameFromContext(c))
	c.JSON(http.StatusOK, bal)
}

func GetTransactions(c *gin.Context) {
	var transactions []model.TransactionDetail
	transactions = services.GetTransactions(GetUsernameFromContext(c))
	c.JSON(http.StatusOK, transactions);
}

func Buy(c *gin.Context) {
	var transaction model.Transaction
	if c.BindJSON(&transaction) == nil {
		username := GetUsernameFromContext(c)
		var balances map[string]interface{}
		balances = services.GetBalance(username)
		rateObj := GetRate(transaction.Symbol)
		rate := rateObj.Ask
		if balances["USD"].(float64) < (rate * transaction.Quantity) + transaction.Fee {
			message := fmt.Sprintf("Insufficient funds (requires %f)", (transaction.Quantity * rate) + transaction.Fee)
			c.JSON(http.StatusForbidden, gin.H{"error":true,"message": message})
		} else {
			balances["USD"] = balances["USD"].(float64) - ((rate * transaction.Quantity) + transaction.Fee)
			if balances[transaction.Symbol] == nil {
				balances[transaction.Symbol] = transaction.Quantity
			} else {
				balances[transaction.Symbol] = balances[transaction.Symbol].(float64) + transaction.Quantity
			}
			err := services.UpdateBalance(username, balances)
			if err != nil {
				panic(err)
			} else {
				message := fmt.Sprintf("Traded %f USD for %f %s at a rate of %f with a transaction total cost of %f", (transaction.Quantity * rate), transaction.Quantity, transaction.Symbol, rate, transaction.Fee)
				err = services.AddTransaction(username, transaction.Symbol, "BUY", transaction.Quantity, rate, transaction.Fee)
				if err != nil {
					panic(err)
				}
				c.JSON(http.StatusOK, gin.H{"error":false,"message": message})
			}
		}
	}
}

func Sell(c *gin.Context) {
	var transaction model.Transaction
	if c.BindJSON(&transaction) == nil {
		username := GetUsernameFromContext(c)
		var balances map[string]interface{}
		balances = services.GetBalance(username)
		if _, ok := balances[transaction.Symbol]; !ok {
			balances[transaction.Symbol] = 0.0
		}
		rateObj := GetRate(transaction.Symbol)
		rate := rateObj.Bid
		if balances[transaction.Symbol].(float64) < transaction.Quantity {
			message := fmt.Sprintf("Cannot sell %f %s because you only have %f", transaction.Quantity, transaction.Symbol, balances[transaction.Symbol])
			c.JSON(http.StatusForbidden, gin.H{"error":true,"message": message})
		} else {
			balances[transaction.Symbol] = balances[transaction.Symbol].(float64) - transaction.Quantity
			balances["USD"] = balances["USD"].(float64) + ((transaction.Quantity * rate) - transaction.Fee)
			err := services.UpdateBalance(username, balances)
			if err != nil {
				panic(err)
			} else {
				message := fmt.Sprintf("Traded %f %s for %f USD at a rate of %f with a transaction total cost of %f", transaction.Quantity, transaction.Symbol, (transaction.Quantity * rate), rate, transaction.Fee)
				err = services.AddTransaction(username, transaction.Symbol, "SELL", transaction.Quantity, rate, transaction.Fee)
				if err != nil {
					panic(err)
				}
				c.JSON(http.StatusOK, gin.H{"error":false,"message":message})
			}
		}
	}
}

func GetRate(symbol string) model.Rate { // TODO: Handle case where symbol does not exist, or is USD
	if rate, ok := services.CurrentRates[symbol]; ok {
		return rate
	}
	if (symbol == "BTC") {
		return services.GetBitcoinRate()[0]
	} else {
		symbolArr := make([]string, 0)
		symbolArr = append(symbolArr, symbol)
		return services.RetrieveRates(symbolArr)[0]
	}
}

func GetBTCPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"price":services.GetBitcoinPriceUSD()})
}

func GetUsernameFromContext(c *gin.Context) string {
	cookie, _ := c.Request.Cookie("tradesim")
	var sesh model.Session
	sesh, _ = services.GetSession(cookie.Value)
	return sesh.Username
}

func ResetBalance(c *gin.Context) {
	username := GetUsernameFromContext(c)
	err := services.ResetBalance(username)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"error":false,"message":services.InitialBalance})
}

func VerifyEmail(c *gin.Context) {
	username := services.DecodeToken(c.Param("token"))
	unencoded, err := url.PathUnescape(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":"Invalid link"})
	} else {
		err = services.UpdateUserLastLogin(unencoded)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error":true,"message":err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"error":false,"message":"User has been successfully verified!"})
		}
	}
	
}

func GetAssetPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	fmt.Printf("\nGetting rate object for symbol %s", symbol)
	var rate model.Rate
	if val, ok := services.CurrentRates[symbol]; ok {
		rate = val
	} else {
		if (strings.ToUpper(symbol) == "BTC") {
			rate = services.GetBitcoinRate()[0]
		} else {
			symbolArr := make([]string, 0)
			symbolArr = append(symbolArr, symbol)
			rates := services.RetrieveRates(symbolArr)
			if len(rates) == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error":true,"message":"Sybol not found"})
				return
			} else {
				rate = rates[0]
			}
		}
	}

	c.JSON(http.StatusOK, rate)
}

func SendResetPasswordEmail(c *gin.Context) {
	username := c.Param("username")
	fmt.Println("Attempting to send reset password email for user " + username)
	user, err := services.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":true,"message":"Username not found"})
	} else {
		services.SendResetPasswordLink(username, user.Email, services.CreateToken(username))
		c.JSON(http.StatusOK, gin.H{"error":false,"message":"Please check your email for instructions to reset your password."})
	}
}

func ResetPassword(c *gin.Context) {
	username := services.DecodeToken(c.Param("token"))
	fmt.Printf("\n\nAttempting to reset password for user %s\n", username)
	// get user's email
	user, err := services.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":"Token not found"});
	} else {
		// generate random password
		newPass := services.RandomPassword()
		// update password in db
		err := services.UpdatePasswordForce(newPass, username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":err.Error()})
		}
		// send email with new password
		services.SendNewPasswordEmail(username, user.Email, newPass)
		c.JSON(http.StatusOK, gin.H{"error":false,"message":"Password has been sucessfully reset"})
	}
}

// Account Value Handlers

func GetMyValue(c *gin.Context) {
	val := GetAccountValueByUsername(GetUsernameFromContext(c))
	c.JSON(http.StatusOK, gin.H{"AccountValueUSD":val})
}

func GetAccountValue(c *gin.Context) {
	username := c.Param("username")
	totalValue := GetAccountValueByUsername(username)
	c.JSON(http.StatusOK, gin.H{"AccountValueUSD":totalValue})
}

func GetAccountValueByUsername(username string) float64 {
	balance := services.GetBalance(username)
	var totalValue float64
	totalValue = 0.0

	// get all symbols besides USD and BTC
	missingSymbols := make(map[string]bool)
	for symbol, _ := range balance {
		if symbol != "BTC" && symbol != "USD" {
			if _, ok := services.CurrentRates[symbol]; !ok {
			    missingSymbols[symbol] = true
			}
		}
	}

	var missArray []string
	for key, _ := range missingSymbols {
		if key != "BTC" && key != "USD" {
			missArray = append(missArray, key)
		}
	}

	// update  rates
	services.UpdateRates(missArray)

	// call stock api for those symbols
	for symbol, quantity := range balance {
		if symbol == "USD" {
			totalValue += quantity.(float64)
		} else {
			totalValue += services.CurrentRates[symbol].Price * quantity.(float64)
		}
	}

	return totalValue
}

func GetAllUserBalances(c *gin.Context) {
	var users []model.User
	var values map[string]float64 // username => account value in usd
	values = make(map[string]float64)
	users = services.GetAllUsers()
	missingSymbols := make(map[string]bool)

	// calculate balance for rates we have already found
	for _, user := range users {
		balances := services.GetBalance(user.Username)
		for symbol, _ := range balances {
			if _, ok := services.CurrentRates[symbol]; !ok {
			    missingSymbols[symbol] = true
			}
		}
	}

	var missArray []string
	for key, _ := range missingSymbols {
		if key != "BTC" && key != "USD" {
			missArray = append(missArray, key)
		}
	}

	// get rates we haven't found
	services.UpdateRates(missArray)

	// calculate balance including rates we hadn't previously found
	for _, user := range users {
		var tempVal float64
		balances := services.GetBalance(user.Username)
		for symbol, quantity := range balances {
			if symbol == "USD" {
				tempVal += quantity.(float64)
			} else {
				tempVal += (services.CurrentRates[symbol].Price * quantity.(float64))
			}
		}
		values[user.Username] = tempVal
	}
	c.JSON(http.StatusOK, values)
}

// ARTICLES

func GetArticle(c *gin.Context) {
	var article model.Article
	articleId := c.Param("articleid")
	article = services.GetArticle(articleId)
	if article.Title == "" {
		c.JSON(http.StatusNotFound, gin.H{"error":true,"message":"Article with provided id (" + articleId + ") could not be found"})
	} else {
		c.JSON(http.StatusOK, article)
	}
}

func GetRecentArticles(c *gin.Context) {
	var articles []model.Article
	limit := c.Query("limit")
	fmt.Println("Limit: " + limit)
	if limit == "" {
		limit = "0"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":"invalid value provided for limit parameter"})
	} else {
		articles = services.GetArticles(limitInt)
		c.JSON(http.StatusOK, articles)
	}
}

// COMMENTS

func GetCommentsForArticle(c *gin.Context) {
	articleid := c.Param("articleid")
	var comments []model.Comment
	comments = services.GetCommentsForArticle(articleid)
	c.JSON(http.StatusOK, comments)
}

func AddComment(c *gin.Context) {
	var comment model.Comment
	if c.BindJSON(&comment) == nil {
		articleid := c.Param("articleid")
		username := GetUsernameFromContext(c)
		fmt.Printf("\nAdding comment for user %s on article with id %s", username, articleid)
		fmt.Printf("\nArticle content: %s", comment.Content)
		services.AddComment(articleid, username, comment.Content)
		var comments []model.Comment
		comments = services.GetCommentsForArticle(articleid)
		c.JSON(http.StatusOK, comments)
	}
	
}

func UpdateComment(c *gin.Context) {
	
}

func DeleteComment(c *gin.Context) {
	commentid := c.Param("commentid")
	username := GetUsernameFromContext(c)
	fmt.Printf("\nDeleting comment with id %s belonging to user %s", commentid, username)
	if username == "k128kevin" {
		services.DeleteCommentForce(commentid)
	} else {
		services.DeleteComment(commentid, username)
	}
	c.JSON(http.StatusOK, gin.H{"error":false,"message":"Successfully deleted comment"})
}









