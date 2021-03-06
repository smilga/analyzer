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

	err = h.Comm.Inspect(website)
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
	websites, _, err := h.WebsiteStorage.ByUser(uid, &api.Pagination{})
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, website := range websites {
		err = h.Comm.Inspect(website)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	h.responseJSON(w, map[string]interface{}{
		"ok":    true,
		"count": len(websites),
	})
}

func (h *Handler) InspectNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites, _, err := h.WebsiteStorage.Where(uid, "inspected_at", "NULL")
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, website := range websites {
		err = h.Comm.Inspect(website)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	h.responseJSON(w, map[string]interface{}{
		"ok":    true,
		"count": len(websites),
	})
}

func (h *Handler) Inspect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	ids := []api.WebsiteID{}
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, id := range ids {
		websites, _, err := h.WebsiteStorage.Where(uid, "id", int(id))
		if err != nil {
			h.responseErr(w, err)
			return
		}
		if len(websites) == 1 {
			err = h.Comm.Inspect(websites[0])
			if err != nil {
				h.responseErr(w, err)
				return
			}
		}
	}

	h.responseJSON(w, "ok")
}
