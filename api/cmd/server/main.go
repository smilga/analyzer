package main

import (
	"fmt"
	"log"
	netHTTP "net/http"
	"os"

	"github.com/smilga/analyzer/api/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	g := http.NewGuard(&http.GuardConfig{
		Auth:    http.NewJWTAuth(os.Getenv("JWT_SECRET")),
		Allowed: []string{"/api/login", "api/logout"},
	})

	h := http.NewHandler()

	router.POST("/api/login", h.Login)
	router.GET("/api/logout", h.Logout)
	router.GET("/api/me", h.Me)

	router.GET("/api/services", h.GetServices)
	router.POST("/api/services", h.CreateService)
	router.GET("/api/services/:id", h.GetService)
	router.POST("/api/services/:id", h.UpdateService)
	router.GET("/api/services/:id/delete", h.DeleteService)

	router.GET("/api/websites", h.GetWebsites)
	router.POST("/api/websites", h.CreateWebsite)
	router.POST("/api/websites/import", h.ImportWebsites)

	router.GET("/api/inspect/websites/:id", h.InspectWebsite)

	port := os.Getenv("API_PORT")
	fmt.Println("Server started on port " + port)
	log.Fatal(netHTTP.ListenAndServe(":"+port, g.Protect(router)))
}
