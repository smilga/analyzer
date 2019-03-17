package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Filters(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filters, err := h.FilterStorage.All()
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, filters)
}

func (h *Handler) Filter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	filter, err := h.FilterStorage.Get(api.FilterID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, filter)
}

func (h *Handler) SaveFilter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filter := &api.Filter{}
	err := json.NewDecoder(r.Body).Decode(filter)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	err = h.FilterStorage.Save(filter)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, filter)
}
