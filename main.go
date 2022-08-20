package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

const (
	REQUEST_TIMEOUT_SECOND   = 10
	MESSAGE_INVALID_ENDPOINT = "invalid endpoint"
	MESSAGE_INVALID_USERKEY  = "invalid userkey"
	MESSAGE_ALREADY_INQUEUE  = "already in the queue"
)

func isValidEndpoint(endpoint string) bool {
	return (strings.HasPrefix(endpoint, "http://") ||
		strings.HasPrefix(endpoint, "https://")) &&
		(!strings.Contains(endpoint, "localhost") &&
			!strings.Contains(endpoint, "127.0.0.1"))
}

func isValidUserKey(userkey string) bool {
	return true
}

func isUserQueued(userkey string) bool {
	return true
}

type Request struct{}

func getRequests() []Request {
	return []Request{}
}

func getRequestForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"requests":       getRequests(),
		"datastudio_url": getEnvDataStudioURL(),
	})
}

func timeoutPostBenchmark(c *gin.Context) {
	c.HTML(http.StatusRequestTimeout, "timeout.html", nil)
}

func postBenchmark(c *gin.Context) {
	userkey := c.PostForm("userkey")
	endpoint := c.PostForm("endpoint")
	projectID := c.PostForm("project_id")

	if !isValidEndpoint(endpoint) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": MESSAGE_INVALID_ENDPOINT,
		})
		return
	}

	if !isValidUserKey(userkey) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": MESSAGE_INVALID_USERKEY,
		})
		return
	}

	if isUserQueued(userkey) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": MESSAGE_ALREADY_INQUEUE,
		})
		return
	}

	c.HTML(http.StatusAccepted, "benchmark.html", gin.H{
		"endpoint":  endpoint,
		"userkey":   userkey,
		"projectid": projectID,
	})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/favicon.ico", "favicon.ico")

	r.GET("/", getRequestForm)
	r.POST("/benchmark", timeout.New(
		timeout.WithTimeout(time.Second*REQUEST_TIMEOUT_SECOND),
		timeout.WithHandler(postBenchmark),
		timeout.WithResponse(timeoutPostBenchmark),
	))

	r.Run(fmt.Sprintf(":%d", getEnvPortalPort()))
}
