package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type (
	Jira interface {
		GetComponents(projectKey string) ([]Component, error)
		CreateVersion(projectKey, versionName string) error
		MapVersionToComponent(componentID, versionName string) error
	}

	DefaultClient struct {
		username   string
		password   string
		baseURL    *url.URL
		httpClient *http.Client
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

func NewClient(username, password string, baseURL *url.URL) Jira {
	return DefaultClient{username: username, password: password, baseURL: baseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (client DefaultClient) GetComponents(projectKey string) ([]Component, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s/components", client.baseURL, projectKey), nil)
	if err != nil {
		return nil, err
	}
	log.Printf("jira.GetComponents URL %s\n", req.URL)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		var reason string = "unhandled reason"
		switch {
		case responseCode == http.StatusBadRequest:
			reason = "Bad request."
		}
		return nil, fmt.Errorf("Error getting project components: %s.  Status code: %d.  Reason: %s\n", string(data), responseCode, reason)
	}

	var r []Component
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (client DefaultClient) CreateVersion(projectKey, versionName string) error {
	return nil
}

func (client DefaultClient) MapVersionToComponent(componentID, versionName string) error {
	return nil
}

func (client DefaultClient) consumeResponse(req *http.Request) (rc int, buffer []byte, err error) {
	response, err := client.httpClient.Do(req)

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
		if e := recover(); e != nil {
			trace := make([]byte, 10*1024)
			_ = runtime.Stack(trace, false)
			log.Printf("%s", trace)
			err = fmt.Errorf("%v", e)
		}
	}()

	if err != nil {
		panic(err)
	}

	if data, err := ioutil.ReadAll(response.Body); err != nil {
		panic(err)
	} else {
		return response.StatusCode, data, nil
	}
}
