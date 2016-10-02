package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/mvdan/xurls"
)

func getFeedForTwitterUser(c *gin.Context, userID string) (string, error) {
	twitterUser, err := getTwitterUserFromDB(userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

	timeline := getTwitterHomeTimeline(twitterUser, url.Values{})
	updateLastRetrievalInDB(c, userID)
	now := time.Now()
	feed := &feeds.Feed{
		Title:       twitterUser.screenName + " Twitter feed",
		Link:        &feeds.Link{Href: "http://twitter.com"},
		Description: "Home Twitter feed for " + twitterUser.screenName,
		Author:      &feeds.Author{Name: "TwiSSR", Email: "info@twissr.com"},
		Created:     now,
	}

	for _, e := range timeline {

		// Basics
		titleText := removeURLs(e.Text)
		createdAtTime, _ := e.CreatedAtTime()
		linkToTweet := "https://twitter.com/" + e.User.ScreenName + "/status/" + e.IdStr
		//title := e.User.Name + " @ " + e.CreatedAt + " | "

		// Images
		img := e.Entities.Media
		imageTag := ""
		if len(img) > 0 {
			//log.Printf("img %v", img[0].Media_url_https)
			imageTag = "<p><img src='" + img[0].Media_url_https + "' /></p>"
		}

		// Links
		links := "<p>"
		for _, e := range e.Entities.Urls {
			//log.Printf("link: %v", e)
			links = links + "<a href='" + e.Url + "'>" + e.Display_url + "</a> | "
		}
		links = links + "</p>"

		// Feed
		feedItem := &feeds.Item{
			Title:       e.User.Name + ": " + titleText,
			Link:        &feeds.Link{Href: linkToTweet},
			Description: e.Text + links + imageTag,
			Author:      &feeds.Author{Name: e.User.Name},
			Created:     createdAtTime,
		}
		feed.Items = append(feed.Items, feedItem)
	}

	return feed.ToRss()
}

func removeURLs(text string) string {

	result := text
	urls := xurls.Relaxed.FindAllString(text, -1)

	for _, url := range urls {
		result = strings.Replace(result, url, "", -1)
	}

	return result

}
