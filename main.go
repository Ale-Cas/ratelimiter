package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func limited(c *gin.Context) {
	// sleep for 1 second
	fmt.Println("Sleeping for 1 second")
	time.Sleep(time.Second)
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "limited",
	})
}

func unlimited(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "unlimited",
	})
}

func main() {
	port := flag.String("port", "8080", "Port for the web server")
	flag.Parse()
	r := gin.Default()
	r.GET("/limited", limited)
	r.GET("/unlimited", unlimited)
	r.Run("localhost:" + *port)
}
