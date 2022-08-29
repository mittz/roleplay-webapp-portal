package queue

import (
	"context"
	"fmt"
	"log"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/mittz/roleplay-webapp-portal/utils"
	"google.golang.org/api/iterator"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const (
	QUEUE_REGION  = "us-central1"
	RANDOM_LENGTH = 6
	QUEUE_NAME    = "queue"
)

type Queue struct {
	id            string
	queueFullPath string
}

type Task struct {
	Name string
}

func GetInstance() Queue {
	return Queue{
		id:            QUEUE_NAME,
		queueFullPath: fmt.Sprintf("projects/%s/locations/%s/queues/%s", utils.GetEnvProjectID(), QUEUE_REGION, QUEUE_NAME),
	}
}

func (q Queue) GetTasks() []Task {
	var tasks []Task

	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return []Task{}
	}
	defer c.Close()

	req := &taskspb.ListTasksRequest{
		Parent: q.queueFullPath,
	}
	it := c.ListTasks(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			return []Task{}
		}

		tasks = append(tasks, Task{Name: resp.Name})
	}

	return tasks
}

func (q Queue) TaskExists(userkey string) bool {
	tasks := q.GetTasks()

	for _, task := range tasks {
		if strings.Contains(task.Name, userkey) {
			return true
		}
	}

	return false
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
