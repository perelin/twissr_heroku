package main

import (
	"log"
	"net/http"
	"os"
    "database/sql"
	//"reflect"

	"github.com/gin-gonic/gin"
    "github.com/ChimeraCoder/anaconda"
    "github.com/gorilla/sessions"
    "github.com/garyburd/go-oauth/oauth"
    "encoding/gob"
    _ "github.com/lib/pq"

)

var (
	store = sessions.NewCookieStore([]byte("something-very-secret"))
	//db *sql.DB = nil
	db = initDB()
	router = gin.New()
)

func getTwitterAuthBasics(c *gin.Context) (string, *oauth.Credentials) {
	authUrl, creds, err := anaconda.AuthorizationURL("http://" + c.Request.Host + "/tc")
	if err != nil {
		http.Error(c.Writer, "Error getting Twitter Creds + Url, "+err.Error(), 500)
        log.Fatalf("Error getting Twitter Creds + Url: %q", err)
    }
	//log.Printf("Url: %s", authUrl)
    //log.Printf("Oauth Token: %s", creds.Token)
    //log.Printf("Oauth Secret: %s", creds.Secret)
    //log.Printf("Url: %v creds: %v", authUrl, creds)
	return authUrl, creds
}

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

func twitterLoginHandler(c *gin.Context) {
	authURL, creds := getTwitterAuthBasics(c)
	setTwitterTempCreds(c, creds)
	c.HTML(http.StatusOK, "twitterlogin.tmpl.html", gin.H{"tlurl": authURL})
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
    //log.Printf("values: %s, %v", reflect.TypeOf(values["user_id"][0]), len(values["user_id"]))

	// Error writing to db, sql: converting Exec argument #0's type: unsupported type []string, a slice
	_, err = db.Exec("INSERT INTO twitter_users VALUES ($1, $2, $3, $4)",
							values["user_id"][0],
							values["screen_name"][0],
							values["oauth_token"][0],
							values["oauth_token_secret"][0])
	if err != nil {
        //c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
		http.Error(c.Writer, "Error writing to db, "+err.Error(), 500)
		//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
        return
    }
	/*
*/
	c.HTML(http.StatusOK, "twittercallback.tmpl.html", values)
    //templates.ExecuteTemplate(c.Writer,"twittercallback.tpl.html", values)
}

func initDB() *sql.DB {
	log.Printf("DATABASE_URL: %s", os.Getenv("DATABASE_URL"))
	db, errd := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if errd != nil {
        log.Fatalf("Error opening database: %q", errd)
    }
	return db
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("$PORT is not set, using default port (8080)")
		port = "5000"
	}
	return port
}

func initAnaconda() {
	anaconda.SetConsumerKey("MBjAJxGxtqI1cysRE8686WPKE")
	anaconda.SetConsumerSecret("0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")
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
	/*
	router.GET("/tl", func(c *gin.Context) {
		c.HTML(http.StatusOK, "twitterlogin.tmpl.html", nil)
	})*/
	router.Run(":" + getPort())
}

func main() {

	printsup()

	//db = initDB()
	initAnaconda()

	gob.Register(&oauth.Credentials{})

	initRouter()
}
