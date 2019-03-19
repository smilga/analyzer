package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Tags(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tags, err := h.TagStorage.All()
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, tags)
}

func (h *Handler) Tag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	tag, err := h.TagStorage.Get(api.TagID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, tag)
}

func (h *Handler) SaveTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tag := &api.Tag{}
	err := json.NewDecoder(r.Body).Decode(tag)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	if int64(tag.ID) == 0 {
		tag = api.NewTag(tag.Value)
	}

	err = h.TagStorage.Save(tag)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, tag)
}
