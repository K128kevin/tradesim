package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/tradesim/handlers"
	"github.com/tradesim/services"
)

func main() {

	services.Initialize()
    router := gin.Default()

    router.GET("/ping", handlers.PingHandler)
    router.GET("/api/verifyEmail/:token", handlers.VerifyEmail)

	/////////////////////
    // User Management //
    /////////////////////

    v2 := router.Group("/api/users")
    {
    	v2.POST("", handlers.CreateUser)
	    v2.POST("/login", handlers.Login)
	    v2.POST("/logout", handlers.Logout)

	    v2.POST("/verifyCookie", handlers.VerifyCookie)
    }
    
    v3 := router.Group("/api/users/me")
    v3.Use(handlers.VerifyCookieMW)
    {
        v3.GET("", handlers.GetMe)
        v3.PATCH("/password", handlers.UpdatePassword)
    }

    ///////////////
    // Trade Sim //
    ///////////////

    v1 := router.Group("/api/tradesim")
    v1.Use(handlers.VerifyCookieMW)
    {
    	v1.GET("/btcrate", handlers.GetBTCPrice)
    	v1.GET("/balance", handlers.GetBalance)
	    v1.GET("/transactions", handlers.GetTransactions)
	    v1.POST("/transactions/buy", handlers.Buy)
	    v1.POST("/transactions/sell", handlers.Sell)
	    v1.POST("/balance/reset", handlers.ResetBalance)
    }

    router.Run()

}