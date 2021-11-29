package gitlab

import (
	wrappedGitlab "github.com/xanzy/go-gitlab"
)

func covertWrappedEnvironment(environment *wrappedGitlab.Environment) *Environment {
	return &Environment{
		Name:     environment.Name,
		Projects: []*Project{convertWrappedProject(environment.Project, environment.LastDeployment)},
	}
}

func convertWrappedProject(project *wrappedGitlab.Project, lastDeployment *wrappedGitlab.Deployment) *Project {
	if project == nil {
		return nil
	}
	return &Project{
		ID:                project.ID,
		Name:              project.Name,
		NameWithNamespace: project.NameWithNamespace,
		AvatarURL:         project.AvatarURL,
		WebURL:            project.WebURL,
		LastDeployment:    convertWrappedDeployment(lastDeployment),
	}
}

func convertWrappedDeployment(deployment *wrappedGitlab.Deployment) *Deployment {
	if deployment == nil {
		return nil
	}
	return &Deployment{
		ID:         deployment.ID,
		Ref:        deployment.Ref,
		User:       convertWrappedProjectUser(deployment.User),
		UpdatedAt:  deployment.UpdatedAt,
		Deployable: convertWrappedDeployable(deployment),
	}
}

func convertWrappedDeployable(deployment *wrappedGitlab.Deployment) *Deployable {
	return &Deployable{
		Pipeline: Pipeline{
			ID:     deployment.Deployable.Pipeline.ID,
			Status: deployment.Deployable.Status,
			User:   convertWrappedUser(deployment.Deployable.User),
		},
	}
}

func convertWrappedUser(user *wrappedGitlab.User) *User {
	if user != nil {
		return &User{
			Name:      user.Name,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		}
	}

	return &User{}
}

func convertWrappedProjectUser(user *wrappedGitlab.ProjectUser) *ProjectUser {
	if user == nil {
		return nil
	}
	return &ProjectUser{
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
	}
}
