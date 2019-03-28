package http

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/smilga/analyzer/api"
)

func (h *Handler) Websites(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	websites := []*api.Website{}

	filterIDs, err := parseFilterString(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	if len(filterIDs) == 0 {
		websites, err = h.WebsiteStorage.ByUser(uid)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	} else {
		websites, err = h.WebsiteStorage.ByFilterID(filterIDs, api.UserID(uid))
		if err != nil {
			h.responseErr(w, err)
			return
		}
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
	uid, err := h.AuthID(r)
	if err != nil {
		h.responseErr(w, err)
		return
	}

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
		website.UserID = uid
		err := h.WebsiteStorage.Save(website)
		if err != nil {
			h.responseErr(w, err)
			return
		}
		websites = append(websites, website)
	}

	h.responseJSON(w, websites)
}

func (h *Handler) DeleteWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ids := []api.WebsiteID{}
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for _, id := range ids {
		err = h.WebsiteStorage.Delete(id)
		if err != nil {
			h.responseErr(w, err)
			return
		}
	}

	h.responseJSON(w, "ok")
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO check if website belongs to user
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		h.responseErr(w, err)
		return
	}
	report, err := h.ReportStorage.ByWebsite(api.WebsiteID(id))
	if err != nil {
		h.responseErr(w, err)
		return
	}

	h.responseJSON(w, report)

}

func parseFilterString(r *http.Request) ([]api.FilterID, error) {
	filterIDs := []api.FilterID{}
	if filterStr, ok := r.URL.Query()["f"]; ok {
		fStr := filterStr[0]
		fs := strings.Split(fStr, ",")

		for _, sid := range fs {
			if len(sid) == 0 {
				continue
			}
			id, err := strconv.Atoi(sid)
			if err != nil {
				return nil, err
			}
			filterIDs = append(filterIDs, api.FilterID(id))
		}
	}
	return filterIDs, nil
}
