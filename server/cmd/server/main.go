package main

import (
	"context"
	"fmt"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/config"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/gitlab"
	"github.com/fr05t1k/gitlab-environment-dashboard/server/pkg/handler"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.CreateConfig()
	gitLabService, err := gitlab.NewClient(
		cfg.GitLabToken,
		cfg.GitLabBaseURL,
		cfg.ProtectedEnvironments,
		cfg.GitLabProjectIDs,
	)
	catchFatalError(err, "cannot create gitlab client: %v", err)
	userService := gitlab.NewUserService(
		gitLabService,
		cfg.GitLabBaseURL,
	)

	scheduleUpdateEnvironments(gitLabService, cfg.UpdateDuration, cfg.GitLabProjectIDs)
	scheduleUpdateBranches(gitLabService, cfg.UpdateDuration, cfg.GitLabProjectIDs)
	err = gitLabService.UpdateEnvironments(cfg.GitLabProjectIDs)
	catchFatalError(err, "cannot update environments: %v", err)
	err = gitLabService.UpdateBranches(cfg.GitLabProjectIDs)
	catchFatalError(err, "cannot update branches: %v", err)

	r := mux.NewRouter()

	addRoutes(r, gitLabService, cfg, userService)
	srv := &http.Server{
		Addr: cfg.ListenAddr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Printf(fmt.Sprintf("listen on: %s", cfg.ListenAddr))

	go func() {
		err = srv.ListenAndServe()
		catchFatalError(err, "cannot start the server: %v", err)
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	catchFatalError(err, "cannot shutdown server: %v", err)
	log.Info("graceful shutting down")
	os.Exit(0)
}

func catchFatalError(err error, format string, args ...interface{}) {
	if err != nil {
		log.Fatalf(format, args...)
	}
}

func addRoutes(r *mux.Router, gitLabService *gitlab.Service, cfg config.Config, userService *gitlab.UserService) {
	wrapWithMiddleware := CreateMiddlewareWrapper(userService)

	// Warning!!!
	// Private area
	// Check token first
	r.Methods("GET").
		Path("/environments").
		Handler(wrapWithMiddleware(
			handler.CreateEnvironmentHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("GET").
		Path("/environments/{environment}/projects/{projectID:[0-9]+}/repository/branches").
		Handler(wrapWithMiddleware(
			handler.CreateBranchHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("POST").
		Path("/environments/{environment}/jobs").
		Handler(wrapWithMiddleware(
			handler.CreatePlayJobsByQueryHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("POST").
		Path("/environments/{environment}/projects/{projectID:[0-9]+}/jobs").
		Handler(wrapWithMiddleware(
			handler.CreatePlayJobHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("GET").
		Path("/environments/{environment}/projects/{projectID:[0-9]+}/jobs").
		Handler(wrapWithMiddleware(
			handler.CreateGetJobHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("GET").
		Path("/environments/{environment}/projects/{projectID:[0-9]+}/deployments").
		Handler(wrapWithMiddleware(
			handler.CreateListDeploymentHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("GET").
		Path("/jobs").
		Handler(wrapWithMiddleware(
			handler.CreateListJobsHandler(gitLabService),
			cfg.OAuthEnabled,
		))

	r.Methods("GET").
		Path("/oauth/logout").
		Handler(wrapWithMiddleware(
			handler.CreateLogoutHandler(cfg.CookieSecured),
			cfg.OAuthEnabled,
		))

	// Warning!!!
	// Public area bellow
	// It's not restricted area
	// We don't check token here
	r.Methods("GET").
		Path("/config").
		Handler(wrapWithMiddleware(
			handler.CreateConfigHandler(userService, cfg),
			false,
		))

	r.Methods("GET").
		Path("/oauth/code").
		Handler(wrapWithMiddleware(
			handler.CreateOauthHandler(cfg.GitLabBaseURL, cfg.GitLabAppID, cfg.GitLabAppSecret, cfg.SslEnabled, cfg.CookieSecured),
			false,
		))

	r.Methods("GET").
		Path("/health").
		Handler(wrapWithMiddleware(
			handler.CreateHealthHandler(),
			false,
		))

	r.Methods("GET", "POST").
		Path("/error").
		Handler(wrapWithMiddleware(
			handler.CreateErrorHandler(),
			false,
		))

	r.Path("/metrics").
		Handler(promhttp.Handler())

	r.PathPrefix("/").
		Handler(wrapWithMiddleware(
			handler.CreateSPAHandler(cfg.PublicDir, "index.html"),
			false,
		))
}

// MiddlewareWrapper wrap handler with basic middlewares
// `restrictedArea` forbids unauthorized actions
type MiddlewareWrapper func(handlerFunc http.HandlerFunc, restrictedArea bool) (wrappedHandler http.Handler)

func CreateMiddlewareWrapper(userService *gitlab.UserService) MiddlewareWrapper {
	return func(handlerFunc http.HandlerFunc, restrictedArea bool) (wrappedHandler http.Handler) {
		wrappedHandler = handlers.CombinedLoggingHandler(os.Stdout, handlerFunc)
		if restrictedArea {
			wrappedHandler = handler.CreateAuthMiddleware(userService, wrappedHandler)
		}
		return
	}
}

func scheduleUpdateEnvironments(service *gitlab.Service, duration time.Duration, projectIDs []int) {
	log.Infof("environments will be updated every: %v\n", duration)
	go func() {
		for {
			log.Info("environments update has been started")
			start := time.Now()
			err := service.UpdateEnvironments(projectIDs)
			status := "environments update has been completed."
			if err != nil {
				status = "error occurred"
				log.Error(err)
			}
			log.Infof("%s. elapsed: %v\n", status, time.Since(start))
			<-time.After(duration)
		}
	}()
}

func scheduleUpdateBranches(service *gitlab.Service, duration time.Duration, projectIDs []int) {
	log.Infof("branches will be updated every: %v\n", duration)
	go func() {
		for {
			log.Info("branches update has been started")
			start := time.Now()
			err := service.UpdateBranches(projectIDs)
			status := "branches update has been completed."
			if err != nil {
				status = "error occurred"
				log.Error(err)
			}
			log.Infof("%s. elapsed: %v\n", status, time.Since(start))
			<-time.After(duration)
		}
	}()
}
