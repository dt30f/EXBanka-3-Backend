package models

type EmployeePermission struct {
	ID         int64  `json:"id"`
	EmployeeID int64  `json:"employee_id"`
	Permission string `json:"permission"`
}