package main

import (
	"fmt"
	"log"
	netHTTP "net/http"
	"os"
	"time"

	"github.com/smilga/analyzer/api/comm"
	"github.com/smilga/analyzer/api/datastore/mysql"
	"github.com/smilga/analyzer/api/htmlp"
	"github.com/smilga/analyzer/api/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	comm := comm.NewComm()

	g := http.NewGuard(&http.GuardConfig{
		Auth:    http.NewJWTAuth(os.Getenv("JWT_SECRET")),
		Allowed: []string{"/api/login", "api/logout", "/api/ws"},
	})

	db := mysql.NewConnection()

	h := http.NewHandler(db, comm)

	go func() {
		for {
			res, err := comm.CollectResults()
			if err != nil {
				fmt.Println("Error collecting results: ", err)
			}
			if res != nil {
				res.HTML.Value = htmlp.Parse(res.HTML.Value)
				err := h.ResultStorage.Save(res)
				if err != nil {
					fmt.Println("Error storing result: ", err)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	// go func() {
	// 	ticker := time.NewTicker(time.Millisecond * 500)
	// 	defer ticker.Stop()
	// 	done := make(chan bool)

	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case <-ticker.C:
	// 			ids := h.Messanger.UsersOnline()
	// 			for _, id := range ids {
	// 				ll, err := h.Analyzer.ListLen(api.PendingList, id)
	// 				if err != nil {
	// 					fmt.Printf("Error reporting list len: %s", err)
	// 				}
	// 				tl, err := h.Analyzer.ListLen(api.TimeoutedList, id)
	// 				if err != nil {
	// 					fmt.Printf("Error reporting list len: %s", err)
	// 				}
	// 				err = h.Messanger.SendToUser(id, &ws.Msg{
	// 					Type:   ws.CommMsg,
	// 					UserID: id,
	// 					Message: map[string]interface{}{
	// 						"action": "report:status",
	// 						"status": api.AnalyzeStatus{ll, tl},
	// 					},
	// 				})
	// 				if err != nil {
	// 					fmt.Println("Error sending report message: ", err)
	// 				}

	// 			}
	// 		}
	// 	}
	// }()

	router.POST("/api/login", h.Login)
	router.GET("/api/logout", h.Logout)
	router.GET("/api/me", h.Me)

	router.GET("/api/patterns", h.Patterns)
	router.POST("/api/patterns", h.SavePattern)
	router.GET("/api/patterns/:id", h.Pattern)
	router.GET("/api/patterns/:id/delete", h.DeletePattern)

	router.GET("/api/services", h.Services)
	router.POST("/api/services", h.SaveService)
	router.GET("/api/services/:id", h.Service)
	router.GET("/api/services/:id/delete", h.DeleteService)

	router.GET("/api/tags", h.Tags)
	router.POST("/api/tags", h.SaveTag)
	router.GET("/api/tags/:id", h.Tag)

	router.GET("/api/features", h.Features)
	router.POST("/api/features", h.SaveFeature)
	router.GET("/api/features/:id", h.Feature)

	router.GET("/api/filters", h.Filters)
	router.POST("/api/filters", h.SaveFilter)
	router.GET("/api/filters/:id", h.Filter)

	router.GET("/api/websites", h.Websites)
	router.GET("/api/websites/:id/report", h.Report)
	router.GET("/api/websites/:id/delete", h.DeleteWebsites)
	router.POST("/api/websites", h.SaveWebsite)
	router.POST("/api/websites/import", h.ImportWebsites)

	router.GET("/api/inspect/websites/:id", h.InspectWebsite)
	router.POST("/api/inspect/websites", h.Inspect)
	router.GET("/api/inspect/websites", h.InspectAll)
	router.GET("/api/inspect/new", h.InspectNew)
	router.POST("/api/inspect/export", h.Export)

	router.GET("/api/ws", h.Upgrade)

	router.PanicHandler = func(w netHTTP.ResponseWriter, r *netHTTP.Request, err interface{}) {
		fmt.Printf("Recovered from Error: %v, URL: %v %v \n", err, r.Method, r.URL)
		w.WriteHeader(netHTTP.StatusInternalServerError)
	}

	port := os.Getenv("API_PORT")
	fmt.Println("Server started on port " + port)
	log.Fatal(netHTTP.ListenAndServe(":"+port, g.Protect(router)))
}
