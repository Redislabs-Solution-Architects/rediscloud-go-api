# CloudFormation

The rediscloud-go-api was written with Terraform in
mind. Unfortunately the CloudFormation processing model is rather
different, and so some changes are required to this api if we are to
integrate with CloudFormation.

In particular, CloudFormation has [Resource type handler
contracts](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/resource-type-test-contract.html)
which dictate the overall flow between CloudFormation and any API
which it calls. Particularly troublesome are these provisions:

1. On Create, a stable resource ID (aka 'primary id') must be returned within 60 seconds
1. In any of Create, Update or Delete operations a reply must be
  returned within 60 seconds, along with any required
  callback. CloudFormation will make the exact same call at a time of
  its choosing, with the callback data, to continue the operation

Redis Labs Cloud API (CAPI) only returns a stable ID after the
resource is completely ready, which could be 10 minutes for a
Subscription.

Redis Labs Cloud API does return a task that can be tracked to see an
operation's progress and (although no guarantees are made in the API
spec) this normally takes no more than a few seconds at most to
return.

Given the above a technical solution (with limits) is:

1. Use the Creation/Update/Deletion task id as the callback value, so
   that the task can be followed.
1. Use the taskID for the Creation task as the primary id for all resources.

The limits here are:

1. Task IDs are cleared out after a couple of weeks. Their lifecycle
   is completely separate from the lifecycle of the resource they
   relate to.
1. There is no guarantee in the CAPI of response times. That being
   said, they're normally short enough that it doesn't matter.
1. There is no guarantee in the CAPI of how network disconnects are
   handled. This is not a robust distributed system. It is possible
   that the client will never hear back from the CAPI, and yet the
   resource has been modified according to the client's request.

Limit #1 above is serious, and means that the CAPI is of limited use
for CloudFormation. In private conversations with Danny Cohen (CAPI
owner) it would appear that it is possible to change the lifecycle of
the task IDs to match the corresponding resources.

# Proposed Changes
The changes proposed here are designed to ensure backwards
compatibility, as described in [Keeping Your Modules
Compatible](https://blog.golang.org/module-compatibility). Particularly
noteworthy are the following observations:

1. "There is, in fact, no backward-compatible change you can make to a
   function’s signature."
1. "Directly adding to an interface is a breaking change"
1. "If you have an exported struct type, you can almost always add a
   field or remove an unexported field without breaking
   compatibility."

The solutions proposed for these three limitations are (respectively):

1. Add a new function
1. Use the [Extension Interface Pattern](https://youtu.be/yx7lmuwUNv8?t=931)
1. Conservatively extend structures, being sure not to change their
   comparability semantic.
   
## Change Overview
From an initial experiment using just the ~cloud_accounts~ package it seems that the following changes will be needed:

1. add function: ~rediscloud_api.NewClientV2~
1. add function: ~rediscloud_api.cloud_accounts.NewClientV2~
1. Extend interface: ~rediscloud_api.cloud_accounts.CloudAccountTask~
1. Extend structure: ~redisclout_api.cloud_accounts.API~

Of these four only the latter is a modification to an original file. Everything else is additional to the current fileset.

We anticipate that a similar set of changes will be required for the other three services (account, databases, subscriptions). 

## Change Details
Each service will have new Create, Read, Update, Delete, List (CRUDL) methods defined. These are hooked into the basic setup in the client via the NewClientV2 functions, and the extended Task and API interfaces and structures.

The timeout functionality can be achieved using with the CUD operations by using the ~context.WithDeadline()~ feature, and then checking the returned error to see if the function errored out and whether the deadline was exceeded or not. Additionally these operations will take a pointer to a taskID, and that can be used as the callback data with the actual CloudFormation handlers.





