package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionTheme struct {
	SessionID string `json:"session_id"`
	DarkMode  bool   `json:"dark_mode"`
}

var ctx = context.Background()

func newRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("error connecting to Redis: %v", err))
	}

	fmt.Println("connected to Redis in", host, port)
	return rdb
}

func loadPublicKey() *rsa.PublicKey {
	pubKeyPEM := os.Getenv("RSA_PUBLIC_KEY")
	if pubKeyPEM == "" {
		panic("RSA_PUBLIC_KEY not set")
	}
	pubKeyPEM = strings.ReplaceAll(pubKeyPEM, `\n`, "\n")

	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		panic("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("failed to parse public key: %v", err))
	}
	return pub.(*rsa.PublicKey)
}

func main() {
	pubKey := loadPublicKey()
	redisDb := newRedisClient()

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.tmpl")
	r.Static("/static/", "./static/")

	r.GET("/", func(c *gin.Context) {
		sessionID := uuid.New().String()
		c.HTML(http.StatusOK, "base", gin.H{
			"title":      "Ping Pong",
			"session_id": sessionID,
			"dark_mode":  true,
		})
	})

	r.GET("/encrypt-decrypt", func(c *gin.Context) {
		sessionID := uuid.New().String()
		c.HTML(http.StatusOK, "encrypt-decrypt/base", gin.H{
			"title":      "Encrypt & Decrypt",
			"session_id": sessionID,
			"dark_mode":  true,
		})
	})

	// Connection with different services

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

	// Using keys for decrypt and encrypt

	r.POST("/encrypt", func(c *gin.Context) {
		text := c.PostForm("text")
		if text == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
			return
		}

		cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(text))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encrypt failed"})
			return
		}

		c.String(http.StatusOK, base64.StdEncoding.EncodeToString(cipher))
	})

	r.POST("/decrypt", func(c *gin.Context) {
		cipherText := c.PostForm("text")
		if cipherText == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text (cipher) is required"})
			return
		}

		body, _ := json.Marshal(map[string]string{"cipher_text": cipherText})
		resp, err := http.Post("http://pongservice:8082/decrypt", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call pongservice"})
			return
		}
		defer resp.Body.Close()

		var pongResp struct {
			DecryptedText string `json:"decrypted_text"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&pongResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse pongservice response"})
			return
		}

		c.String(http.StatusOK, pongResp.DecryptedText)
	})

	// Redis
	r.PUT("/theme", func(c *gin.Context) {
		ctx := c.Request.Context()

		sessionID := c.PostForm("session_id")
		if sessionID == "" {
			sessionID = uuid.New().String()
		}

		darkMode := c.PostForm("darkMode") == "on"

		var session SessionTheme
		val, err := redisDb.Get(ctx, "session:"+sessionID).Result()
		if err == redis.Nil {
			session = SessionTheme{SessionID: sessionID, DarkMode: darkMode}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		} else {
			if err := json.Unmarshal([]byte(val), &session); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
				return
			}
			session.DarkMode = darkMode
		}

		data, _ := json.Marshal(session)
		if err := redisDb.Set(ctx, "session:"+sessionID, data, 0).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save theme"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session_id": session.SessionID,
			"dark_mode":  session.DarkMode,
		})
	})

	r.Run(":8080")
}
