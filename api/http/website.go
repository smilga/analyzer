package http

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) GetWebsites(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(uid)
	ID, ok := id.(uuid.UUID)
	if !ok {
		h.responseErr(w, errors.New("Error getting context value"))
		return
	}

	websites, err := h.WebsiteStorage.ByUser(ID)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, websites)
}

func (h *Handler) CreateWebsite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(uid)
	ID, ok := id.(uuid.UUID)
	if !ok {
		h.responseErr(w, errors.New("Error getting context value"))
		return
	}
	
	website := &api.Website{}
	err := json.NewDecoder(r.Body).Decode(website)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	
	website.UserID = ID

	err = h.WebsiteStorage.Save(website)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, website)
}

func (h *Handler) ImportWebsites(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, _, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		h.responseErr(w, err)
		return
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites := []*api.Website{}
	for _, r := range records {
		if len(r) != 1 {
			// App dont know how to handle this yet
			continue
		}
		website := &api.Website{URL: r[0]}
		err := h.WebsiteStorage.Save(website)
		if err != nil {
			h.responseErr(w, err)
			return
		}
		websites = append(websites, website)
	}

	h.responseJSON(w, websites)
}
