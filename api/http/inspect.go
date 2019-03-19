package http

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) InspectWebsite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	website, err := h.WebsiteStorage.Get(api.WebsiteID(id))
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
