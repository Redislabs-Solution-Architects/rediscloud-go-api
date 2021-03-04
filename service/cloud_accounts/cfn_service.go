package cloud_accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/kvstore"
)

type CloudAccountTask interface {
	Task
	GetExistingCloudAccounts(ctx context.Context) ([]string, error)
}

func NewAPIV2(client HttpClient, task CloudAccountTask, logger Log, kvstore kvstore.KVStore) *API {
	return &API{client: client, task: Task(task), logger: logger, kvstore: kvstore}
}

const longDuration time.Duration = 55 * time.Second
const shortDuration time.Duration = 25 * time.Second

// CreateWithTask will initiate or continue trying to create the given cloud account. If this is the initial call
// it will setup the taskId for completion on subsequent calls.
// The returned values will be
// - (0,DeadlineExceeded) for a successful call that is ongoing;
// - (non-zero, nil) for a created completion
// - (0, error) if there's an unrecoverable error.
func (a *API) CreateWithTask(ctx context.Context, account CreateCloudAccount, primaryID *string) (rid int, err error) {
	var response taskResponse
	if *primaryID == "" { // Initial create call
		if err = a.client.Post(ctx, "cloud account", "/cloud-accounts", account, &response); err != nil {
			return
		}
		*primaryID = *response.ID
	}
	// subsequent create calls
	rid, err = a.task.WaitForResourceId(ctx, *primaryID)
	if !errors.Is(err, context.DeadlineExceeded) {
		a.kvstore.Put(*primaryID, rid)
	}

	return
}

// DeleteWithPrimaryID will initiate or continue a cloud account deletion given a cloud account primary id
// Returns:
// - context.DeadlineExceeded if the deletion is to be continued; in which case the taskID will be set and must be used in subsequent calls
// - nil if the deletion is completed
// - error if the deletion failed.
func (a *API) DeleteWithPrimaryID(ctx context.Context, primaryID string, taskID *string) (err error) {
	var rid int

	if *taskID == "" {
		rid, err = a.getResourceID(ctx, primaryID)
		if err != nil {
			return err
		}
		var response taskResponse
		if err = a.client.Delete(ctx,
			fmt.Sprintf("delete cloud account %d", rid),
			fmt.Sprintf("/cloud-accounts/%d", rid),
			&response); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				return wrap404Error(rid, err)
			}
		}
		*taskID = *response.ID
	}

	a.logger.Printf("Waiting for cloud account %d to finish being deleted", rid)

	err = a.task.Wait(ctx, *taskID)
	if !errors.Is(err, context.DeadlineExceeded) {
		a.kvstore.Delete(primaryID)
	}
	return

}

//getResourceID returns the resource id given the primaryID
// Returns error if the taskID does not correspond to a successful creation and/or no resource id can be found
func (a *API) getResourceID(ctx context.Context, primaryID string) (id int, err error) {
	id = a.kvstore.Get(primaryID)
	if id == 0 {
		err = &NotFoundPrimary{primaryID}
	}
	return
}

//ListByPrimaryID returns a slice of cloud accounts by their primary ID (i.e. the id of the task that created them)
//Only existing cloud accounts are returned.
//An error is returned if errors are found
func (a *API) ListByPrimaryID(ctx context.Context) []string {
	return a.kvstore.Keys()
}

//ReadByPrimaryID returns a cloud account with the given primary ID, or an error
func (a *API) ReadByPrimaryID(ctx context.Context, primaryID string) (cloudAccount *CloudAccount, err error) {
	rid, err := a.getResourceID(ctx, primaryID)
	if err != nil {
		return
	}
	return a.Get(ctx, rid)
}

//UpdateWithPrimaryID updates the given CloudAccount.
//If taskID is null then sets taskID for continuation call.
//Error is context.DeadlineExceeded if continuation needed.
func (a *API) UpdateWithPrimaryID(ctx context.Context, cfnCloudAccount CfnUpdateCloudAccount, taskID *string) (err error) {
	var rid int

	if *taskID == "" {
		rid, err = a.getResourceID(ctx, *cfnCloudAccount.PrimaryID)
		if err != nil {
			return
		}
		var response taskResponse
		if err := a.client.Put(ctx, fmt.Sprintf("update cloud account %d", rid), fmt.Sprintf("/cloud-accounts/%d", rid), *cfnCloudAccount.UpdateCloudAccount, &response); err != nil {
			return wrap404ErrorPrimary(*cfnCloudAccount.PrimaryID, err)
		}
		*taskID = *response.ID
	}

	a.logger.Printf("Waiting for cloud account %d to finish being updated", rid)

	return a.task.Wait(ctx, *taskID)

}

func wrap404ErrorPrimary(primary string, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFoundPrimary{primary: primary}
	}
	return err
}
