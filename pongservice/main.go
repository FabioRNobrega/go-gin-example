package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func loadPrivateKey() *rsa.PrivateKey {
	privKeyPEM := os.Getenv("RSA_PRIVATE_KEY")
	if privKeyPEM == "" {
		panic("RSA_PRIVATE_KEY not set")
	}
	privKeyPEM = strings.ReplaceAll(privKeyPEM, `\n`, "\n")

	block, _ := pem.Decode([]byte(privKeyPEM))
	if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
		panic("failed to decode PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse private key: " + err.Error())
	}
	return privKey.(*rsa.PrivateKey)
}

func main() {
	privKey := loadPrivateKey()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/pong", func(c *gin.Context) {
		c.String(http.StatusOK, "🏓 Hello World, I'm the pong service")
	})

	r.POST("/decrypt", func(c *gin.Context) {
		var req struct {
			CipherText string `json:"cipher_text"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		cleanCipher := strings.TrimSpace(req.CipherText)

		cipherBytes, err := base64.StdEncoding.DecodeString(cleanCipher)
		if err != nil {
			log.Printf("Base64 decode failed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64"})
			return
		}

		plainBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, cipherBytes)
		if err != nil {
			log.Printf("RSA decrypt failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"decrypted_text": string(plainBytes),
		})
	})

	r.Run(":8082")
}
