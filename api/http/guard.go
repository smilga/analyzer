package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/smilga/analyzer/api"
)

var (
	ErrInvalidToken = errors.New("Error invalid token")
)

type contextKey string

const (
	uidKey contextKey = "uid"
)

type Guard struct {
	*GuardConfig
}

func (g *Guard) Protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.isAllowedRoute(r.RequestURI) {
			h.ServeHTTP(w, r)
			return
		}

		id, err := g.authorize(r)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			res, err := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		ctx := context.WithValue(r.Context(), uidKey, id)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *Guard) isAllowedRoute(URI string) bool {
	for _, r := range g.GuardConfig.Allowed {
		if r == URI {
			return true
		}
	}
	return false
}

func (g *Guard) authorize(r *http.Request) (api.UserID, error) {
	token, err := parseBearerHeader(r)
	if err != nil {
		return 0, err
	}

	valid, ID, err := g.Auth.Valid(token)
	if err != nil {
		return 0, err
	}

	if !valid {
		return 0, err
	}

	return ID, nil
}

func parseBearerHeader(r *http.Request) (string, error) {
	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) > 1 && fields[1] != "null" {
		return fields[1], nil
	}
	return "", fmt.Errorf("authorization header invalid")
}

type GuardConfig struct {
	Auth    api.Auth
	Allowed []string
}

func NewGuard(config *GuardConfig) *Guard {
	return &Guard{config}
}
