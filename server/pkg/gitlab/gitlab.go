/*
Package gitlab is an application abstraction upon github.com/xanzy/go-gitlab
*/
package gitlab

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	wrappedGitLab "github.com/xanzy/go-gitlab"
	"gitlab-environment-dashboard/server/pkg/utils"
	"sort"
	"sync"
	"time"
)

var (
	DeniedForProtectedEnvironment = errors.New("cannot perform the action for an protected environment")
	JobNotFound                   = errors.New("job not found")
	JobIsNotReady                 = errors.New("job is not ready")
)

// Service operates with gitlab API
type Service struct {
	git             *wrappedGitLab.Client
	environments    map[string]*Environment
	environmentsMtx sync.RWMutex

	// branches by project ID
	branches    map[int][]*wrappedGitLab.Branch
	branchesMtx sync.RWMutex

	// We store jobs which was run from the dashboard
	jobs                  map[string]map[int]*wrappedGitLab.Job
	jobsMtx               sync.RWMutex
	protectedEnvironments []string
	projectIDs            []int
	// Sometimes could have scheduled pipeline which doesn't have environments
	// We we try to run a job it finds first pipeline with the expected environment
	jobRecursiveSearchLimit int
}

// Environment represents a wrapper for wrappedGitLab.Environment
type Environment struct {
	Name     string     `json:"name"`
	Projects []*Project `json:"projects"`
}

// Project represents a wrapper for wrappedGitLab.Project
type Project struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	AvatarURL         string      `json:"avatarURL"`
	WebURL            string      `json:"webURL"`
	NameWithNamespace string      `json:"nameWithNamespace"`
	LastDeployment    *Deployment `json:"lastDeployment"`
}

// Deployment represents a wrapper for wrappedGitLab.Deployment
type Deployment struct {
	ID         int          `json:"id"`
	Ref        string       `json:"ref"`
	User       *ProjectUser `json:"user"`
	UpdatedAt  *time.Time   `json:"updatedAt"`
	Deployable *Deployable  `json:"deployable"`
}

// Deployable represents a status of a deploy launch
type Deployable struct {
	Pipeline Pipeline `json:"pipeline"`
}

// Pipeline represents a pipeline :)
type Pipeline struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	User   *User  `json:"user"`
}

type User struct {
	Username  string
	Name      string
	AvatarURL string
}

// ProjectUser represents a wrapper for wrappedGitLab.ProjectUser
type ProjectUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

// Sorting branches by committed_date desc
type ByCommitDateDesc []*wrappedGitLab.Branch

func (s ByCommitDateDesc) Len() int {
	return len(s)
}

func (s ByCommitDateDesc) Less(i, j int) bool {
	return s[i].Commit.CommittedDate.After(*s[j].Commit.CommittedDate)
}

func (s ByCommitDateDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// GetBranches returns all branches for given projectID
func (c *Service) GetBranches(projectID int) ([]*wrappedGitLab.Branch, error) {
	var branches []*wrappedGitLab.Branch

	c.branchesMtx.RLock()
	defer c.branchesMtx.RUnlock()

	if val, ok := c.branches[projectID]; ok {
		branches = val
	}

	return branches, nil
}

func (c *Service) UpdateBranches(projectIDs []int) error {
	branchesByProjectID := make(map[int][]*wrappedGitLab.Branch, len(projectIDs))
	for _, projectID := range projectIDs {
		var branches []*wrappedGitLab.Branch
		page := 1
		for {
			remoteBranches, resp, err := c.git.Branches.ListBranches(
				projectID,
				&wrappedGitLab.ListBranchesOptions{
					ListOptions: wrappedGitLab.ListOptions{PerPage: 100, Page: page},
				},
			)

			if err != nil {
				log.Println("Error when getting branches list, ", err)
				break
			}

			for _, branch := range remoteBranches {
				branches = append(branches, branch)
			}

			if page >= resp.TotalPages {
				break
			}

			page++
		}
		sort.Sort(ByCommitDateDesc(branches))
		branchesByProjectID[projectID] = branches
	}

	c.branchesMtx.Lock()
	c.branches = branchesByProjectID
	c.branchesMtx.Unlock()

	return nil
}

const (
	JobStatusPreparing = "preparing"
	JobStatusCreated   = "created"
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusFailed    = "failed"
	JobStatusSuccess   = "success"
	JobStatusCanceled  = "canceled"
	JobStatusSkipped   = "skipped"
	JobStatusManual    = "manual"
)

var finishedJobStatus = []string{
	JobStatusFailed,
	JobStatusSuccess,
	JobStatusCanceled,
	JobStatusSkipped,
	JobStatusManual,
}

var inProcessJobStatus = []string{
	JobStatusCreated,
	JobStatusPreparing,
	JobStatusPending,
	JobStatusRunning,
}

var neverStartedJobStatus = []string{
	JobStatusManual,
}

// PlayOrRetryJob play a job or retries a job for given criteria
// Affected job will be tracker by a watcher until finished status
// Affected job will be places in job list (Service.jobs) forever
func (c *Service) PlayOrRetryJob(projectID int, environment string, ref string) (*wrappedGitLab.Job, error) {
	if utils.StringsContainString(c.protectedEnvironments, environment) {
		return nil, DeniedForProtectedEnvironment
	}

	job, err := c.findJobForGivenCriteriaRecursive(projectID, environment, ref)
	if err != nil {
		return nil, err
	}

	// Skip if it's already running
	if utils.StringsContainString(inProcessJobStatus, job.Status) {
		return nil, errors.New("job already running")
	}

	// Play or Retry jobs depends on current status
	var runJob *wrappedGitLab.Job
	if utils.StringsContainString(neverStartedJobStatus, job.Status) {
		runJob, _, err = c.git.Jobs.PlayJob(projectID, job.ID)
	} else {
		runJob, _, err = c.git.Jobs.RetryJob(projectID, job.ID)
	}
	if err != nil {
		return nil, err
	}

	// Store job to the job list
	c.jobsMtx.Lock()
	if _, ok := c.jobs[environment]; !ok {
		c.jobs[environment] = map[int]*wrappedGitLab.Job{}
	}
	c.jobs[environment][projectID] = runJob
	c.jobsMtx.Unlock()

	// Run watcher
	c.runJobWatcher(environment, projectID, runJob)

	return job, nil
}

// runJobWatcher checks the job and replace in jobs map in case of status changing
// It checks every 3 seconds
// When status became on of finished we stop the watcher
func (c *Service) runJobWatcher(environment string, projectId int, runJob *wrappedGitLab.Job) {
	go func() {
		for {
			watchedJob, _, err := c.git.Jobs.GetJob(projectId, runJob.ID)
			if err != nil {
				log.Println(err)
				return
			}

			// If the jobs changed the status we need to replace with a new one
			if runJob.Status != watchedJob.Status {
				c.jobsMtx.Lock()
				c.jobs[environment][projectId] = watchedJob
				c.jobsMtx.Unlock()
			}
			if utils.StringsContainString(finishedJobStatus, watchedJob.Status) {
				return
			}
			<-time.After(time.Second * 3)
		}
	}()
}

func (c *Service) GetDeployment(projectId, deploymentId int) (*Deployment, error) {
	deployment, _, err := c.git.Deployments.GetProjectDeployment(
		projectId,
		deploymentId,
	)
	return convertWrappedDeployment(deployment), err
}

func (c *Service) GetJob(environment string, projectID int) (*wrappedGitLab.Job, bool) {
	if utils.StringsContainString(c.protectedEnvironments, environment) {
		return nil, true
	}

	c.jobsMtx.RLock()
	defer c.jobsMtx.RUnlock()

	job, ok := c.jobs[environment][projectID]
	return job, ok
}

func (c *Service) ListProjectDeployments(environment string, projectID int) ([]*Deployment, error) {
	deployments, _, err := c.git.Deployments.ListProjectDeployments(
		projectID,
		&wrappedGitLab.ListProjectDeploymentsOptions{
			Environment: wrappedGitLab.String(environment),
			OrderBy:     wrappedGitLab.String("id"),
			Sort:        wrappedGitLab.String("desc"),
			Status:      wrappedGitLab.String(JobStatusSuccess),
			ListOptions: wrappedGitLab.ListOptions{
				PerPage: 5,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	wrappedDeployments := make([]*Deployment, len(deployments))
	for i, deployment := range deployments {
		wrappedDeployments[i] = convertWrappedDeployment(deployment)
	}
	return wrappedDeployments, nil
}

// UpdateEnvironments updates environments cache
func (c *Service) UpdateEnvironments(projectIds []int) error {
	environments := map[string]*Environment{}
	for _, projectId := range projectIds {
		remoteEnvironments, _, err := c.git.Environments.ListEnvironments(
			projectId,
			&wrappedGitLab.ListEnvironmentsOptions{PerPage: 100},
		)
		if err != nil {
			return err
		}

		for _, remoteEnvironment := range remoteEnvironments {
			// Skip protected environments
			if utils.StringsContainString(c.protectedEnvironments, remoteEnvironment.Name) {
				continue
			}
			// We store it because
			// GetEnvironment returns env without Project field :(
			remoteProject := remoteEnvironment.Project

			// We need to fetch env by ID to get last deployment
			// because ListEnvironments doesn't return it
			remoteEnvironment, _, err = c.git.Environments.GetEnvironment(projectId, remoteEnvironment.ID)
			fmt.Println(remoteEnvironment.Name)
			fmt.Println(projectId)
			if err != nil {
				return err
			}

			// Put the project back
			remoteEnvironment.Project = remoteProject
			// Add a new environment or add the remoteProject to the existed env
			if _, ok := environments[remoteEnvironment.Name]; !ok {
				environments[remoteEnvironment.Name] = covertWrappedEnvironment(remoteEnvironment)
			} else {
				env := environments[remoteEnvironment.Name]
				env.Projects = append(env.Projects, convertWrappedProject(remoteEnvironment.Project, remoteEnvironment.LastDeployment))
			}

		}
	}
	c.environmentsMtx.Lock()
	c.environments = environments
	c.environmentsMtx.Unlock()

	return nil
}

// GetEnvironments returns cached environments by UpdateEnvironments function
func (c *Service) GetEnvironments() []*Environment {
	var environments []*Environment

	c.environmentsMtx.RLock()
	for _, environment := range c.environments {
		environments = append(environments, environment)
	}
	c.environmentsMtx.RUnlock()

	return environments
}

func (c *Service) GetJobs() map[string]map[int]*wrappedGitLab.Job {
	jobs := map[string]map[int]*wrappedGitLab.Job{}
	c.jobsMtx.RLock()
	defer c.jobsMtx.RUnlock()

	for env, project := range c.jobs {
		for projectId, job := range project {
			if _, ok := jobs[env]; !ok {
				jobs[env] = map[int]*wrappedGitLab.Job{}
			}

			jobs[env][projectId] = job
		}
	}

	return jobs
}

func (c *Service) findJobForGivenCriteriaRecursive(projectId int, environment string, ref string) (*wrappedGitLab.Job, error) {
	job, err := c.findJobForGivenCriteria(projectId, environment, ref, 1, 0)

	// If we didn't find a job
	// Let's try to find it in previous pipelines
	// because some pipelines couldn't have environments (i.e. scheduled piplines)
	if err == JobNotFound {
		limit := c.jobRecursiveSearchLimit
		page := 0
		perPage := 10
		for page <= limit {
			limit -= 1
			job, err := c.findJobForGivenCriteria(projectId, environment, ref, perPage, page)
			if err == JobNotFound {
				page += 1
				continue
			}

			return job, err
		}
	}

	return job, err
}

func (c *Service) findJobForGivenCriteria(projectId int, environment string, ref string, perPage int, page int) (*wrappedGitLab.Job, error) {
	pipelines, _, err := c.git.Pipelines.ListProjectPipelines(projectId, &wrappedGitLab.ListProjectPipelinesOptions{
		Ref: &ref,
		ListOptions: wrappedGitLab.ListOptions{
			PerPage: perPage,
			Page:    page,
		},
		OrderBy: wrappedGitLab.String("id"),
		Sort:    wrappedGitLab.String("desc"),
	})
	if err != nil {
		return nil, err
	}

	for _, pipeline := range pipelines {
		jobs, _, err := c.git.Jobs.ListPipelineJobs(projectId, pipeline.ID, &wrappedGitLab.ListJobsOptions{
			ListOptions: wrappedGitLab.ListOptions{PerPage: 100},
			Scope: []wrappedGitLab.BuildStateValue{
				wrappedGitLab.Manual,
				wrappedGitLab.Success,
				wrappedGitLab.Created,
				wrappedGitLab.Failed,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, job := range jobs {
			if job.Name == environment {
				if job.Status == JobStatusCreated {
					return nil, JobIsNotReady
				}
				return job, nil
			}
		}
	}
	return nil, JobNotFound
}

func (c *Service) PlayOrRetryJobsWithQuery(environment string, query string) error {
	count := 0
	for _, projectId := range c.projectIDs {
		branches, _, err := c.git.Branches.ListBranches(projectId, &wrappedGitLab.ListBranchesOptions{
			ListOptions: wrappedGitLab.ListOptions{PerPage: 1},
			Search:      wrappedGitLab.String(fmt.Sprintf("^%s", query)),
		})
		if err != nil {
			return err
		}
		// We did't find a branch go next
		if len(branches) == 0 {
			continue
		}

		_, err = c.PlayOrRetryJob(projectId, environment, branches[0].Name)

		// It's a normal behavior If job is not ready
		// We just skip this project for query deployment
		// If jobs not found we skipp it as well
		// Because project couldn't have environment
		if err != nil && err != JobIsNotReady && err != JobNotFound {
			return err
		}
		if err == nil {
			count += 1
		}
	}
	if count == 0 {
		return errors.New("nothing was run")
	}

	return nil
}

// NewClient creates a new Service
func NewClient(gitLabToken, gitLabBaseURL string, protectedEnvironments []string, projectIDs []int) (*Service, error) {
	git, err := wrappedGitLab.NewClient(gitLabToken, wrappedGitLab.WithBaseURL(gitLabBaseURL))
	if err != nil {
		return nil, err
	}

	return &Service{
		git:                     git,
		environments:            map[string]*Environment{},
		environmentsMtx:         sync.RWMutex{},
		branches:                map[int][]*wrappedGitLab.Branch{},
		branchesMtx:             sync.RWMutex{},
		jobs:                    map[string]map[int]*wrappedGitLab.Job{},
		jobsMtx:                 sync.RWMutex{},
		protectedEnvironments:   protectedEnvironments,
		projectIDs:              projectIDs,
		jobRecursiveSearchLimit: 10,
	}, nil
}
