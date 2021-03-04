package cloud_accounts

import "fmt"

//CfnCloudAccount is the regular CloudAccount wrapped with the create task ID
// The enclosed CloudAccount might have a 0 ID which requires a look up using the PrimaryID to resolve it.
type CfnCloudAccount struct {
	PrimaryID    *string `json:"omitempty"`
	CloudAccount *CloudAccount
}

//CfnUpdateCloudAccount is a wrapper around UpdateCloudAccount
type CfnUpdateCloudAccount struct {
	PrimaryID          *string
	UpdateCloudAccount *UpdateCloudAccount
}

//NotFoundPrimary is an error indicating that a primary ID was not found
type NotFoundPrimary struct {
	primary string
}

func (f *NotFoundPrimary) Error() string {
	return fmt.Sprintf("cloud account %s not found", f.primary)
}
