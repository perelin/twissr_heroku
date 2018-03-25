package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func twitterHomeTimelineRSSHandler(c *gin.Context) {
	rss, err := getFeedForTwitterUser(c, c.Param("userID"))
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	//c.String(http.StatusOK, rss)
	c.Render(
		http.StatusOK, render.Data{
			ContentType: "application/xml",
			Data:        []byte(rss),
		})
}

func startPageHandler(c *gin.Context) {
	authURL, creds := getTwitterAuthBasics(c)
	setTwitterTempCreds(c, creds)
	c.HTML(http.StatusOK, "twissr.tmpl.html", gin.H{"authURL": authURL})
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

	log.Printf("Identified Twitter userID: %q", values["user_id"][0])

	saveTwitterCredsToDB(c, values)

	/*
		_, err := getTwitterUserFromDB(values["user_id"][0])
		if err != nil {
			saveTwitterCredsToDB(c, values)
		} else {
			log.Printf("We could update that user")
		}*/

	content := gin.H{
		"rssURL":             "http://" + c.Request.Host + "/hometimeline/" + values["user_id"][0],
		"user_id":            values["user_id"][0],
		"screen_name":        values["screen_name"][0],
		"oauth_token":        values["oauth_token"][0],
		"oauth_token_secret": values["oauth_token_secret"][0],
	}

	c.HTML(http.StatusOK, "twissr_cb.tmpl.html", content)
}

func lists(c *gin.Context) {

	userID := c.Param("userID")

	twitterUser, err := getTwitterUserFromDB(userID)

	if err != nil {
		c.String(http.StatusOK, "Couldn´t find user in DB")
	}

	lists := getTwitterLists(twitterUser, url.Values{})

	for _, list := range lists {
		log.Printf("%#v", list.Id)
		//spew.Dump(list)

	}

	c.String(http.StatusOK, twitterUser.screenName)

}
