package handler

import (
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/gitlab"
	"net/http"
)

type environmentsResponse struct {
	Environments []*gitlab.Environment `json:"environments"`
}

// CreateEnvironmentHandler provides all environments
func CreateEnvironmentHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		environments := git.GetEnvironments()
		writeResponse(w, &environmentsResponse{Environments: environments})
	}
}
