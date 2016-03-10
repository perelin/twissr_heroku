package main

import (
	"log"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func setTwitterTempCreds(c *gin.Context, creds *oauth.Credentials) {

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(c.Request, "twitterTempCreds")
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	//log.Printf("Session values: %v", session.Values)

	// Set some session values.
	session.Values["tempCreds"] = creds
	session.Values["Token"] = creds.Token
	session.Values["Secret"] = creds.Secret
	// Save it before we write to the response/return from the handler.
	session.Save(c.Request, c.Writer)

	//reconstructedCreds := session.Values["tempCreds"]
	//log.Printf("creds: %v", creds)
	//log.Printf("reconstructedCreds: %v", reconstructedCreds)
}

func getTwitterTempCreds(c *gin.Context) *oauth.Credentials {
	session, err := store.Get(c.Request, "twitterTempCreds")
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		log.Fatalf("Error getting Twitter TempCreds from Session: %q", err)
	}
	reconstructedCreds := session.Values["tempCreds"]
	log.Printf("reconstructedCreds: %v", reconstructedCreds)
	return reconstructedCreds.(*oauth.Credentials)
}
