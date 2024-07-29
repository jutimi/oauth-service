package utils

import (
	"oauth-server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPayload struct {
	ID    uuid.UUID `json:"id"`
	Scope string    `json:"scopes"`
	jwt.RegisteredClaims
}

type WorkspacePayload struct {
	ID              uuid.UUID `json:"id"`
	Scope           string    `json:"scopes"`
	WorkspaceID     uuid.UUID `json:"workspace_id"`
	UserWorkspaceID uuid.UUID `json:"user_workspace_id"`
	jwt.RegisteredClaims
}

/*
Parameters:

- payload: The data payload to be included in the token.

- key: The secret key used for signing the token.

- expireTime: The expiration time of the token in seconds.

Returns:

string: The generated token.

error: An error if the token generation fails.
*/
func GenerateToken(payload interface{}, key string, expireTime int) (string, error) {
	conf := config.GetConfiguration().Jwt

	claims := struct {
		data interface{}
		jwt.RegisteredClaims
	}{
		data: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime))),
			Issuer:    conf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func VerifyToken(tokenString string, key string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return token, nil
}

func ParseWSToken(tokenString string) (*WorkspacePayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WorkspacePayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*WorkspacePayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
