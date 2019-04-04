package main

import (
	"fmt"
	"log"
	netHTTP "net/http"
	"os"

	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/datastore/mysql"
	"github.com/smilga/analyzer/api/http"
	"github.com/smilga/analyzer/api/ws"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	g := http.NewGuard(&http.GuardConfig{
		Auth:    http.NewJWTAuth(os.Getenv("JWT_SECRET")),
		Allowed: []string{"/api/login", "api/logout", "/api/ws"},
	})

	db := mysql.NewConnection()

	h := http.NewHandler(db)
	go func() {
		h.Analyzer.StartReporting(func(w *api.Website) {
			err := h.Messanger.SendToUser(w.UserID, &ws.Msg{
				Type:   ws.CommMsg,
				UserID: w.UserID,
				Message: map[string]interface{}{
					"action":  "update:website",
					"website": w,
				},
			})
			if err != nil {
				fmt.Println("Error sending update website message: ", err)
			}
		})
	}()

	router.POST("/api/login", h.Login)
	router.GET("/api/logout", h.Logout)
	router.GET("/api/me", h.Me)

	router.GET("/api/patterns", h.Patterns)
	router.POST("/api/patterns", h.SavePattern)
	router.GET("/api/patterns/:id", h.Pattern)
	router.GET("/api/patterns/:id/delete", h.DeletePattern)

	router.GET("/api/tags", h.Tags)
	router.POST("/api/tags", h.SaveTag)
	router.GET("/api/tags/:id", h.Tag)

	router.GET("/api/filters", h.Filters)
	router.POST("/api/filters", h.SaveFilter)
	router.GET("/api/filters/:id", h.Filter)

	router.GET("/api/websites", h.Websites)
	router.GET("/api/websites/:id/report", h.Report)
	router.POST("/api/websites/delete", h.DeleteWebsites)
	router.POST("/api/websites", h.SaveWebsite)
	router.POST("/api/websites/import", h.ImportWebsites)

	router.GET("/api/inspect/websites/:id", h.InspectWebsite)
	router.POST("/api/inspect/websites/", h.Inspect)
	router.GET("/api/inspect/websites/", h.InspectAll)

	router.GET("/api/ws", h.Upgrade)

	router.PanicHandler = func(w netHTTP.ResponseWriter, r *netHTTP.Request, err interface{}) {
		fmt.Printf("Error: %v, URL: %v %v \n", err, r.Method, r.URL)
		w.WriteHeader(netHTTP.StatusInternalServerError)
	}

	port := os.Getenv("API_PORT")
	fmt.Println("Server started on port " + port)
	log.Fatal(netHTTP.ListenAndServe(":"+port, g.Protect(router)))
}
