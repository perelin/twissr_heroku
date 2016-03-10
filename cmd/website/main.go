package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func twitterHomeTimelineRSSHandler(c *gin.Context) {
	userID := c.Param("userID")
	//c.String(http.StatusOK, "UserID: %s", userID)
	twitterUser, err := getTwitterUserFromDB(userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
	timeline := getTwitterHomeTimeline(twitterUser, url.Values{})
	now := time.Now()
	feed := &feeds.Feed{
		Title:       twitterUser.screenName + " home timeline (provided by TwiSSR)",
		Link:        &feeds.Link{Href: "http://twitter.com"},
		Description: "Home timeline Twitter feed for " + twitterUser.screenName,
		Author:      &feeds.Author{Name: "TwiSSR", Email: "info@twissr.com"},
		Created:     now,
	}

	for _, e := range timeline {
		//log.Printf("%v", e.crea)
		createdAtTime, _ := e.CreatedAtTime()
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       e.User.Name,
			Link:        &feeds.Link{Href: e.Source},
			Description: e.Text,
			Author:      &feeds.Author{Name: e.User.Name, Email: e.User.Name},
			Created:     createdAtTime,
		})
	}

	rss, err := feed.ToRss()
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

	getTwitterHomeTimeline(twitterUser, url.Values{})

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "jmoiron.net blog",
		Link:        &feeds.Link{Href: "http://jmoiron.net/blog"},
		Description: "discussion about tech, footie, photos",
		Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	//saveTwitterCredsToDB(c, values)
	//c.String(http.StatusOK, "%s", feed)

	c.XML(http.StatusOK, rss)
	//c.HTML(http.StatusOK, "twittercallback.tmpl.html", values)
}

func main() {
	initAnaconda()
	gob.Register(&oauth.Credentials{})
	initRouter()
}
