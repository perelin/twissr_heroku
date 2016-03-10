package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.New()
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("$PORT is not set, using default port (8080)")
		port = "5000"
	}
	return port
}

func initRouter() {
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/tl", twitterLoginHandler)
	router.GET("/tc", twitterCallbackHandler)
	router.GET("/hometimeline/:userID", twitterHomeTimelineRSSHandler)
	/*
		router.GET("/tl", func(c *gin.Context) {
			c.HTML(http.StatusOK, "twitterlogin.tmpl.html", nil)
		})*/
	router.Run(":" + getPort())
}
