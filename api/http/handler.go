package http

import (
	"os"

	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/datastore/inmemory"
)

type Handler struct {
	Auth           api.Auth
	ServiceStorage api.ServiceStorage
	WebsiteStorage api.WebsiteStorage
	UserStorage    api.UserStorage
	Analyzer       *api.Analyzer
}

func NewHandler() *Handler {
	return &Handler{
		Auth:           NewJWTAuth(os.Getenv("JWT_SECRET")),
		ServiceStorage: inmemory.NewServiceStore(),
		WebsiteStorage: inmemory.NewWebsiteStore(),
		UserStorage:    inmemory.NewUserStore(),
		Analyzer: &api.Analyzer{
			inmemory.NewResultStore(),
			inmemory.NewServiceStore(),
		},
	}
}
