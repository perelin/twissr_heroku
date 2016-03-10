package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
)

type TwitterUser struct {
	userID, screenName, oauthToken, oauthTokenSecret, createDate string
}

func initAnaconda() {
	anaconda.SetConsumerKey("MBjAJxGxtqI1cysRE8686WPKE")
	anaconda.SetConsumerSecret("0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")
}

func getTwitterAuthBasics(c *gin.Context) (string, *oauth.Credentials) {
	authURL, creds, err := anaconda.AuthorizationURL("http://" + c.Request.Host + "/tc")
	if err != nil {
		http.Error(c.Writer, "Error getting Twitter Creds + Url, "+err.Error(), 500)
		log.Fatalf("Error getting Twitter Creds + Url: %q", err)
	}
	//log.Printf("Url: %s", authURL)
	//log.Printf("Oauth Token: %s", creds.Token)
	//log.Printf("Oauth Secret: %s", creds.Secret)
	//log.Printf("Url: %v creds: %v", authUrl, creds)
	return authURL, creds
}

func verifyTwitterOAuthCreds(c *gin.Context, tempCreds *oauth.Credentials, verifier string) url.Values {
	_, values, err := anaconda.GetCredentials(tempCreds, verifier)
	if err != nil {
		http.Error(c.Writer, "Error getting request token, "+err.Error(), 500)
		log.Fatalf("Error verifying Twitter TempCreds: %q", err)
	}
	return values
}

func getTwitterHomeTimeline(user TwitterUser, v url.Values) (timeline []anaconda.Tweet) {
	api := anaconda.NewTwitterApi(user.oauthToken, user.oauthTokenSecret)
	timeline, err := api.GetHomeTimeline(v)
	if err != nil {
		//http.Error(c.Writer, "Couldn´t get twitter API, "+err.Error(), 500)
		log.Fatalf("Couldn´t get twitter API: %q", err)
	}
	//log.Printf("Timeline: %v", timeline)
	return timeline
}
