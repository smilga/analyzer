package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Services(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	services := []*api.Service{}
	var total int

	pagination, err := parsePagination(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	featureIDs, err := parseFeatureString(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	if len(featureIDs) == 0 {
		services, total, err = h.ServiceStorage.All(pagination)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	} else {
		services, total, err = h.ServiceStorage.ByFeatures(featureIDs, pagination)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	h.responseJSON(w, map[string]interface{}{
		"services": services,
		"total":    total,
	})
}

func (h *Handler) Service(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	services, err := h.ServiceStorage.Get(api.ServiceID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, services)
}

func (h *Handler) SaveService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service := &api.Service{}

	err := json.NewDecoder(r.Body).Decode(service)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, t := range service.Features {
		if int64(t.ID) == 0 {
			err := h.FeatureStorage.Save(t)
			if err != nil {
				h.responseErr(w, err)
				return
			}
		}
	}

	err = h.ServiceStorage.Save(service)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, service)
}

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	err = h.ServiceStorage.Delete(api.ServiceID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, "ok")
}

func parseFeatureString(r *http.Request) ([]api.FeatureID, error) {
	featureIDs := []api.FeatureID{}
	if featureStr, ok := r.URL.Query()["f"]; ok {
		fStr := featureStr[0]
		fs := strings.Split(fStr, ",")

		for _, sid := range fs {
			if len(sid) == 0 {
				continue
			}
			id, err := strconv.Atoi(sid)
			if err != nil {
				return nil, err
			}
			featureIDs = append(featureIDs, api.FeatureID(id))
		}
	}
	return featureIDs, nil
}
