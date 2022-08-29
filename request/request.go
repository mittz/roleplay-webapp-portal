package request

import (
	"strings"

	"github.com/mittz/roleplay-webapp-portal/user"
)

type Request struct {
	Userkey   string
	Endpoint  string
	ProjectID string
}

func NewRequest(userkey string, endpoint string, projectID string) Request {
	return Request{
		Userkey:   userkey,
		Endpoint:  endpoint,
		ProjectID: projectID,
	}
}

func (r Request) IsValidEndpoint() bool {
	return (strings.HasPrefix(r.Endpoint, "http://") ||
		strings.HasPrefix(r.Endpoint, "https://")) &&
		(!strings.Contains(r.Endpoint, "localhost") &&
			!strings.Contains(r.Endpoint, "127.0.0.1"))
}

func (r Request) IsValidUserKey() bool {
	user := user.GetUser(r.Userkey)
	return user.Userkey == r.Userkey
}
