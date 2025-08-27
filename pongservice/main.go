package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// rota hello
	r.GET("/pong", func(c *gin.Context) {
		c.String(http.StatusOK, "🏓 Hello World, I'm the pong service")
	})

	r.Run(":8082")
}
