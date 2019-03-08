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
}

func NewHandler() *Handler {
	return &Handler{
		Auth:           NewJWTAuth(os.Getenv("JWT_SECRET")),
		ServiceStorage: inmemory.NewServiceStore(),
		WebsiteStorage: inmemory.NewWebsitesStore(),
		UserStorage:    inmemory.NewUserStore(),
	}
}
