// JIRA API with Oguz Component Mappings
package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var debug bool

type (
	// https://docs.atlassian.com/jira/REST/latest/
	Jira interface {
		GetProject(projectKey string) (Project, error)
		GetComponents(projectID string) (map[string]Component, error)
		GetVersions(projectID string) (map[string]Version, error)
		CreateVersion(projectID, versionName string) (Version, error)
	}

	// http://jiraplugins.denizoguz.com/wp-content/uploads/2014/09/REST-Manual-v0.1.pdf
	ComponentVersions interface {
		GetMappings() error
		GetVersionsForComponent(projectID, componentID int) error
		UpdateReleaseDate(mappingID int, releaseDate string) error
		UpdateReleasedFlag(mappingID int, released bool) error
		CreateMapping(projectID string, componentID string, versionID string) error
		DeleteMapping(mappingID int) error
	}

	Project struct {
		ID string `json:"id"`
	}

	DefaultClient struct {
		username   string
		password   string
		baseURL    *url.URL
		httpClient *http.Client
		Jira
		ComponentVersions
	}

	Component struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	Version struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Project     string `json:"project"`
		ProjectID   int    `json:"projectId"`
		Archived    bool   `json:"archived"`
		Released    bool   `json:"released"`
		ReleaseDate string `json:"releaseDate"`
	}
)

func init() {
	debug = strings.ToLower(os.Getenv("JIRA_DEBUG")) == "true"
}

// NewClient returns a new default Jira client.
func NewClient(username, password string, baseURL *url.URL) Jira {
	return DefaultClient{username: username, password: password, baseURL: baseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

// GetProject returns a representation of a Jira project.
func (client DefaultClient) GetProject(projectKey string) (Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s", client.baseURL, projectKey), nil)
	if err != nil {
		return Project{}, err
	}
	if debug {
		log.Printf("jira.GetComponents URL %s\n", req.URL)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return Project{}, err
	}
	if responseCode != http.StatusOK {
		return Project{}, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var r Project
	if err := json.Unmarshal(data, &r); err != nil {
		return Project{}, err
	}

	return r, nil
}

// GetComponents returns a map of Component indexed by component name.
func (client DefaultClient) GetComponents(projectID string) (map[string]Component, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s/components", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("jira.GetComponents URL %s\n", req.URL)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var r []Component
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[string]Component)
	for _, c := range r {
		m[c.Name] = c
	}

	return m, nil
}

// GetVersions returns a map of Version indexed by version name.
func (client DefaultClient) GetVersions(projectID string) (map[string]Version, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s/versions", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("jira.GetVersions URL %s\n", req.URL)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var r []Version
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[string]Version)
	for _, c := range r {
		m[c.Name] = c
	}

	return m, nil
}

// CreateVersion creates a new version in Jira for the given project ID and version name.
func (client DefaultClient) CreateVersion(projectID, versionName string) (Version, error) {
	i, err := strconv.Atoi(projectID)
	if err != nil {
		return Version{}, err
	}
	version := Version{Name: versionName, Description: "Version " + versionName, ProjectID: i, Archived: false, Released: true, ReleaseDate: time.Now().Format("2006-01-02")}
	data, err := json.Marshal(&version)
	if err != nil {
		return Version{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/version", client.baseURL), bytes.NewBuffer(data))
	if debug {
		log.Printf("jira.CreateVersion URL %s\n", req.URL)
	}
	if err != nil {
		return Version{}, err
	}
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return Version{}, err
	}
	if responseCode != http.StatusCreated {
		return Version{}, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var v Version
	if err := json.Unmarshal(data, &v); err != nil {
		return Version{}, err
	}
	return v, nil
}

func (client DefaultClient) CreateMapping(projectID, componentID, versionID string) error {
	// POST http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/
	/*
		body:
			   {
			    "projectId":10000,
			    "componentId":10003,
			    "versionId":10001,
			    "released":false
			   }
	*/
	return nil
}

func (client DefaultClient) GetMappings() error {
	// GET http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/mappings
	return nil
}

func (client DefaultClient) GetVersionsForComponent(projectID, componentID string) error {
	// GET http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/applicable_versions?projectId=10000&projectKey=&selectedComponentIds=10000
	/*
	   [ { "description" : "Unknown",
	       "id" : -1,
	       "isReleased" : false,
	       "name" : "Unknown"
	     },
	     { "id" : 10001,
	       "isReleased" : true,
	       "name" : "v2"
	     },
	     { "id" : 10000,
	       "isReleased" : true,
	       "name" : "v1"
	     },
	     { "id" : 10002,
	       "isReleased" : true,
	       "name" : "v3"
	     }
	   ]
	*/
	return nil
}

func (client DefaultClient) UpdateReleaseDate(mappingID int, releaseDate string) error {
	// PUT http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/releaseDate/5?releaseDate=16%2FSep%2F14
	return nil
}

func (client DefaultClient) UpdateReleasedFlag(mappingID int, released bool) error {
	// PUT http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/releaseFlag/5?isReleased=true
	return nil
}

func (client DefaultClient) DeleteMapping(mappingID int) error {
	// DELETE http://localhost:2990/jira/rest/com.deniz.jira.mapping/latest/5
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
