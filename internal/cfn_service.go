package internal

import (
	"context"
)

// resourceType denotes the types of resources
type resourceType string

// Resource Types
const (
	cloudAccount resourceType = "cloudAccount"
)

//GetExistingCloudAccounts returns the primary IDs for Cloud Account resources, or an error
func (a *api) GetExistingCloudAccounts(ctx context.Context) (primaryIDs []string, err error) {
	return a.getExistingResources(ctx, cloudAccount)
}

//getExistingResources returns the primary IDs for resources of the given type
func (a *api) getExistingResources(ctx context.Context, resourceType resourceType) (taskIDs []string, err error) {
	var tasks struct {
		Tasks []task `json:"tasks"`
	}

	if err := a.client.Get(ctx, "retrieving all tasks", "/tasks/", &tasks); err != nil {
		return nil, err
	}

	creates := make(map[int]string)
	var deletes []int

	for _, t := range tasks.Tasks {
		if isCompletedCreate(t, resourceType) {
			creates[*t.Response.ID] = *t.ID
		} else {
			if isCompletedDelete(t, resourceType) {
				deletes = append(deletes, *t.Response.ID)
			}
		}
	}

	for _, d := range deletes {
		delete(creates, d)
	}
	for _, pid := range creates {
		taskIDs = append(taskIDs, pid)
	}
	return
}

func isCompletedCreate(t task, resourceType resourceType) bool {
	return *t.Status == processedState && *t.CommandType == string(resourceType)+"CreateRequest"
}

func isCompletedDelete(t task, resourceType resourceType) bool {
	return *t.Status == processedState && *t.CommandType == string(resourceType)+"DeleteRequest"
}
