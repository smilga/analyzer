package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	id, err := strconv.Atoi(idStr)
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
		if int64(t.ID) == 0 {
			err := h.TagStorage.Save(t)
			if err != nil {
				h.responseErr(w, err)
				return
			}
		}
	}

	err = h.PatternStorage.Save(pattern)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, pattern)
}

func (h *Handler) DeletePattern(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
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
