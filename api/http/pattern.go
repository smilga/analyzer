package http

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Patterns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	patterns, err := h.PatternStorage.All()
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, patterns)
}

func (h *Handler) Pattern(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	patterns, err := h.PatternStorage.Get(api.PatternID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, patterns)
}

func (h *Handler) SavePattern(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pattern := &api.Pattern{}

	err := json.NewDecoder(r.Body).Decode(pattern)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, t := range pattern.Tags {
		if t.ID == api.TagID(uuid.UUID{}) {
			tag := api.NewTag(t.Value)
			err := h.TagStorage.Save(tag)
			if err != nil {
				h.responseErr(w, err)
				return
			}
			t = tag
		}
	}

	spew.Dump(pattern)

	err = h.PatternStorage.Save(pattern)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, pattern)
}

func (h *Handler) DeletePattern(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	err = h.PatternStorage.Delete(api.PatternID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, "ok")
}
