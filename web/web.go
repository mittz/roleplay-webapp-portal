package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/mittz/roleplay-webapp-portal/database"
	"github.com/mittz/roleplay-webapp-portal/image"
	"github.com/mittz/roleplay-webapp-portal/job"
	"github.com/mittz/roleplay-webapp-portal/queue"
	"github.com/mittz/roleplay-webapp-portal/request"
	"github.com/mittz/roleplay-webapp-portal/user"
	"github.com/mittz/roleplay-webapp-portal/utils"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	REQUEST_TIMEOUT_SECOND         = 30
	MESSAGE_INVALID_ENDPOINT       = "invalid endpoint"
	MESSAGE_INVALID_USERKEY        = "invalid userkey"
	MESSAGE_ALREADY_INQUEUE        = "already in the queue"
	MESSAGE_FAIL_HANDLEJOB         = "failed to create or replace job"
	MESSAGE_FAIL_ENQUEUE           = "failed to enqueue task"
	MESSAGE_INVALID_DATA_STRUCTURE = "invalid data structure"
	MESSAGE_FAILED_BULK_IMPORT     = "failed to bulk-import data"
	MESSAGE_INVALID_API_KEY        = "invalid api key"
)

func getRequestForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"ldaps":          job.GetLDAPsOfRunningExecutions(),
		"userkey":        c.DefaultQuery("userkey", ""),
		"datastudio_url": utils.GetEnvDataStudioURL(),
	})
}

func timeoutPostBenchmark(c *gin.Context) {
	c.HTML(http.StatusRequestTimeout, "timeout.html", nil)
}

func postBenchmark(c *gin.Context) {
	r := request.NewRequest(c.PostForm("userkey"), c.PostForm("endpoint"), c.PostForm("project_id"))
	if !r.IsValidEndpoint() {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"userkey":        r.Userkey,
			"datastudio_url": utils.GetEnvDataStudioURL(),
			"message":        MESSAGE_INVALID_ENDPOINT,
		})
		return
	}

	if !r.IsValidUserKey() {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"userkey":        r.Userkey,
			"datastudio_url": utils.GetEnvDataStudioURL(),
			"message":        MESSAGE_INVALID_USERKEY,
		})
		return
	}

	if job.IsRunning(r.Userkey) {
		c.HTML(http.StatusNotAcceptable, "error.html", gin.H{
			"userkey":        r.Userkey,
			"datastudio_url": utils.GetEnvDataStudioURL(),
			"message":        MESSAGE_ALREADY_INQUEUE,
		})
		return
	}

	if err := job.CreateOrReplace(r.Userkey, r.Endpoint, r.ProjectID); err != nil {
		c.HTML(http.StatusNotAcceptable, "error.html", gin.H{
			"userkey":        r.Userkey,
			"datastudio_url": utils.GetEnvDataStudioURL(),
			"message":        MESSAGE_FAIL_HANDLEJOB,
		})
		return
	}

	q := queue.GetInstance()
	if err := q.EnqueueTask(r.Userkey); err != nil {
		c.HTML(http.StatusNotAcceptable, "error.html", gin.H{
			"userkey":        r.Userkey,
			"datastudio_url": utils.GetEnvDataStudioURL(),
			"message":        MESSAGE_FAIL_ENQUEUE,
		})
		return
	}

	c.HTML(http.StatusAccepted, "benchmark.html", gin.H{
		"endpoint":       r.Endpoint,
		"userkey":        r.Userkey,
		"projectid":      r.ProjectID,
		"datastudio_url": utils.GetEnvDataStudioURL(),
	})
}

func isValidAdminAPIKey(apikey string) bool {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", utils.GetEnvProjectID(), "portal-admin-api-key")
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Println(err)
	}

	r := &secretmanagerpb.AccessSecretVersionRequest{Name: name}
	s, err := c.AccessSecretVersion(ctx, r)
	if err != nil {
		log.Println(err)
	}

	return apikey == string(s.Payload.Data)
}

func bulkImportUsers(c *gin.Context) {
	if !isValidAdminAPIKey(c.Request.Header.Get("Admin-API-Key")) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": MESSAGE_INVALID_API_KEY,
		})
		return
	}

	b := user.NewBulkUsers()
	if err := c.ShouldBind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": MESSAGE_INVALID_DATA_STRUCTURE,
		})
		return
	}

	if err := b.OverrideDatabase(); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": MESSAGE_FAILED_BULK_IMPORT,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Successfully the data was imported.",
	})
}

func bulkImportImageHashes(c *gin.Context) {
	b := image.NewBulkImageHashes()
	if err := c.ShouldBind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": MESSAGE_INVALID_DATA_STRUCTURE,
		})
		return
	}

	if err := b.OverrideDatabase(); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": MESSAGE_FAILED_BULK_IMPORT,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Successfully the data was imported.",
	})
}

func initJobHistory(c *gin.Context) {
	dbPool := database.GetDatabaseConnection()

	queryTableCreation := `
	DROP TABLE IF EXISTS job_histories;
	CREATE TABLE job_histories (
		id SERIAL NOT NULL,
		userkey character varying(40) NOT NULL,
		ldap character varying(20),
		score INTEGER,
		score_by_cost double precision,
		performance INTEGER,
		availability_rate INTEGER,
		message character varying(200),
		cost double precision,
		executed_at timestamp,
		PRIMARY KEY(id)
	);
	GRANT ALL ON job_histories TO PUBLIC;
	GRANT USAGE ON SEQUENCE job_histories_id_seq TO PUBLIC;
	`

	if _, err := dbPool.Exec(context.Background(), queryTableCreation); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": fmt.Errorf("Table recreation failed: %v\n", err),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Successfully the table was inited.",
	})
}

func initRanking(c *gin.Context) {
	dbPool := database.GetDatabaseConnection()

	queryTableCreation := `
	DROP TABLE IF EXISTS rankings;
	CREATE TABLE rankings (
		ldap character varying(20),
		score INTEGER,
		score_by_cost double precision,
		executed_at timestamp,
		PRIMARY KEY(ldap)
	);
	GRANT ALL ON rankings TO PUBLIC;
	`

	if _, err := dbPool.Exec(context.Background(), queryTableCreation); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": fmt.Errorf("Table recreation failed: %v\n", err).Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Successfully the table was inited.",
	})
}

func StartApp() {
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.StaticFile("/favicon.ico", "favicon.ico")

	r.GET("/", getRequestForm)
	r.POST("/benchmark", timeout.New(
		timeout.WithTimeout(time.Second*REQUEST_TIMEOUT_SECOND),
		timeout.WithHandler(postBenchmark),
		timeout.WithResponse(timeoutPostBenchmark),
	))

	r.POST("/admin/bulk/users", bulkImportUsers)
	r.POST("/admin/bulk/imagehashes", bulkImportImageHashes)
	r.POST("/admin/init/jobhistory", initJobHistory)
	r.POST("/admin/init/ranking", initRanking)

	r.Run(fmt.Sprintf(":%d", utils.GetEnvPortalPort()))
}
