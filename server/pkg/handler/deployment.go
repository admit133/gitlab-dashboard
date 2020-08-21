package handler

import (
	"fmt"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/gitlab"
	"github.com/gorilla/mux"
	"net/http"
)

type deploymentsResponse struct {
	Deployments []*gitlab.Deployment `json:"deployments"`
}

// CreateListDeploymentHandler provides list of deployments for given projectID and environment
func CreateListDeploymentHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID, err := getRequiredIntFromVars(w, vars, "projectID")
		if err != nil {
			return
		}
		environment, err := getRequiredStringFromVars(w, vars, "environment")
		if err != nil {
			return
		}

		deployments, err := git.ListProjectDeployments(environment, projectID)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot get deployments: %v", err))
			return
		}

		writeResponse(w, &deploymentsResponse{Deployments: deployments})
		return
	}
}
