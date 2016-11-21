package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkRequestForm(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		log.Printf("checkRequestForm failed")
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}
