package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the authenticity of the request based on X-UserId and X-Digest headers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-UserId")
		digest := c.GetHeader("X-Digest")

		if userID == "" || digest == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Read the body of the request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		// Restore the request body for future handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Compute the HMAC of the request body
		computedDigest := computeHMAC(body)
		if digest != computedDigest {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// computeHMAC computes the HMAC-SHA1 hash of a message
func computeHMAC(message []byte) string {
	key := []byte("secret")
	h := hmac.New(sha1.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}
