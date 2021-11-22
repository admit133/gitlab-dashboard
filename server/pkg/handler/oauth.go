package handler

import (
	"encoding/json"
	"fmt"
	"gitlab-environment-dashboard/server/pkg/gitlab"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	TokenCookieName = "token"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func CreateOauthHandler(gitLabBaseUrl string, gitLabAppID string, gitLabSecret string, sslEnabled bool, cookieSecured bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		// Check errors first
		// If any error occurred we redirect to rhe main page
		errorMessage := query.Get("error")
		if errorMessage != "" {
			log.Errorf("Error authentication: %s", query.Encode())
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			return
		}

		code := query.Get("code")

		scheme := "http"
		if sslEnabled {
			scheme = "https"
		}

		redirectUrl := scheme + "://" + request.Host + request.URL.Path
		log.Info(redirectUrl)

		requestUrl := fmt.Sprintf(
			"%s/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
			gitLabBaseUrl,
			gitLabAppID,
			gitLabSecret,
			code,
			redirectUrl,
		)
		resp, err := http.Post(requestUrl, "application/json", strings.NewReader(""))
		if err != nil {
			log.Error(err)
			badRequest(writer, "cannot get token")
			return
		}
		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			log.Printf("%v,%v", string(body), err)
			badRequest(writer, "authorization failed, url: "+requestUrl+",code: "+string(resp.StatusCode)+",body: "+string(body))
			return
		}

		tokenResponse := tokenResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		err = json.NewDecoder(strings.NewReader(string(body))).Decode(&tokenResponse)
		if err != nil {
			log.Error(err)
			badRequest(writer, "cannot decode token")
		}

		http.SetCookie(writer, &http.Cookie{
			Name:     TokenCookieName,
			Value:    tokenResponse.AccessToken,
			Domain:   "*",
			Path:     "/",
			Expires:  time.Now().Add(7600 * time.Second),
			Secure:   cookieSecured,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
	}
}

func CreateLogoutHandler(cookieSecured bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.SetCookie(
			writer,
			&http.Cookie{
				Name:     TokenCookieName,
				Value:    "",
				Expires:  time.Now(),
				Domain:   "*",
				Path:     "/",
				Secure:   cookieSecured,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
	}
}

func CreateAuthMiddleware(userService *gitlab.UserService, handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user, err := userService.GetUserFromRequest(request)
		if err != nil {
			badRequest(writer, fmt.Sprintf("cannot authorize user: %v", err))
			return
		}
		if user == nil {
			unauthorizedRequest(writer, "unauthorized")
			return
		}

		handler.ServeHTTP(writer, request)
	}
}
