package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Features(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	features, err := h.FeatureStorage.All()
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, features)
}

func (h *Handler) Feature(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	feature, err := h.FeatureStorage.Get(api.FeatureID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, feature)
}

func (h *Handler) SaveFeature(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	feature := &api.Feature{}
	err := json.NewDecoder(r.Body).Decode(feature)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	if int64(feature.ID) == 0 {
		feature = api.NewFeature(feature.Value)
	}

	err = h.FeatureStorage.Save(feature)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, feature)
}
