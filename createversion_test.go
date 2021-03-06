package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateVersion(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("wanted POST but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/version" {
			t.Fatalf("Want /rest/api/2/version but got %s\n", url.Path)
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

		var v Version
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if v.Name != "1.0" {
			t.Fatalf("Want 1.0 but got %s\n", v.Name)
		}
		if v.Description != "Version 1.0" {
			t.Fatalf("Want Version 1.0 but got %s\n", v.Description)
		}
		if v.ProjectID != 1 {
			t.Fatalf("Want 1 but got %d\n", v.ProjectID)
		}
		if v.Archived {
			t.Fatalf("Want false\n")
		}
		if v.Released {
			t.Fatalf("Want false\n")
		}

		v.ID = "9999"

		data, err = json.Marshal(&v)
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", string(data))
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	version, err := client.CreateVersion("1", "1.0")
	if err != nil {
		t.Fatalf("Unexpected error:  %v\n", err)
	}
	if version.ID != "9999" {
		t.Fatalf("Want 9999 but got %s\n", version.ID)
	}
	if version.Name != "1.0" {
		t.Fatalf("Want 1.0 but got %s\n", version.Name)
	}
}

func TestCreateVersionNon201(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.CreateVersion("1", "1.0")
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
