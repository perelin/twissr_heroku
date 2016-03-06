package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
    "github.com/ChimeraCoder/anaconda"
    "github.com/gorilla/sessions"
    "github.com/garyburd/go-oauth/oauth"
    "encoding/gob"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))


func twitterLoginHandler(c *gin.Context) {
	// , w http.ResponseWriter, r *http.Request
    log.Printf("Host %v", c.Request.Host)
    log.Printf("RequestURI %v", c.Request.RequestURI)
    log.Printf("Twitter login page request %s", c.Request.URL.Path)


    authUrl, creds, _ := anaconda.AuthorizationURL("http://" + c.Request.Host + "/tc")
    log.Printf("Url: %s", authUrl)
    log.Printf("Oauth Token: %s", creds.Token)
    log.Printf("Oauth Secret: %s", creds.Secret)

    log.Printf("Url: %v creds: %v", authUrl, creds)

    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.
    session, err := store.Get(c.Request, "twitterTempCreds")
    if err != nil {
        log.Printf("Error: %s", err)
        return
    }

    log.Printf("Session values: %v", session.Values)


    // Set some session values.
    session.Values["tempCreds"] = creds
    session.Values["Token"] = creds.Token
    session.Values["Secret"] = creds.Secret
    // Save it before we write to the response/return from the handler.
    session.Save(c.Request, c.Writer)

    reconstructedCreds := session.Values["tempCreds"]

    log.Printf("creds: %v", creds)
    log.Printf("reconstructedCreds: %v", reconstructedCreds)

	c.HTML(http.StatusOK, "twitterlogin.tmpl.html", gin.H{"tlurl": authUrl})
	/*
    templates.ExecuteTemplate(c.Writer,"twitterlogin.tpl.html", authUrl)
	*/
}

func twitterCallbackHandler(c *gin.Context) {
    log.Printf("Twitter callback page request %s", c.Request.URL.Path)
    //log.Printf("Request %v", r)
    //func GetCredentials(tempCred *oauth.Credentials, verifier string) (*oauth.Credentials, url.Values, error)
    //creds, values, err = GetCredentials()

    err := c.Request.ParseForm()
    if err != nil {
        http.Error(c.Writer, err.Error(), 500)
        return
    }
    verifier := c.Request.Form["oauth_verifier"][0]

    session, err := store.Get(c.Request, "twitterTempCreds")
    if err != nil {
        http.Error(c.Writer, err.Error(), 500)
        return
    }
    reconstructedCreds := session.Values["tempCreds"]

    log.Printf("reconstructedCreds: %v", reconstructedCreds)

    creds, values, err := anaconda.GetCredentials(reconstructedCreds.(*oauth.Credentials), verifier)
    if err != nil {
        http.Error(c.Writer, "Error getting request token, "+err.Error(), 500)
        return
    }

    log.Printf("creds: %v", creds)
    log.Printf("values: %v", values)

	c.HTML(http.StatusOK, "twittercallback.tmpl.html", values)
    //templates.ExecuteTemplate(c.Writer,"twittercallback.tpl.html", values)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Printf("$PORT is not set, using default port (8080)")
		port = "8080"
	}

	/*
	twitterConsumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
    twitterConsumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")

	anaconda.SetConsumerKey(twitterConsumerKey)
    anaconda.SetConsumerSecret(twitterConsumerSecret)
	*/
	anaconda.SetConsumerKey("MBjAJxGxtqI1cysRE8686WPKE")
	anaconda.SetConsumerSecret("0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")

	gob.Register(&oauth.Credentials{})

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/tl", twitterLoginHandler)
	router.GET("/tc", twitterCallbackHandler)

	/*
	router.GET("/tl", func(c *gin.Context) {
		c.HTML(http.StatusOK, "twitterlogin.tmpl.html", nil)
	})*/


	router.Run(":" + port)
	//router.Run(":" + "8080")
}
