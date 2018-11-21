package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/garyburd/go-oauth/oauth"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("working dir:")
	fmt.Println(dir)

	initAnaconda()
	gob.Register(&oauth.Credentials{})
	initRouter()
}
