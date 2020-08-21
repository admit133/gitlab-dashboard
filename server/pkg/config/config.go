// Package config provides all configuration for the application
package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config provides general application configuration
type Config struct {
	PublicDir             string
	GitLabBaseURL         string
	GitLabToken           string
	GitLabProjectIDs      []int
	UserLinkTemplate      string
	UpdateDuration        time.Duration
	ListenAddr            string
	ProtectedEnvironments []string
	GitLabAppID           string
	GitLabAppSecret       string
	CookieSecured         bool
	SslEnabled            bool
	OAuthEnabled          bool
}

// CreateConfig creates the application configuration
func CreateConfig() Config {
	var err error
	// We don't require .env
	_ = godotenv.Load()

	config := Config{}

	config.GitLabBaseURL = os.Getenv("GITLAB_BASE_URL")
	config.GitLabAppID = os.Getenv("GITLAB_APP_ID")
	config.GitLabAppSecret = os.Getenv("GITLAB_APP_SECRET")
	config.GitLabToken = os.Getenv("GITLAB_TOKEN")
	config.ListenAddr = os.Getenv("LISTEN_ADDRESS")
	config.UserLinkTemplate = os.Getenv("USER_LINK_TEMPLATE")
	config.PublicDir = os.Getenv("PUBLIC_DIR")
	config.CookieSecured = os.Getenv("COOKIE_SECURED") == "1"
	config.SslEnabled = os.Getenv("SSL_ENABLED") == "1"
	config.OAuthEnabled = os.Getenv("OAUTH_ENABLED") == "1"

	if os.Getenv("GITLAB_PROJECT_IDS") == "" {
		log.Fatalln("GITLAB_PROJECT_IDS should have at least one ID")
	}
	projectIDsStrings := strings.Split(os.Getenv("GITLAB_PROJECT_IDS"), ",")
	for _, idString := range projectIDsStrings {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Fatalf("GITLAB_PROJECT_IDS should have integers. %s given", idString)
		}
		config.GitLabProjectIDs = append(config.GitLabProjectIDs, int(id))
	}

	// Set Default duration
	config.UpdateDuration, err = time.ParseDuration("30s")
	if err != nil {
		log.Fatalln(err)
	}

	// Check if we have custom duration
	updateDuration := os.Getenv("ENVIRONMENT_UPDATE_DURATION")
	if updateDuration != "" {
		updateEach, err := time.ParseDuration(os.Getenv("ENVIRONMENT_UPDATE_DURATION"))
		if err != nil {
			log.Fatalf("Wrond duration format: %v\n", err)
		}

		config.UpdateDuration = updateEach
	}

	// Set default public DIR
	if config.PublicDir == "/" {
		config.PublicDir = "/public"
	}

	config.ProtectedEnvironments = strings.Split(os.Getenv("PROTECTED_ENVIRONMENTS"), ",")

	return config
}
