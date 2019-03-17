package http

import (
	"encoding/csv"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Websites(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites, err := h.WebsiteStorage.ByUser(uid)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, websites)
}

func (h *Handler) SaveWebsite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	website := &api.Website{}
	err = json.NewDecoder(r.Body).Decode(website)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	website.UserID = uid

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
