package main

import (
	"encoding/gob"

	"github.com/garyburd/go-oauth/oauth"
)

func main() {
	initAnaconda()
	gob.Register(&oauth.Credentials{})
	initRouter()
}
