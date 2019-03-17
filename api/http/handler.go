package http

import (
	"net/http"
	"os"

	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/datastore/inmem"
)

type Handler struct {
	Auth           api.Auth
	WebsiteStorage api.WebsiteStorage
	UserStorage    api.UserStorage
	PatternStorage api.PatternStorage
	TagStorage     api.TagStorage
	FilterStorage  api.FilterStorage
	Analyzer       *api.Analyzer
}

func (h *Handler) AuthID(r *http.Request) (api.UserID, error) {
	id := r.Context().Value(uidKey)
	uid, ok := id.(api.UserID)
	if !ok {
		return uid, api.ErrTokenError
	}

	return uid, nil
}

func NewTestHandler() *Handler {
	return &Handler{
		Auth:           NewJWTAuth(os.Getenv("JWT_SECRET")),
		WebsiteStorage: inmem.NewWebsiteStore(),
		UserStorage:    inmem.NewUserStore(),
		PatternStorage: inmem.NewPatternStore(),
		TagStorage:     inmem.NewTagStore(),
		FilterStorage:  inmem.NewFilterStore(),
		Analyzer: &api.Analyzer{
			PatternStorage: inmem.NewPatternStore(),
			WebsiteStorage: inmem.NewWebsiteStore(),
		},
	}
}
