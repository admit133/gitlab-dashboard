package handler

import (
	"fmt"
	"gitlab-environment-dashboard/server/pkg/gitlab"
	"github.com/gorilla/mux"
	wrappedGitlab "github.com/xanzy/go-gitlab"
	"net/http"
)

type branchesResponse struct {
	Branches []*wrappedGitlab.Branch `json:"branches"`
}

// CreateEnvironmentHandler provides all environments
func CreateBranchHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID, err := getRequiredIntFromVars(w, vars, "projectID")
		if err != nil {
			return
		}

		branches, err := git.GetBranches(projectID)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot get branches: %v", err))
			return
		}

		writeResponse(w, &branchesResponse{Branches: branches})
		return
	}
}
