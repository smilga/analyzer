package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
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

	// TODO run in goroutine and when ready post to socket channel
	// Socket could be service on handler
	result, err := h.Analyzer.Inspect(website, services)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, result)
}
