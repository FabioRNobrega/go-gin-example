package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get some books
	r.GET("/books", func(c *gin.Context) {
		books := []map[string]interface{}{
			{"id": 1, "title": "The Go Programming Language", "author": "Alan A. A. Donovan"},
			{"id": 2, "title": "Introducing Go", "author": "Caleb Doxsey"},
			{"id": 3, "title": "Go in Action", "author": "William Kennedy"},
		}
		c.JSON(http.StatusOK, gin.H{"books": books})
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic Zm9vOmJhcg==' \
		-H 'content-type: application/json' \
		-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	/*
		Gin templates
		To load templates, you can use LoadHTMLFiles or LoadHTMLGlob.
		LoadHTMLFiles loads the specified files, while LoadHTMLGlob loads all files matching the specified pattern.
		For example, to load all templates in the "templates" directory, you can use LoadHTMLGlob like this:
		router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	*/
	r.LoadHTMLGlob("templates/*")
	r.Static("/static/", "./static/")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Home",
		})
	})

	return r
}

func main() {
	// Entry point of the application: sets up the router and starts the HTTP server
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
