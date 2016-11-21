package main

import (
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
		//http.Error(c.Writer, err.Error(), 404)
		//c.AbortWithStatus(404)
		return "", err
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

		// Text
		text := linkURLs(e.Text)

		// Feed
		feedItem := &feeds.Item{}
		if e.Text == auth_failed {
			feedItem = &feeds.Item{
				Title:       "Twitter authentication failed!",
				Link:        &feeds.Link{Href: "http://www.twissr.com"},
				Description: "Please visit <a href='http://www.twissr.com'>http://www.twissr.com</a> to reauthenticate.",
				Author:      &feeds.Author{Name: "TwiSSR"},
				Created:     createdAtTime,
			}
		} else {
			feedItem = &feeds.Item{
				Title:       e.User.Name + ": " + titleText,
				Link:        &feeds.Link{Href: linkToTweet},
				Description: text + links + imageTag,
				Author:      &feeds.Author{Name: e.User.Name},
				Created:     createdAtTime,
			}
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

func linkURLs(text string) string {
	result := text
	urls := xurls.Relaxed.FindAllString(text, -1)
	for _, url := range urls {
		result = strings.Replace(result, url,
			"<a href='"+url+"' target='_blank'>"+url+"</a>", -1)
	}
	return result
}
