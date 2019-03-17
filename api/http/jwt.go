package http

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

// Error definitions
var (
	ErrParsingClaims = errors.New("Error parsing token claims")
)

type Claims struct {
	UserID string
	jwt.StandardClaims
}

// JWTAuth uses jwt token for authentification
type JWTAuth struct {
	Secret string
}

// Valid returns if auth token is valid
func (a *JWTAuth) Valid(tokenString string) (bool, api.UserID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})

	if err != nil {
		return false, api.UserID{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	if userID, ok := claims["UserID"].(string); !ok {
		return false, api.UserID{}, ErrParsingClaims
	} else {
		uid, err := uuid.FromString(userID)
		if err != nil {
			return false, api.UserID{}, err
		}

		return true, api.UserID(uid), nil
	}
}

// Sign returns new access token
func (a *JWTAuth) Sign(ID api.UserID) (string, error) {
	claims := Claims{
		ID.String(),
		jwt.StandardClaims{
			Issuer: "Analyzer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.Secret))
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{secret}
}
