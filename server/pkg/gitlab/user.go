package gitlab

import (
	"errors"
	log "github.com/sirupsen/logrus"
	wrappedGitlab "github.com/xanzy/go-gitlab"
	"net/http"
)

// UserService using for authenticate current user
// It keeps token and do not ask GitLab twice with same token
// Probably we need to track all tokens and clean it
// but for first implementation it's ok
type UserService struct {
	tokens        map[string]*ProjectUser
	gitlabService *Service
	gitlabBaseURL string
}

func NewUserService(
	service *Service,
	gitlabBaseURL string,
) *UserService {
	return &UserService{
		tokens:        map[string]*ProjectUser{},
		gitlabBaseURL: gitlabBaseURL,
		gitlabService: service,
	}
}

func (s *UserService) GetUserFromRequest(request *http.Request) (user *ProjectUser, err error) {
	tokenCookie, err := request.Cookie("token")
	if err != nil && err != http.ErrNoCookie {
		return nil, errors.New("cannot read cookies")
	}
	if tokenCookie == nil {
		return nil, nil
	}

	// Check stored token firstly
	user = s.getStoredUser(tokenCookie.Value)
	if user != nil {
		return user, nil
	}

	client, err := wrappedGitlab.NewOAuthClient(tokenCookie.Value, wrappedGitlab.WithBaseURL(s.gitlabBaseURL))
	if err != nil {
		return nil, errors.New("cannot create gitlab client")
	}

	remoteUser, _, err := client.Users.CurrentUser()
	if err != nil {
		log.Error(err)
		// Most likely is unauthorized
		return nil, nil
	}

	if remoteUser != nil {
		user = &ProjectUser{
			Username:  remoteUser.Username,
			Name:      remoteUser.Name,
			AvatarURL: remoteUser.AvatarURL,
		}
		s.storeUserToken(tokenCookie.Value, user)
	}

	return
}

func (s *UserService) storeUserToken(token string, user *ProjectUser) {
	s.tokens[token] = user
}

func (s *UserService) getStoredUser(token string) (user *ProjectUser) {
	user, ok := s.tokens[token]
	if !ok {
		return nil
	}

	return user
}
