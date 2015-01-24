package jira

import "net/url"

type (
	Jira interface {
		GetComponents(projectKey string) ([]Component, error)
		CreateVersion(projectKey, versionName string) error
		MapVersionToComponent(componentID, versionName string) error
	}

	DefaultClient struct {
		username string
		password string
		baseURL  url.URL
		Jira
	}

	Component struct {
		ID          string `json:"id"`
		Name        string `json:"name:`
		Description string `json:"description"`
	}

	Version struct {
		Name        string `json:"name:`
		Description string `json:"description"`
		Project     string `json:"project"`
		ProjectID   int    `json:"projectId"`
	}
)

func NewClient(username, password string, baseURL url.URL) Jira {
	return DefaultClient{username: username, password: password, baseURL: baseURL}
}

func (client DefaultClient) GetComponents(projectKey string) ([]Component, error) {
	l := make([]Component, 0)
	return l, nil
}

func (client DefaultClient) CreateVersion(projectKey, versionName string) error {
	return nil
}

func (client DefaultClient) MapVersionToComponent(componentID, versionName string) error {
	return nil
}
