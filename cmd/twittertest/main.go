package main

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

func main() {

	//var api *anaconda.TwitterApi

	api := anaconda.NewTwitterApiWithCredentials("	7559392-6AMV7putpKu6MAhh0mbYXGBmoLZbKQOUZWOUC0pWm5", "	fTOwcpzYo6zt37WF2wYEKKOZJf2hjQDBKgQp49n9N9Fh2", "MBjAJxGxtqI1cysRE8686WPKE", "0eXlKiGYf5TpzOo67y7H0bCzJETgeoRRT1YY4iCqC85wdkZVbZ")

	authURL, _, err := api.AuthorizationURL("http://localhost:5000/tc")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(authURL)
}
