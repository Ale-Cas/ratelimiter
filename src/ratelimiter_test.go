package src

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	router *gin.Engine
	req    *http.Request
)

const bucketCapacity = uint64(1)

// setup initializes the router and the request.
func setup() {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.Use(TokenBucketLimiter(bucketCapacity))
	router.GET("/", func(c *gin.Context) {})

	req, _ = http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:1234"

	buckets = []TokenBucket{}
}

// testRequest sends a request to the server and checks the status code.
func testRequest(t *testing.T, expectedStatus int) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code)
}

// TestTokenBucketLimiter_NewClient tests the case where a new client makes a request.
func TestTokenBucketLimiter_NewClient(t *testing.T) {
	setup()
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 1, len(buckets))
	assert.Equal(t, "192.0.2.1", buckets[0].ipAddr)
	assert.Equal(t, bucketCapacity, buckets[0].tokens)
}

// TestTokenBucketLimiter_ExistingClient tests the case where an existing client makes a request.
func TestTokenBucketLimiter_ExistingClient(t *testing.T) {
	setup()
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	router.ServeHTTP(w, req)

	assert.Equal(t, 1, len(buckets))
	assert.Equal(t, "192.0.2.1", buckets[0].ipAddr)
	assert.Equal(t, uint64(0), buckets[0].tokens)
}

// TestTokenBucketLimiter_TooManyRequests tests the case where a client makes too many requests.
func TestTokenBucketLimiter_TooManyRequests(t *testing.T) {
	setup()
	testRequest(t, http.StatusOK)
	testRequest(t, http.StatusOK)
	testRequest(t, http.StatusTooManyRequests)
	testRequest(t, http.StatusTooManyRequests)
	testRequest(t, http.StatusTooManyRequests)

	// refill bucket
	time.Sleep(2 * time.Second)

	testRequest(t, http.StatusOK)
}
