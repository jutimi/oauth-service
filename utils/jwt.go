package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPayload struct {
	Id    uuid.UUID `json:"id"`
	Scope string    `json:"scopes"`
	jwt.RegisteredClaims
}

type WorkspacePayload struct {
	Id              uuid.UUID `json:"id"`
	Scope           string    `json:"scopes"`
	WorkspaceId     uuid.UUID `json:"workspace_id"`
	UserWorkspaceId uuid.UUID `json:"user_workspace_id"`
	jwt.RegisteredClaims
}

/*
Parameters:

- claims: The data payload to be included in the token.

- key: The secret key used for signing the token.

Returns:

string: The generated token.

error: An error if the token generation fails.
*/
func GenerateToken(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func VerifyUserToken(tokenString, key string) (*UserPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func VerifyWorkspaceToken(tokenString, key string) (*WorkspacePayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WorkspacePayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*WorkspacePayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func GenerateExpireTime(expireTime int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expireTime)))
}
