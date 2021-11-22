package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	gitlab2 "github.com/xanzy/go-gitlab"
	"gitlab-environment-dashboard/server/pkg/gitlab"
	"net/http"
)

type jobResponse struct {
	Job *gitlab2.Job `json:"job"`
}

type playJobRequestBody struct {
	Ref string `json:"ref"`
}

type playJobsRequestBody struct {
	Query string `json:"query"`
}

type jobsListResponse struct {
	Jobs map[string]map[int]*gitlab2.Job `json:"jobs"`
}

// CreatePlayJobHandler plays or retries a job for given projectId and environment
func CreatePlayJobHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		environment, err := getRequiredStringFromVars(w, vars, "environment")
		if err != nil {
			return
		}
		projectID, err := getRequiredIntFromVars(w, vars, "projectID")
		if err != nil {
			return
		}

		requestBody := playJobRequestBody{}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot parse request body: %v", err))
			return
		}
		if requestBody.Ref == "" {
			badRequest(w, "ref is empty")
			return
		}
		deployment, err := git.PlayOrRetryJob(projectID, environment, requestBody.Ref)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot create job: %v", err))
			return
		}

		writeResponse(w, &jobResponse{Job: deployment})
		return
	}
}

// CreateGetJobHandler provides a job for given environment and given projectID
func CreateGetJobHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		environment, err := getRequiredStringFromVars(w, vars, "environment")
		if err != nil {
			return
		}
		projectID, err := getRequiredIntFromVars(w, vars, "projectID")
		if err != nil {
			return
		}

		job, _ := git.GetJob(environment, projectID)

		writeResponse(w, &jobResponse{Job: job})
		return
	}
}

// CreatePlayJobsByQueryHandler plays or retries a job for given query
// Query is substring for branch name
func CreatePlayJobsByQueryHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		environment, err := getRequiredStringFromVars(w, vars, "environment")
		if err != nil {
			return
		}
		requestBody := playJobsRequestBody{}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot parse request body: %v", err))
			return
		}
		if requestBody.Query == "" {
			badRequest(w, "query is empty")
			return
		}
		if len(requestBody.Query) < 3 {
			badRequest(w, "query is too small (min 3 symbols)")
			return
		}

		err = git.PlayOrRetryJobsWithQuery(environment, requestBody.Query)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot start jobs: %v", err))
			return
		}

		writeResponse(w, &jobsListResponse{Jobs: git.GetJobs()})
		return
	}
}

// CreateListJobsHandler provides list of all jobs
func CreateListJobsHandler(git *gitlab.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := git.GetJobs()

		writeResponse(w, &jobsListResponse{Jobs: jobs})
		return
	}
}
