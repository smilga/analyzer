package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) InspectWebsite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	wid, err := uuid.FromString(id)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	website, err := h.WebsiteStorage.Get(api.WebsiteID(wid))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	// TODO run in goroutine and when ready post to socket channel
	// Socket could be service on handler
	err = h.Analyzer.Inspect(website)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, website)
}
