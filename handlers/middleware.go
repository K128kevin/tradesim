package handlers

import(
	"github.com/tradesim/services"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func VerifyCookieMW(c *gin.Context) {
	cookie, err := c.Request.Cookie("tradesim")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":true,"message":"access token cookie required"})
		c.Abort()
		return
	}
	_, err = services.GetSession(cookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":true,"message":err.Error()})
		c.Abort()
		return
	}
}