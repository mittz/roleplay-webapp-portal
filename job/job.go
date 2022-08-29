package job

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mittz/roleplay-webapp-portal/utils"
)

const (
	IMAGE_URL               = "us-central1-docker.pkg.dev/role-play-web-app-host-project/roleplay-webapp/assess:latest"
	ASSESS_SERVICEA_ACCOUNT = "benchmark@role-play-web-app-host-project.iam.gserviceaccount.com"
	MAX_RETRIES             = 0
	ASSESS_CPU              = 4
	ASSESS_MEMORY_GB        = 2
)

type Job struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name string `json:"name"`
}

func jobExists(jobName string) bool {
	outDescribe, err := exec.Command(
		"gcloud",
		"beta",
		"run",
		"jobs",
		"list",
		"--region=us-central1",
		fmt.Sprintf("--project=%s", utils.GetEnvProjectID()),
		"--format=json",
	).Output()
	if err != nil {
		log.Printf("jobExists: %v", err)
		return false
	}

	var jobs []Job
	if err := json.Unmarshal(outDescribe, &jobs); err != nil {
		log.Println(err)
		return false
	}

	for _, job := range jobs {
		if job.Metadata.Name == jobName {
			return true
		}
	}

	return false
}

func createJob(jobName string, userkey string, endpoint string, projectID string) error {
	if _, err := exec.Command(
		"gcloud",
		"beta",
		"run",
		"jobs",
		"create",
		jobName,
		"--region=us-central1",
		"--image",
		IMAGE_URL,
		"--max-retries",
		strconv.Itoa(MAX_RETRIES),
		"--set-env-vars",
		fmt.Sprintf("USER_KEY=%s", userkey),
		"--set-env-vars",
		fmt.Sprintf("ENDPOINT=%s", endpoint),
		"--set-env-vars",
		fmt.Sprintf("PROJECT_ID=%s", projectID),
		"--set-env-vars",
		fmt.Sprintf("INSTANCE_CONNECTION_NAME=%s", utils.GetEnvInstanceConnectionName()),
		"--set-env-vars",
		fmt.Sprintf("DATABASE_NAME=%s", utils.GetEnvDatabaseName()),
		"--set-env-vars",
		fmt.Sprintf("DATABASE_USER=%s", "benchmark@role-play-web-app-host-project.iam"),
		"--service-account",
		ASSESS_SERVICEA_ACCOUNT,
		"--cpu",
		fmt.Sprintf("%d", ASSESS_CPU),
		"--memory",
		fmt.Sprintf("%dGi", ASSESS_MEMORY_GB),
		"--format=json",
	).Output(); err != nil {
		log.Printf("createJob: %v", err)
		return err
	}

	return nil
}

func updateJob(jobName string, userkey string, endpoint string, projectID string) error {
	if _, err := exec.Command(
		"gcloud",
		"beta",
		"run",
		"jobs",
		"update",
		jobName,
		"--region=us-central1",
		"--set-env-vars",
		fmt.Sprintf("USER_KEY=%s", userkey),
		"--set-env-vars",
		fmt.Sprintf("ENDPOINT=%s", endpoint),
		"--set-env-vars",
		fmt.Sprintf("PROJECT_ID=%s", projectID),
		"--set-env-vars",
		fmt.Sprintf("INSTANCE_CONNECTION_NAME=%s", utils.GetEnvInstanceConnectionName()),
		"--set-env-vars",
		fmt.Sprintf("DATABASE_NAME=%s", utils.GetEnvDatabaseName()),
		"--set-env-vars",
		fmt.Sprintf("DATABASE_USER=%s", "benchmark@role-play-web-app-host-project.iam"),
		"--service-account",
		ASSESS_SERVICEA_ACCOUNT,
		"--cpu",
		fmt.Sprintf("%d", ASSESS_CPU),
		"--memory",
		fmt.Sprintf("%dGi", ASSESS_MEMORY_GB),
		"--format=json",
	).Output(); err != nil {
		log.Printf("updateJob: %v", err)
		return err
	}

	return nil
}

func CreateOrReplace(userkey string, endpoint string, projectID string) error {
	jobName := fmt.Sprintf("assess-%s", strings.ToLower(userkey))

	if !jobExists(jobName) {
		if err := createJob(jobName, userkey, endpoint, projectID); err != nil {
			return err
		}
	} else {
		if err := updateJob(jobName, userkey, endpoint, projectID); err != nil {
			return err
		}
	}

	return nil
}
