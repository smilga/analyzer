package http

import (
	"encoding/csv"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Export(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	filterIDs, err := parseFilterString(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites := []*api.Website{}

	if len(filterIDs) == 0 {
		websites, _, err = h.WebsiteStorage.ByUser(uid, &api.Pagination{})
		if err != nil {
			h.responseErr(w, err)
			return
		}
	} else {
		websites, _, err = h.WebsiteStorage.ByFilterID(filterIDs, api.UserID(uid), &api.Pagination{})
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	// h.responseJSON(w, websites)
	csvData := api.WebsitesToCsv(websites)
	csvWriter := csv.NewWriter(w)
	csvWriter.WriteAll(csvData)
}
