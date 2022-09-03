package queue

import (
	"context"
	"fmt"
	"log"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/mittz/roleplay-webapp-portal/utils"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const (
	QUEUE_REGION  = "us-central1"
	RANDOM_LENGTH = 6
)

type Queue struct {
	id            string
	queueFullPath string
}

func GetInstance() Queue {
	queueName := utils.GetEnvQueueName()
	return Queue{
		id:            queueName,
		queueFullPath: fmt.Sprintf("projects/%s/locations/%s/queues/%s", utils.GetEnvProjectID(), QUEUE_REGION, queueName),
	}
}

func (q Queue) EnqueueTask(userkey string) error {
	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()

	taskFullPath := fmt.Sprintf("%s/tasks/%s-%s", q.queueFullPath, userkey, utils.RandomString(RANDOM_LENGTH))
	req := &taskspb.CreateTaskRequest{
		Parent: q.queueFullPath,
		Task: &taskspb.Task{
			Name: taskFullPath,
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("https://us-central1-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/%s/jobs/assess-%s:run", utils.GetEnvProjectID(), strings.ToLower(userkey)),
					AuthorizationHeader: &taskspb.HttpRequest_OauthToken{
						OauthToken: &taskspb.OAuthToken{
							ServiceAccountEmail: "portal@role-play-web-app-host-project.iam.gserviceaccount.com",
						},
					},
				},
			},
		},
	}

	if _, err := c.CreateTask(ctx, req); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
