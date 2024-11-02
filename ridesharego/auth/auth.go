package auth

import (
	"errors"
	"main/core"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	sessionDuration = 24 * time.Hour
	roleAdmin       = "admin"
	roleUser        = "user"
)

var permissions = map[string]int{
	roleAdmin: 1,
	roleUser:  2,
}

func GetUserDetailsAndValidate(tokenString, role, secret string) (*core.UserAuth, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
		if time.Unix(exp, 0).Before(time.Now()) {
			return nil, errors.New("token expired")
		}

		requiredPermission, ok := permissions[role]
		if !ok {
			return nil, errors.New("invalid role")
		}

		requestRole := claims["aud"].(string)
		if requestRole == "" {
			return nil, errors.New("invalid role")
		}

		if permissions[requestRole] > requiredPermission {
			return nil, errors.New("insufficient role")
		}

		return &core.UserAuth{
			Role:   claims["aud"].(string),
			UserID: int(claims["sub"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token")
}

func createToken(user *core.User, issueTime time.Time, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": "rideshare-go",
		"aud": user.Role,
		"exp": issueTime.Add(sessionDuration).Unix(),
		"iat": issueTime.Unix(),
	})

	return claims.SignedString([]byte(secret))
}
