package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
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
	id, err := uuid.FromString(idStr)
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

	if tag.ID == api.TagID(uuid.UUID{}) {
		tag = api.NewTag(tag.Value)
	}

	err = h.TagStorage.Save(tag)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, tag)
}
