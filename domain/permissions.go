package domain

import "errors"

type Permission string

func (p Permission) Validate() error {
	for _, permission := range Permissions {
		if permission == p {
			return nil
		}
	}
	return errors.New("Invalid permission")
}

const (
	CustomerPermission Permission = "customer"
	APIPermission      Permission = "api"
	AdminPermission    Permission = "admin"
)

var Permissions = []Permission{
	CustomerPermission,
	APIPermission,
	AdminPermission,
}
