package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	db = initDB()
)

func initDB() *sql.DB {
	log.Printf("DATABASE_URL: %s", os.Getenv("DATABASE_URL"))
	db, errd := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	return db
}

func getTwitterUserFromDB(twitterUserID string) (TwitterUser, error) {
	// make struct for twitter user and use that, or use url value
	var twitterUser TwitterUser
	err := db.QueryRow(
		"SELECT * FROM twitter_users WHERE user_id = $1",
		twitterUserID).Scan(
		&twitterUser.userID,
		&twitterUser.screenName,
		&twitterUser.oauthToken,
		&twitterUser.oauthTokenSecret,
		&twitterUser.createDate)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Values %q", twitterUser)
	}
	return twitterUser, err
}

func saveTwitterCredsToDB(c *gin.Context, values url.Values) {
	getTwitterUserFromDB(values["user_id"][0])

	_, err := db.Exec("INSERT INTO twitter_users VALUES ($1, $2, $3, $4, now())",
		values["user_id"][0],
		values["screen_name"][0],
		values["oauth_token"][0],
		values["oauth_token_secret"][0])
	if err != nil {
		//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
		http.Error(c.Writer, "Error saving Twitter creds to DB, "+err.Error(), 500)
		//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
		log.Fatalf("Error saving Twitter creds to DB: %q", err)
	}
}
