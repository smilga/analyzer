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

	err = h.Analyzer.Inspect(website)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, website)
}

func (h *Handler) InspectAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites, err := h.WebsiteStorage.ByUser(uid)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for i := 0; i < len(websites); i++ {
		err = h.Analyzer.Inspect(websites[i])
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	h.responseJSON(w, "ok")
}
