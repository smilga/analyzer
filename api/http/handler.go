package http

import (
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/comm"
	"github.com/smilga/analyzer/api/datastore/cache"
	"github.com/smilga/analyzer/api/datastore/mysql"
	"github.com/smilga/analyzer/api/ws"
)

type Handler struct {
	Auth           api.Auth
	WebsiteStorage api.WebsiteStorage
	UserStorage    api.UserStorage
	PatternStorage api.PatternStorage
	TagStorage     api.TagStorage
	FeatureStorage api.FeatureStorage
	ServiceStorage api.ServiceStorage
	FilterStorage  api.FilterStorage
	ReportStorage  api.ReportStorage
	ResultStorage  api.ResultStorage
	Messanger      *ws.Messanger
	Comm           *comm.Comm
}

func (h *Handler) AuthID(r *http.Request) (api.UserID, error) {
	id := r.Context().Value(uidKey)
	uid, ok := id.(api.UserID)
	if !ok {
		return uid, api.ErrTokenError
	}

	return uid, nil
}

func NewHandler(db *gorm.DB, comm *comm.Comm) *Handler {
	wstore := mysql.NewWebsiteStore(db)
	ps := mysql.NewPatternStore(db.DB())
	us := mysql.NewUserStore(db)
	ts := mysql.NewTagStore(db)
	rs := mysql.NewReportStore(db.DB())
	fs := mysql.NewFilterStore(db.DB())
	xs := mysql.NewFeatureStore(db)
	ss := mysql.NewServiceStore(db)
	rss := mysql.NewResultStore(db)
	pcache := cache.NewPatternCache(ps)

	return &Handler{
		Auth:           NewJWTAuth(os.Getenv("JWT_SECRET")),
		WebsiteStorage: wstore,
		UserStorage:    us,
		PatternStorage: pcache,
		TagStorage:     ts,
		FeatureStorage: xs,
		ServiceStorage: ss,
		FilterStorage:  fs,
		ReportStorage:  rs,
		Messanger:      ws.NewMessanger(),
		ResultStorage:  rss,
		Comm:           comm,
	}
}

// func NewTestHandler() *Handler {
// 	return &Handler{
// 		Auth:           NewJWTAuth(os.Getenv("JWT_SECRET")),
// 		WebsiteStorage: inmem.NewWebsiteStore(),
// 		UserStorage:    inmem.NewUserStore(),
// 		PatternStorage: inmem.NewPatternStore(),
// 		TagStorage:     inmem.NewTagStore(),
// 		FilterStorage:  inmem.NewFilterStore(),
// 		Analyzer: &api.Analyzer{
// 			PatternStorage: inmem.NewPatternStore(),
// 			WebsiteStorage: inmem.NewWebsiteStore(),
// 		},
// 	}
// }
