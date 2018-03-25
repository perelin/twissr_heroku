package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/gin"
)

// TwitterUser needs to have a comment...
type TwitterUser struct {
	userID,
	screenName,
	oauthToken,
	oauthTokenSecret,
	createDate,
	lastRetrievalDate string
}

func initAnaconda() {
	anaconda.SetConsumerKey("MBjAJxGxtqI1cysRE8686WPKE")
	anaconda.SetConsumerSecret("0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")
}

func getApiWithCredentials() *anaconda.TwitterApi {

	api := anaconda.NewTwitterApiWithCredentials("7559392-6AMV7putpKu6MAhh0mbYXGBmoLZbKQOUZWOUC0pWm5", "	fTOwcpzYo6zt37WF2wYEKKOZJf2hjQDBKgQp49n9N9Fh2", "MBjAJxGxtqI1cysRE8686WPKE", "0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")

	return api
}

func getTwitterAuthBasics(c *gin.Context) (string, *oauth.Credentials) {
	spew.Dump(c.Request.Host)

	var api = getApiWithCredentials()

	authURL, creds, err := api.AuthorizationURL("http://" + c.Request.Host + "/tc")

	spew.Dump("######################################")

	if err != nil {
		http.Error(c.Writer, "Error getting Twitter Creds + Url, "+err.Error(), 500)
		log.Fatalf("Error getting Twitter Creds + Url: %q", err)
	}
	//log.Printf("Url: %s", authURL)
	//log.Printf("Oauth Token: %s", creds.Token)
	//log.Printf("Oauth Secret: %s", creds.Secret)
	//log.Printf("Url: %v creds: %v", authUrl, creds)
	return authURL, creds
}

func verifyTwitterOAuthCreds(c *gin.Context, tempCreds *oauth.Credentials, verifier string) url.Values {

	var api = getApiWithCredentials()

	_, values, err := api.GetCredentials(tempCreds, verifier)
	if err != nil {
		http.Error(c.Writer, "Error getting request token, "+err.Error(), 500)
		//c.Redirect(http.StatusFound, "/")
		//log.Fatalf("Error verifying Twitter TempCreds: %q", err)
	}
	return values
}

func getTwitterHomeTimeline(user TwitterUser, v url.Values) (timeline []anaconda.Tweet) {

	spew.Dump("### getTwitterHomeTimeline")
	spew.Dump(user)

	api := anaconda.NewTwitterApi(user.oauthToken, user.oauthTokenSecret)

	timeline, err := api.GetHomeTimeline(v)

	if err != nil {
		//http.Error(c.Writer, "Couldn´t get twitter API, "+err.Error(), 500)
		//log.Fatalf("Couldn´t get twitter API: %q", err)
		log.Printf("Couldn´t get twitter API: %q", err)
		timeline = append(timeline, anaconda.Tweet{Text: auth_failed})
	}
	//log.Printf("Timeline: %v", timeline)
	return timeline
}

func getTwitterLists(user TwitterUser, v url.Values) []anaconda.List {
	api := anaconda.NewTwitterApi(user.oauthToken, user.oauthTokenSecret)
	id, err := strconv.ParseInt(user.userID, 10, 64)
	if err != nil {
		panic(err)
	}
	lists, err := api.GetListsOwnedBy(id, v)
	if err != nil {
		panic(err)
	}

	return lists
}

// func getTwitterListMemebers(user TwitterUser, listID int64, v url.Values) (c UserCursor, err error) {

// 	api := anaconda.NewTwitterApi(user.oauthToken, user.oauthTokenSecret)

// 	return anaconda.UserCursor()

// 	//c, err := api.get ..... upgrade to newest anaconda version + get list members

// }
