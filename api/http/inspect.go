package http

import (
	"encoding/json"
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

func (h *Handler) Inspect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	// TODO refacture this dont request all websites
	ids := []api.WebsiteID{}
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	idMap := make(map[api.WebsiteID]bool, len(ids))
	for _, id := range ids {
		idMap[id] = true
	}

	websites, err := h.WebsiteStorage.ByUser(uid)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, website := range websites {
		if _, ok := idMap[website.ID]; ok {
			err = h.Analyzer.Inspect(website)
			if err != nil {
				h.responseErr(w, err)
				return
			}
		}
	}

	h.responseJSON(w, "ok")
}
