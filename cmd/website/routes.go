package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.New()
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("$PORT is not set, using default port (5000)")
		port = "5000"
	}
	return port
}

func initRouter() {

	//router.Use(NewRelic("175defc71b45f0e19f5073564299f2c9779da80b", "TwiSSR", false))

	//config := newrelic.NewConfig("TwiSSR", "175defc71b45f0e19f5073564299f2c9779da80b")
	//App, err := newrelic.NewApplication(config)

	router.Use(gin.Logger())
	//router.LoadHTMLGlob(os.Getenv("TEMPLATE_FOLDER_PREFIX") + "templates/*")
	router.LoadHTMLGlob("templates/*")
	/*router.LoadHTMLFiles("templates/twitterlogin.tmpl.html",
	"templates/twittercallback.tmpl.html",
	"templates/h5bp/twissr.html")*/
	//router.LoadHTMLGlob("templates/h5bp/*.html")
	router.Static("/static", "static")
	router.Static("/js", "static/h5bp/js")
	router.Static("/css", "static/h5bp/css")
	router.Static("/img", "static/h5bp/img")
	router.Static("/fonts", "static/h5bp/fonts")

	/*
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl.html", nil)
		})*/

	router.GET("/", startPageHandler)
	router.GET("/tl", twitterLoginHandler)
	router.GET("/tc", twitterCallbackHandler)
	router.GET("/hometimeline/:userID", twitterHomeTimelineRSSHandler)
	router.GET("/lists/:userID", lists)
	/*
		router.GET("/tl", func(c *gin.Context) {
			c.HTML(http.StatusOK, "twitterlogin.tmpl.html", nil)
		})*/
	router.Run(":" + getPort())
}
