package handler

import (
	"fmt"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/config"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/gitlab"
	"net/http"
)

type configResponse struct {
	GitLabBaseURL    string              `json:"gitLabBaseURL"`
	GitLabAppID      string              `json:"gitLabAppId"`
	UserLinkTemplate string              `json:"userLinkTemplate"`
	OAuthEnabled     bool                `json:"oAuthEnabled"`
	User             *gitlab.ProjectUser `json:"user"`
}

// CreateConfigHandler provides basing configuration for GUI
func CreateConfigHandler(userService *gitlab.UserService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		user, err := userService.GetUserFromRequest(request)
		if err != nil {
			badRequest(w, fmt.Sprintf("cannot get user from gitlab: %v", err))
			return
		}

		response := configResponse{
			GitLabBaseURL:    cfg.GitLabBaseURL,
			UserLinkTemplate: cfg.UserLinkTemplate,
			GitLabAppID:      cfg.GitLabAppID,
			OAuthEnabled:     cfg.OAuthEnabled,
			User:             user,
		}

		writeResponse(w, &response)
	}
}
