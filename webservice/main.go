package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static/", "./static/")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Ping Pong",
		})
	})

	r.GET("/call-ping", func(c *gin.Context) {
		resp, err := http.Get("http://pingservice:8081/ping")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error calling pingservice: %v", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.String(http.StatusOK, string(body))
	})

	r.GET("/call-pong", func(c *gin.Context) {
		resp, err := http.Get("http://pongservice:8082/pong")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error calling pongservice: %v", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.String(http.StatusOK, string(body))
	})

	r.Run(":8080") // webservice na porta 8080
}
