package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkRequestForm(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}
