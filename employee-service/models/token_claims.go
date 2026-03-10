package models

type TokenClaims struct {
	CredentialID int64  `json:"credential_id"`
	EmployeeID   int64  `json:"employee_id"`
	Email        string `json:"email"`
	TokenType    string `json:"token_type"`
}