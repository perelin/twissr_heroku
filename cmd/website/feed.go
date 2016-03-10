package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func getFeedForTwitterUser(c *gin.Context, userID string) (string, error) {
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
		linkToTweet := "https://twitter.com/" + e.User.ScreenName + "/status/" + e.IdStr
		title := e.User.Name + " @ " + e.CreatedAt
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: linkToTweet},
			Description: e.Text,
			Author:      &feeds.Author{Name: e.User.Name},
			Created:     createdAtTime,
		})
	}

	return feed.ToRss()
}
