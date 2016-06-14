package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetTransitions(t *testing.T) {
	response := `
{
  "expand": "transitions",
  "transitions": [
    {
      "id": "391",
      "name": "Close",
      "to": {
        "self": "https://example.com/rest/api/2/status/6",
        "description": "The issue is considered finished, the resolution is correct. Issues which are closed can be reopened.",
        "iconUrl": "https://example.com/images/icons/statuses/closed.png",
        "name": "Closed",
        "id": "6",
        "statusCategory": {
          "self": "https://example.com/rest/api/2/statuscategory/3",
          "id": 3,
          "key": "done",
          "colorName": "green",
          "name": "Done"
        }
      }
    },
    {
      "id": "341",
      "name": "Merge To Master",
      "to": {
        "self": "https://example.com/rest/api/2/status/10455",
        "description": "",
        "iconUrl": "https://example.com/images/icons/statuses/generic.png",
        "name": "Merged To Master",
        "id": "10455",
        "statusCategory": {
          "self": "https://example.com/rest/api/2/statuscategory/3",
          "id": 3,
          "key": "done",
          "colorName": "green",
          "name": "Done"
        }
      }
    }
  ]
}
`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/issue/JIRA-5580/transitions" {
			t.Fatalf("Want /rest/api/2/issue/JIRA-5580/transitions but got %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		fmt.Fprintln(w, response)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	transitions, err := client.GetTransitions("JIRA-5580")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if transitions.AvailableTransitions[1].ID != "341" {
		t.Fatalf("Want 341 but got %s\n", transitions.AvailableTransitions[1].ID)
	}

	if transitions.AvailableTransitions[1].Name != "Merge To Master" {
		t.Fatalf("Want Resolved but got %s\n", transitions.AvailableTransitions[1].Name)
	}
}
