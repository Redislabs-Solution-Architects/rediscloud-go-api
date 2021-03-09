# rediscloud-go-api

This repository is a Go SDK for the [Redis Cloud REST API](https://docs.redislabs.com/latest/rc/api/).

The cloudformation-proposed-changes branch contains the changes proposed to integrate with CloudFormation using resource types.

Unfortunately that proposal is on hold due to issues on the AWS CloudFormation side. The basic issue being that, as of 2021-03, the testing path doesn't work and the handler API is scantily documented. The consequence being that it will take some time to complete, so better to wait until the CloudFormation code is more mature and documented.

## Getting Started

### Installing
You can use this module by using `go get` to add it to either your `GOPATH` workspace or
the project's dependencies.
```shell script
go get github.com/RedisLabs/rediscloud-go-api
```

### Example
This is an example of using the SDK
```go
package main

import (
	"context"
	"fmt"

	rediscloud_api "github.com/RedisLabs/rediscloud-go-api"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
)

func main() {
	// The client will use the credentials from `REDISCLOUD_ACCESS_KEY` and `REDISCLOUD_SECRET_KEY` by default
	client, err := rediscloud_api.NewClient()
	if err != nil {
		panic(err)
	}

	id, err := client.Subscription.Create(context.TODO(), subscriptions.CreateSubscription{
		// ...
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created subscription: %d", id)
}
```
