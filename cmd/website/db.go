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
	db_handle, errd := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	return db_handle
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
		&twitterUser.createDate,
		&twitterUser.lastRetrievalDate)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Printf("Unknown error (is postgress running?): %q", err)
	default:
		fmt.Printf("Values %q", twitterUser)
	}
	return twitterUser, err
}

func saveTwitterCredsToDB(c *gin.Context, values url.Values) {

	twitterUser, err := getTwitterUserFromDB(values["user_id"][0])
	if err != nil {
		_, err := db.Exec("INSERT INTO twitter_users VALUES ($1, $2, $3, $4, now(), now())",
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
		log.Printf("Initialy saved twitter user creds to db: %q", values["screen_name"][0])
	} else {
		_, err := db.Exec("UPDATE twitter_users SET oauth_token = $1, oauth_token_secret = $2 WHERE user_id = $3",
			values["oauth_token"][0],
			values["oauth_token_secret"][0],
			twitterUser.userID)
		if err != nil {
			//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
			http.Error(c.Writer, "Error updateing Twitter creds to DB, "+err.Error(), 500)
			//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
			log.Fatalf("Error updating Twitter creds to DB: %q", err)
		}
		log.Printf("updated twitter user creds to db: %q", twitterUser.screenName)
	}
}

func updateLastRetrievalInDB(c *gin.Context, userID string) {

	_, err := db.Exec("UPDATE twitter_users SET last_retrieval_date = now() WHERE user_id = $1;",
		userID)
	/*
		getTwitterUserFromDB(values["user_id"][0])

		_, err := db.Exec("INSERT INTO twitter_users VALUES ($1, $2, $3, $4, now())",
			values["user_id"][0],
			values["screen_name"][0],
			values["oauth_token"][0],
			values["oauth_token_secret"][0])
	*/

	if err != nil {
		//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
		http.Error(c.Writer, "Error saving last retrieval date to DB, "+err.Error(), 500)
		//c.String(http.StatusInternalServerError, log.Printf("Error updateing DB tick: %v", err))
		log.Fatalf("Error saving last retrieval date to DB: %q", err)
	}
}
