package jira

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func TestAddFixVersion(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("wanted PUT but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/issue/JIRA-5580" {
			t.Fatalf("Want /rest/api/2/issue/JIRA-5580 but got %s\n", url.Path)
		}
		if r.Header.Get("Content-type") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Content-type"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading PUT body: %v\n", err)
		}

		fmt.Printf("@@@ %s\n", string(data))

		var issue updateIssue
		if err := json.Unmarshal(data, &issue); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}

		fmt.Printf("%+v\n", issue)

		if issue.update.fixVersions[0].add != "1.0" {
			t.Fatalf("Want 1.0 but got %s\n", issue.update.fixVersions[0].add)
		}

		w.WriteHeader(204)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)

	rc, err := client.AddFixVersion("JIRA-5580", "1.0")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if rc != http.StatusOK {
		t.Fatalf("Expected 200, but got %d", rc)
	}
}
