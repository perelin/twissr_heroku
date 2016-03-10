package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
)

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

	getTwitterHomeTimeline(twitterUser, url.Values{})

	//saveTwitterCredsToDB(c, values)
	c.HTML(http.StatusOK, "twittercallback.tmpl.html", values)
}

func main() {
	initAnaconda()
	gob.Register(&oauth.Credentials{})
	initRouter()
}
