package rediscloud_api

import (
	"fmt"
	"testing"
)

func getCloudAccountCreateCompleted(t *testing.T, primaryID string, resourceID int) endpointRequest {
	return getRequest(t, "/tasks/"+primaryID, fmt.Sprintf(`{
	"taskId": "%v",
	"commandType": "cloudAccountCreateRequest",
	"status": "processing-completed",
	"timestamp": "2020-10-28T09:58:16.798Z",
	"response": {
	  "resourceId": %d
	},
	"_links": {
	  "self": {
		"href": "https://example.com",
		"type": "GET"
	  }
	}
	}`, primaryID, resourceID))
}

func getCloudAccountCreateInProgress(t *testing.T, primaryID string, alternativeID int) endpointRequest {
	return getRequest(t, "/tasks/"+primaryID, fmt.Sprintf(`{
		"taskId": "%v",
		"commandType": "cloudAccountCreateRequest",
		"status": "processing-in-progress",
		"timestamp": "2020-11-02T09:05:35.1Z",
		"response": {
		  "resourceId": %d
		},
		"_links": {
		  "self": {
			"href": "https://example.com",
			"type": "GET"
		  }
		}
	  }`, primaryID, alternativeID))
}

func getCloudAccountUpdateCompleted(t *testing.T, taskID string, resourceID int) endpointRequest {
	return getRequest(t, fmt.Sprintf("/tasks/%s", taskID), fmt.Sprintf(`{
	"taskId": "%s",
	"commandType": "cloudAccountUpdateRequest",
	"status": "processing-completed",
	"timestamp": "2020-10-28T09:58:16.798Z",
	"response": {
		"resourceId": %d
	},
	"_links": {
	  "self": {
		"href": "https://example.com",
		"type": "GET"
	  }
	}
  }`, taskID, resourceID))
}

func getCloudAccountUpdateInProgress(t *testing.T, taskID string, resourceID int) endpointRequest {
	return getRequest(t, "/tasks/"+taskID, fmt.Sprintf(`{
    "taskId": "%v",
    "commandType": "cloudAccountUpdateRequest",
    "status": "processing-in-progress",
    "timestamp": "2020-11-02T09:05:35.1Z",
    "response": {
      "resourceId": %d
    },
    "_links": {
      "self": {
        "href": "https://example.com",
        "type": "GET"
      }
    }
  }`, taskID, resourceID))
}

func getCloudAccountDeleteInProgress(t *testing.T, taskID string) endpointRequest {
	return getRequest(t, "/tasks/"+taskID, fmt.Sprintf(`{
		"taskId": "%v",
		"commandType": "cloudAccountDeleteRequest",
		"status": "processing-in-progress",
		"timestamp": "2020-10-28T09:58:16.798Z",
		"response": {
		},
		"_links": {
		  "self": {
			"href": "https://example.com",
			"type": "GET"
		  }
		}
	  }`, taskID))
}
