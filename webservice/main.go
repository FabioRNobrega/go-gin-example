package main

import (
	"bytes"
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
)

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

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.tmpl")
	r.Static("/static/", "./static/")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Ping Pong",
		})
	})

	r.GET("/encrypt-decrypt", func(c *gin.Context) {
		c.HTML(http.StatusOK, "encrypt-decrypt/base", gin.H{
			"title": "Encrypt & Decrypt",
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

	r.Run(":8080")
}
