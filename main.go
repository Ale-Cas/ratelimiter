package main

import (
	"flag"
	"net/http"

	"github.com/Ale-Cas/ratelimiter/src"
	"github.com/gin-gonic/gin"
)

// limited is a simple request handler
// that returns a JSON response with a message.
func limited(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "limited",
	})
}

func main() {
	port := flag.String("port", "8080", "Port for the web server.")
	reqPerSec := flag.Uint64("reqPerSec", 10, "Number of allowed requests per second.")
	flag.Parse()
	r := gin.Default()
	r.Use(src.TokenBucketLimiter(*reqPerSec))
	r.GET("/limited", limited)
	r.Run(":" + *port)
}
