package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetIssue(t *testing.T) {
	response := `
{
  "expand": "renderedFields,names,schema,transitions,operations,editmeta,changelog",
  "id": "176587",
  "self": "https://example.com/rest/api/2/issue/176587",
  "key": "JIRA-5580",
  "fields": {
    "customfield_13100": null,
    "customfield_15400": {
      "self": "https://example.com/rest/api/2/customFieldOption/26405",
      "value": "No",
      "id": "26405"
    },
    "customfield_13102": null,
    "customfield_13101": [
      {
        "self": "https://example.com/rest/api/2/user?username=jdoe",
        "name": "jdoe",
        "key": "jdoe",
        "emailAddress": "jdoe@example.com",
        "avatarUrls": {
          "48x48": "https://example.com/secure/useravatar?avatarId=10122",
          "24x24": "https://example.com/secure/useravatar?size=small&avatarId=10122",
          "16x16": "https://example.com/secure/useravatar?size=xsmall&avatarId=10122",
          "32x32": "https://example.com/secure/useravatar?size=medium&avatarId=10122"
        },
        "displayName": "Jon Doe",
        "active": true,
        "timeZone": "America/Guatemala"
      }
    ],
    "resolution": {
      "self": "https://example.com/rest/api/2/resolution/10000",
      "id": "10000",
      "description": "",
      "name": "Done"
    },
    "lastViewed": null,
    "labels": [],
    "aggregatetimeoriginalestimate": null,
    "issuelinks": [
      {
        "id": "92424",
        "self": "https://example.com/rest/api/2/issueLink/92424",
        "type": {
          "id": "10300",
          "name": "Test Case",
          "inward": "Is a Test Case of",
          "outward": "Test Case",
          "self": "https://example.com/rest/api/2/issueLinkType/10300"
        },
        "outwardIssue": {
          "id": "186599",
          "key": "JIRA-6785",
          "self": "https://example.com/rest/api/2/issue/186599",
          "fields": {
            "summary": "A nice summary of the issue",
            "status": {
              "self": "https://example.com/rest/api/2/status/10044",
              "description": "synapseRT Passed status",
              "iconUrl": "https://example.com/download/resources/com.go2group.jira.plugin.synapse/synapse/css/images/icons/status_passed.gif",
              "name": "Passed",
              "id": "10044",
              "statusCategory": {
                "self": "https://example.com/rest/api/2/statuscategory/3",
                "id": 3,
                "key": "done",
                "colorName": "green",
                "name": "Done"
              }
            },
            "issuetype": {
              "self": "https://example.com/rest/api/2/issuetype/10101",
              "id": "10101",
              "description": "For Go2Group SYNAPSE Test Case issue type",
              "iconUrl": "https://jirabeta.xoom.com/images/icons/genericissue.gif",
              "name": "Test Case",
              "subtask": false
            }
          }
        }
      }
    ],
    "assignee": {
      "self": "https://example.com/rest/api/2/user?username=jsmith",
      "name": "jsmith",
      "key": "jsmith",
      "emailAddress": "jsmith@example.com",
      "avatarUrls": {
        "48x48": "https://example.com/secure/useravatar?avatarId=10122",
        "24x24": "https://example.com/secure/useravatar?size=small&avatarId=10122",
        "16x16": "https://example.com/secure/useravatar?size=xsmall&avatarId=10122",
        "32x32": "https://example.com/secure/useravatar?size=medium&avatarId=10122"
      },
      "displayName": "Jane Smith",
      "active": true,
      "timeZone": "America/Los_Angeles"
    },
    "components": [
      {
        "self": "https://example.com/rest/api/2/component/15053",
        "id": "15053",
        "name": "trunk",
        "description": "https://git.example.com/projects/PROJ/repos/program/browse"
      }
    ],
    "subtasks": [
      {
        "id": "186598",
        "key": "JIRA-6784",
        "self": "https://example.com/rest/api/2/issue/186598",
        "fields": {
          "summary": "Do the needful.",
          "status": {
            "self": "https://example.com/rest/api/2/status/5",
            "description": "A resolution has been taken, and it is awaiting verification by reporter. From here issues are either reopened, or are closed.",
            "iconUrl": "https://example.com/images/icons/statuses/resolved.png",
            "name": "Resolved",
            "id": "5",
            "statusCategory": {
              "self": "https://example.com/rest/api/2/statuscategory/3",
              "id": 3,
              "key": "done",
              "colorName": "green",
              "name": "Done"
            }
          },
          "issuetype": {
            "self": "https://example.com/rest/api/2/issuetype/5",
            "id": "5",
            "description": "The sub-task of the issue",
            "iconUrl": "https://example.com/images/icons/issuetypes/subtask_alternate.png",
            "name": "Sub-task",
            "subtask": true
          }
        }
      }
    ],
    "reporter": {
      "self": "https://example.com/rest/api/2/user?username=jdoe",
      "name": "jdoe",
      "key": "jdoe",
      "emailAddress": "jdoe@example.com",
      "avatarUrls": {
        "48x48": "https://example.com/secure/useravatar?avatarId=10122",
        "24x24": "https://example.com/secure/useravatar?size=small&avatarId=10122",
        "16x16": "https://example.com/secure/useravatar?size=xsmall&avatarId=10122",
        "32x32": "https://example.com/secure/useravatar?size=medium&avatarId=10122"
      },
      "displayName": "Jon Doe",
      "active": true,
      "timeZone": "America/Los_Angeles"
    },
    "customfield_12100": [
      {
        "self": "https://example.com/rest/api/2/customFieldOption/30066",
        "value": "JIRA-6785",
        "id": "30066"
      },
      {
        "self": "https://example.com/rest/api/2/customFieldOption/30069",
        "value": "JIRA-6792",
        "id": "30069"
      }
    ],
    "customfield_10840": 0,
    "customfield_10838": "Please Provide",
    "progress": {
      "progress": 0,
      "total": 0
    },
    "worklog": {
      "startAt": 0,
      "maxResults": 20,
      "total": 0,
      "worklogs": []
    },
    "issuetype": {
      "self": "https://example.com/rest/api/2/issuetype/7",
      "id": "7",
      "description": "Created by JIRA Agile - do not edit or delete. Issue type for a user story.",
      "iconUrl": "https://example.com/images/icons/issuetypes/story.png",
      "name": "Story",
      "subtask": false
    },
    "project": {
      "self": "https://example.com/rest/api/2/project/14203",
      "id": "14203",
      "key": "JIRA",
      "name": "JIRA Project",
      "avatarUrls": {
        "48x48": "https://example.com/secure/projectavatar?avatarId=10011",
        "24x24": "https://example.com/secure/projectavatar?size=small&avatarId=10011",
        "16x16": "https://example.com/secure/projectavatar?size=xsmall&avatarId=10011",
        "32x32": "https://example.com/secure/projectavatar?size=medium&avatarId=10011"
      },
      "projectCategory": {
        "self": "https://example.com/rest/api/2/projectCategory/10200",
        "id": "10200",
        "description": "",
        "name": "Engineering- Microservices"
      }
    },
    "customfield_13300": "",
    "resolutiondate": "2016-04-13T13:28:02.000-0700",
    "watches": {
      "self": "https://example.com/rest/api/2/issue/JIRA-5580/watchers",
      "watchCount": 4,
      "isWatching": true
    },
    "customfield_16000": "2016-04-12 00:00:00.0",
    "customfield_14500": {
      "self": "https://example.com/rest/api/2/customFieldOption/18926",
      "value": "No",
      "id": "18926"
    },
    "updated": "2016-04-14T17:55:47.000-0700",
    "timeoriginalestimate": null,
    "description": "An excellent description of the work to be done.",
    "timetracking": {},
    "summary": "A concise summary",
    "environment": null,
    "duedate": null,
    "comment": {
      "startAt": 0,
      "maxResults": 0,
      "total": 0,
      "comments": []
    },
    "fixVersions": [
      {
        "self": "https://example.com/rest/api/2/version/17261",
        "id": "17261",
        "name": "PROJ.1.0.1",
        "archived": true,
        "released": true,
        "releaseDate": "2016-04-12"
      }
    ],
    "timeestimate": null,
    "versions": [],
    "status": {
      "self": "https://example.com/rest/api/2/status/10455",
      "description": "",
      "iconUrl": "https://example.com/images/icons/statuses/generic.png",
      "name": "Resolved",
      "id": "10455",
      "statusCategory": {
        "self": "https://example.com/rest/api/2/statuscategory/3",
        "id": 3,
        "key": "done",
        "colorName": "green",
        "name": "Done"
      }
    },
    "aggregatetimeestimate": 0,
    "creator": {
      "self": "https://example.com/rest/api/2/user?username=jdoe",
      "name": "jdoe",
      "key": "jdoe",
      "emailAddress": "jdoe@example.com",
      "avatarUrls": {
        "48x48": "https://example.com/secure/useravatar?avatarId=10122",
        "24x24": "https://example.com/secure/useravatar?size=small&avatarId=10122",
        "16x16": "https://example.com/secure/useravatar?size=xsmall&avatarId=10122",
        "32x32": "https://example.com/secure/useravatar?size=medium&avatarId=10122"
      },
      "displayName": "Jon Doe",
      "active": true,
      "timeZone": "America/Los_Angeles"
    },
    "aggregateprogress": {
      "progress": 18000,
      "total": 18000,
      "percent": 100
    },
    "timespent": null,
    "aggregatetimespent": 18000,
    "workratio": -1,
    "created": "2016-03-08T16:44:58.000-0800",
    "attachment": []
  }
}
`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/issue/JIRA-5580" {
			t.Fatalf("Want /rest/api/2/issue/JIRA-5580 but got %s\n", url.Path)
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
	issue, err := client.GetIssue("JIRA-5580")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if issue.ID != "176587" {
		t.Fatalf("Want 176587 but got %s\n", issue.ID)
	}

	if (issue.Fields.Status.Name != "Resolved") {
		t.Fatalf("Want Resolved but got %s\n", issue.Fields.Status.Name)
	}
}
