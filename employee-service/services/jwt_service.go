package services

import (
	"errors"

	"employee-service/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		Secret: secret,
	}
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	tokenType, _ := claimsMap["token_type"].(string)
	if tokenType != "access" {
		return nil, errors.New("provided token is not an access token")
	}

	credentialID, _ := claimsMap["credential_id"].(float64)
	employeeID, _ := claimsMap["employee_id"].(float64)
	email, _ := claimsMap["email"].(string)

	return &models.TokenClaims{
		CredentialID: int64(credentialID),
		EmployeeID:   int64(employeeID),
		Email:        email,
		TokenType:    tokenType,
	}, nil
}