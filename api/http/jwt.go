package http

import (
	"errors"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
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
		return false, 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	if userID, ok := claims["UserID"].(string); !ok {
		return false, 0, ErrParsingClaims
	} else {
		id, err := strconv.Atoi(userID)
		if err != nil {
			return false, 0, ErrParsingClaims
		}
		return true, api.UserID(id), nil
	}
}

// Sign returns new access token
func (a *JWTAuth) Sign(id api.UserID) (string, error) {
	claims := Claims{
		strconv.Itoa(int(id)),
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
