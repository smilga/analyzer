package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api/analyzer"
)

func (h *Handler) InspectWebsite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	services, err := h.ServiceStorage.All(true)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	website, err := h.WebsiteStorage.Get(uid)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	analyzer.Analyze(website.URL, services[0])

	return
	// TODO run ispection in go routine
}
