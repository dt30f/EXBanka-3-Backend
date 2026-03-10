package models

type UpdateEmployeePermissionsRequest struct {
	Permissions []string `json:"permissions"`
}