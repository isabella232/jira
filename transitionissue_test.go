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

func TestTransitionIssue(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("wanted POST but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/issue/JIRA-5580/transitions" {
			t.Fatalf("Want /rest/api/2/issue/JIRA-5580/transitions but got %s\n", url.Path)
		}
		if r.Header.Get("Content-type") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Content-type"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading POST body: %v\n", err)
		}

		fmt.Printf("@@@ %s\n", string(data))

		var issue Issue
		if err := json.Unmarshal(data, &issue); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}

		fmt.Printf("%+v\n", issue)

		if issue.Transition.ID != "341" {
			t.Fatalf("Want 341 but got %s\n", issue.Transition.ID)
		}

		if issue.Fields.FixVersions[0].Name != "1.0" {
			t.Fatalf("Want 1.0 but got %s\n", issue.Fields.FixVersions[0].Name)
		}

		w.WriteHeader(204)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	rc, err := client.TransitionIssue("JIRA-5580", "341", "1.0")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if (rc != http.StatusNoContent) {
		t.Fatalf("Excpected 204, but was %s", rc)
	}
}
