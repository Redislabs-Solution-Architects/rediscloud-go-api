package rediscloud_api

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/kvstore"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/cloud_accounts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TestCreateInitialization that the CreateWithTask returns the primary id when the creation is started.
func TestCreateInitialization(t *testing.T) {
	expected := 0
	primaryID := "f15aced1-187c-4ff7-8109-ff82168a8f45"
	s := httptest.NewServer(infiniteTestServer("key", "secret", postRequest(t, "/cloud-accounts", `{
  "accessKeyId": "123456",
  "accessSecretKey": "765432",
  "consoleUsername": "foo",
  "consolePassword": "bar",
  "name": "cumulus nimbus",
  "provider": "AWS",
  "signInLoginUrl": "http://example.org/foo"
}`, fmt.Sprintf(`{
  "taskId": "%v",
  "commandType": "cloudAccountCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`, primaryID)),
		getCloudAccountCreateInProgress(t, primaryID, expected)))

	subject, err := clientFromTestServerV2(s, "key", "secret")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()

	taskID := ""
	actual, err := subject.CloudAccount.CreateWithTask(ctx, cloud_accounts.CreateCloudAccount{
		AccessKeyID:     redis.String("123456"),
		AccessSecretKey: redis.String("765432"),
		ConsoleUsername: redis.String("foo"),
		ConsolePassword: redis.String("bar"),
		Name:            redis.String("cumulus nimbus"),
		Provider:        redis.String("AWS"),
		SignInLoginURL:  redis.String("http://example.org/foo"),
	}, &taskID)
	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
	assert.Equal(t, expected, actual)
	assert.Equal(t, primaryID, taskID)
}

// TestCreateContinuation tests that intermediate create calls return a 0 alternative id
func TestCreateContinuation(t *testing.T) {
	alternativeID := 0
	primaryID := "f15aced1-187c-4ff7-8109-ff82168a8f45"
	s := httptest.NewServer(infiniteTestServer("key", "secret",
		getRequest(t, "/tasks/"+primaryID, fmt.Sprintf(`{
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
}`, primaryID, alternativeID))))

	subject, err := clientFromTestServerV2(s, "key", "secret")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	taskID := primaryID
	alternateID, err := subject.CloudAccount.CreateWithTask(ctx, cloud_accounts.CreateCloudAccount{
		AccessKeyID:     redis.String("123456"),
		AccessSecretKey: redis.String("765432"),
		ConsoleUsername: redis.String("foo"),
		ConsolePassword: redis.String("bar"),
		Name:            redis.String("cumulus nimbus"),
		Provider:        redis.String("AWS"),
		SignInLoginURL:  redis.String("http://example.org/foo"),
	}, &taskID)

	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
		assert.NotEqual(t, 0, alternateID)
	}
	// However the task finished the taskID should be primaryID
	assert.Equal(t, primaryID, taskID)
}

// TestCreateComplete tests that the final create call returns the alternative id
func TestCreateComplete(t *testing.T) {
	alternativeID := 12345
	primaryID := "f15aced1-187c-4ff7-8109-ff82168a8f45"
	s := httptest.NewServer(testServer("key", "secret",
		getCloudAccountCreateCompleted(t, primaryID, alternativeID)))

	subject, err := clientFromTestServerV2(s, "key", "secret")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	taskID := primaryID
	actual, err := subject.CloudAccount.CreateWithTask(ctx, cloud_accounts.CreateCloudAccount{
		AccessKeyID:     redis.String("123456"),
		AccessSecretKey: redis.String("765432"),
		ConsoleUsername: redis.String("foo"),
		ConsolePassword: redis.String("bar"),
		Name:            redis.String("cumulus nimbus"),
		Provider:        redis.String("AWS"),
		SignInLoginURL:  redis.String("http://example.org/foo"),
	}, &taskID)
	require.NoError(t, err)
	assert.Equal(t, alternativeID, actual)
	assert.Equal(t, primaryID, taskID)
}

//TestDeleteInitiation initiates the deletion task
func TestDeleteWithInitiation(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	deleteRequestTaskID := "DELETE_ID"
	s := httptest.NewServer(infiniteTestServer("apiKey", "secret",
		deleteRequest(t, fmt.Sprintf("/cloud-accounts/%d", resourceID), fmt.Sprintf(`{
  "taskId": "%v",
  "commandType": "cloudAccountDeleteRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "DELETE"
    }
  }
}`, deleteRequestTaskID)),
		getCloudAccountDeleteInProgress(t, deleteRequestTaskID)))

	kvstore := kvstore.NewKVMap()
	kvstore.Put(primaryID, resourceID)
	subject, err := clientFromTestServerV2(s, "apiKey", "secret", KVStoreOpt(kvstore))
	require.NoError(t, err)

	taskID := ""

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	err = subject.CloudAccount.DeleteWithPrimaryID(ctx, primaryID, &taskID)
	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
	assert.Equal(t, taskID, deleteRequestTaskID)
}

//TestDeleteCompletion tests that another call completes the deletion task
func TestDeleteCompletion(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	deleteRequestTaskID := "DELETE_ID"
	s := httptest.NewServer(testServer("apiKey", "secret",
		getRequest(t, "/tasks/"+deleteRequestTaskID, fmt.Sprintf(`{
    "taskId": "%v",
    "commandType": "cloudAccountDeleteRequest",
    "status": "processing-completed",
    "timestamp": "2020-10-28T09:58:16.798Z",
    "response": {
    },
    "_links": {
      "self": {
        "href": "https://example.com",
        "type": "GET"
      }
    }
  }`, deleteRequestTaskID))))

	subject, err := clientFromTestServerV2(s, "apiKey", "secret")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	err = subject.CloudAccount.DeleteWithPrimaryID(ctx, primaryID, &deleteRequestTaskID)
	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
}

//TestDeleteContinuation tests that deletion can be deferred until completion
func TestDeleteContinuation(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	deleteRequestTaskID := "DELETE_ID"
	s := httptest.NewServer(infiniteTestServer("apiKey", "secret",
		getCloudAccountDeleteInProgress(t, deleteRequestTaskID)))

	subject, err := clientFromTestServerV2(s, "apiKey", "secret")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	err = subject.CloudAccount.DeleteWithPrimaryID(ctx, primaryID, &deleteRequestTaskID)

	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
}

//TestListPrimaryIds tests that all current Cloud Account Primary IDs are returned
func TestListPrimaryIDs(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret"))
	kvstore := kvstore.NewKVMap()
	kvstore.Put("pid1", 1)
	kvstore.Put("pid2", 2)
	uut, err := clientFromTestServerV2(s, "apiKey", "secret", KVStoreOpt(kvstore))
	require.NoError(t, err)

	primaryIDs := uut.CloudAccount.ListByPrimaryID(context.TODO())
	expected := []string{"pid1", "pid2"}
	assert.ElementsMatch(t, expected, primaryIDs)

}

//TestReadCloudAccount tests that a cloud account can be read, or an error is returned.
func TestReadCloudAccount(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	s := httptest.NewServer(testServer("apiKey", "secret",
		getRequest(t,
			fmt.Sprintf("/cloud-accounts/%v", resourceID),
			fmt.Sprintf(`{
  "id": %v,
  "name": "Frank",
  "provider": "AWS",
  "status": "active",
  "accessKeyId": "keyId",
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`, resourceID))))
	kvstore := kvstore.NewKVMap()
	kvstore.Put(primaryID, resourceID)
	kvstore.Put("pid2", 2)
	subject, err := clientFromTestServerV2(s, "apiKey", "secret", KVStoreOpt(kvstore))
	require.NoError(t, err)

	actual, err := subject.CloudAccount.ReadByPrimaryID(context.TODO(), primaryID)
	require.NoError(t, err)

	assert.Equal(t, &cloud_accounts.CloudAccount{
		ID:          redis.Int(resourceID),
		Name:        redis.String("Frank"),
		Provider:    redis.String("AWS"),
		Status:      redis.String("active"),
		AccessKeyID: redis.String("keyId"),
	}, actual)
}

//TestReadCloudAccountError tests that a cloud account returns a *error when the cloud account doesn't exist.
func TestReadCloudAccountError(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	s := httptest.NewServer(testServer("apiKey", "secret",
		getRequestWithStatus(t, fmt.Sprintf("/cloud-accounts/%v", resourceID),
			404, // error injected here - 404 returned when the resourceID is not found
			fmt.Sprintf(`{
  "id": %v,
  "name": "Frank",
  "provider": "AWS",
  "status": "active",
  "accessKeyId": "keyId",
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`, resourceID))))

	subject, err := clientFromTestServerV2(s, "apiKey", "secret")
	require.NoError(t, err)

	_, err = subject.CloudAccount.ReadByPrimaryID(context.TODO(), primaryID)

	require.Error(t, err)
	//require.IsType(t, &cloud_accounts.NotFound{}, err)
}

//TestUpdateInitiation tests that an update initiates the updating process
func TestUpdateInitiation(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	updateTaskID := "update-task-id"
	s := httptest.NewServer(infiniteTestServer("key", "secret",
		putRequest(t, fmt.Sprintf("/cloud-accounts/%d", resourceID), `{
  "accessKeyId": "tfvbjuyg",
  "accessSecretKey": "gyujmnbvgy",
  "consoleUsername": "baz",
  "consolePassword": "bar",
  "name": "stratocumulus",
  "signInLoginUrl": "http://example.org/foo"
}`, fmt.Sprintf(`{
  "taskId": "%s",
  "commandType": "cloudAccountUpdateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`, updateTaskID)),
		getCloudAccountUpdateInProgress(t, updateTaskID, resourceID)))
	kvstore := kvstore.NewKVMap()
	kvstore.Put(primaryID, resourceID)
	subject, err := clientFromTestServerV2(s, "key", "secret", KVStoreOpt(kvstore))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	taskID := ""
	err = subject.CloudAccount.UpdateWithPrimaryID(ctx, cloud_accounts.CfnUpdateCloudAccount{
		PrimaryID: &primaryID,
		UpdateCloudAccount: &cloud_accounts.UpdateCloudAccount{
			AccessKeyID:     redis.String("tfvbjuyg"),
			AccessSecretKey: redis.String("gyujmnbvgy"),
			ConsoleUsername: redis.String("baz"),
			ConsolePassword: redis.String("bar"),
			Name:            redis.String("stratocumulus"),
			SignInLoginURL:  redis.String("http://example.org/foo"),
		},
	},
		&taskID)
	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
	assert.Equal(t, updateTaskID, taskID)
}

//TestUpdateContinuation tests that an update continues when not finished the updating process
func TestUpdateContinuation(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	updateTaskID := "update-task-id"
	s := httptest.NewServer(infiniteTestServer("key", "secret",
		getCloudAccountUpdateInProgress(t, updateTaskID, resourceID)))

	subject, err := clientFromTestServerV2(s, "key", "secret")
	require.NoError(t, err)

	taskID := "update-task-id"
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	err = subject.CloudAccount.UpdateWithPrimaryID(ctx, cloud_accounts.CfnUpdateCloudAccount{
		PrimaryID: &primaryID,
		UpdateCloudAccount: &cloud_accounts.UpdateCloudAccount{
			AccessKeyID:     redis.String("tfvbjuyg"),
			AccessSecretKey: redis.String("gyujmnbvgy"),
			ConsoleUsername: redis.String("baz"),
			ConsolePassword: redis.String("bar"),
			Name:            redis.String("stratocumulus"),
			SignInLoginURL:  redis.String("http://example.org/foo"),
		},
	},
		&taskID)

	if errors.Is(err, context.DeadlineExceeded) {
		assert.Equal(t, updateTaskID, taskID)
	} else {
		require.NoError(t, err)
	}
}

//TestUpdateCompletion tests that an update completes and returns finished == true
func TestUpdateCompletion(t *testing.T) {
	primaryID := "e02b40d6-1395-4861-a3b9-ecf829d835fd"
	resourceID := 98765
	updateTaskID := "update-task-id"
	s := httptest.NewServer(testServer("key", "secret",
		getCloudAccountUpdateCompleted(t, updateTaskID, resourceID)))

	subject, err := clientFromTestServerV2(s, "key", "secret")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()

	taskID := "update-task-id"
	err = subject.CloudAccount.UpdateWithPrimaryID(ctx, cloud_accounts.CfnUpdateCloudAccount{
		PrimaryID: &primaryID,
		UpdateCloudAccount: &cloud_accounts.UpdateCloudAccount{
			AccessKeyID:     redis.String("tfvbjuyg"),
			AccessSecretKey: redis.String("gyujmnbvgy"),
			ConsoleUsername: redis.String("baz"),
			ConsolePassword: redis.String("bar"),
			Name:            redis.String("stratocumulus"),
			SignInLoginURL:  redis.String("http://example.org/foo"),
		},
	},
		&taskID)
	require.NoError(t, err)
}
