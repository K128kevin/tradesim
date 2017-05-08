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
	user = services.GetUserByUsername(GetUsernameFromContext(c))
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
		rate, err := GetRate(transaction.Symbol)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":err.Error()})
		} else {
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
}

func Sell(c *gin.Context) {
	var transaction model.Transaction
	if c.BindJSON(&transaction) == nil {
		username := GetUsernameFromContext(c)
		var balances map[string]interface{}
		balances = services.GetBalance(username)
		rate, err := GetRate(transaction.Symbol)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":true,"message":err.Error()})
		} else {
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
}

func GetRate(symbol string) (float64, error) {
	if (symbol == "BTC") {
		return services.GetBitcoinPriceUSD(symbol), nil
	} else {
		return services.GetStockPriceUSD(symbol)
	}
}

func GetBTCPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"price":services.GetBitcoinPriceUSD("BTC")})
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









