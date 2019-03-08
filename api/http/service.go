package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) GetServices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(uid)
	ID, ok := id.(uuid.UUID)
	if !ok {
		h.responseErr(w, errors.New("Error getting context value"))
		return
	}

	services, err := h.ServiceStorage.ByUser(ID)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, services)
}

func (h *Handler) GetService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := uuid.Must(uuid.FromString(ps.ByName("id")))
	service, err := h.ServiceStorage.Get(ID)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	h.responseJSON(w, service)
}

func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(uid)
	ID, ok := id.(uuid.UUID)
	if !ok {
		h.responseErr(w, errors.New("Error getting context value"))
		return
	}

	service := &api.Service{}
	err := json.NewDecoder(r.Body).Decode(service)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	service.UserID = ID

	err = h.ServiceStorage.Save(service)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, "ok")
}

func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := uuid.Must(uuid.FromString(ps.ByName("id")))
	err := h.ServiceStorage.Delete(ID)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	h.responseJSON(w, "ok")
}
