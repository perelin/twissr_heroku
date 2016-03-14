package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
)

func startPageHandler(c *gin.Context) {
	authURL, creds := getTwitterAuthBasics(c)
	setTwitterTempCreds(c, creds)
	c.HTML(http.StatusOK, "twissr.tmpl.html", gin.H{"authURL": authURL})
}

func twitterHomeTimelineRSSHandler(c *gin.Context) {
	rss, err := getFeedForTwitterUser(c, c.Param("userID"))
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusOK, rss)
}

func twitterLoginHandler(c *gin.Context) {
	authURL, creds := getTwitterAuthBasics(c)
	setTwitterTempCreds(c, creds)
	c.HTML(http.StatusOK, "twitterlogin.tmpl.html", gin.H{"tlurl": authURL})
}

func twitterCallbackHandler(c *gin.Context) {
	log.Printf("Twitter callback page request %s", c.Request.URL.Path)
	checkRequestForm(c)
	reconstructedTempCreds := getTwitterTempCreds(c)
	verifier := c.Request.Form["oauth_verifier"][0]
	values := verifyTwitterOAuthCreds(c, reconstructedTempCreds, verifier)

	twitterUser := TwitterUser{
		userID:           values["user_id"][0],
		screenName:       values["screen_name"][0],
		oauthToken:       values["oauth_token"][0],
		oauthTokenSecret: values["oauth_token_secret"][0]}

	log.Printf("values: %v", twitterUser)

	twitterUser, err := getTwitterUserFromDB(values["user_id"][0])
	if err != nil {
		saveTwitterCredsToDB(c, values)
	}

	content := gin.H{
		"rssURL":             "http://" + c.Request.Host + "/hometimeline/" + values["user_id"][0],
		"user_id":            values["user_id"][0],
		"screen_name":        values["screen_name"][0],
		"oauth_token":        values["oauth_token"][0],
		"oauth_token_secret": values["oauth_token_secret"][0],
	}

	c.HTML(http.StatusOK, "twittercallback.tmpl.html", content)
}

func twitterCallbackHandler2(c *gin.Context) {
	log.Printf("Twitter callback page request %s", c.Request.URL.Path)
	checkRequestForm(c)
	reconstructedTempCreds := getTwitterTempCreds(c)
	verifier := c.Request.Form["oauth_verifier"][0]
	values := verifyTwitterOAuthCreds(c, reconstructedTempCreds, verifier)

	_, err := getTwitterUserFromDB(values["user_id"][0])
	if err != nil {
		saveTwitterCredsToDB(c, values)
	}

	content := gin.H{
		"rssURL":             "http://" + c.Request.Host + "/hometimeline/" + values["user_id"][0],
		"user_id":            values["user_id"][0],
		"screen_name":        values["screen_name"][0],
		"oauth_token":        values["oauth_token"][0],
		"oauth_token_secret": values["oauth_token_secret"][0],
	}

	c.HTML(http.StatusOK, "twissr_cb.tmpl.html", content)
}

func main() {
	initAnaconda()
	gob.Register(&oauth.Credentials{})
	initRouter()
}
