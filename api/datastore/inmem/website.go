package inmem

import (
	"github.com/smilga/analyzer/api"
)

type WebsiteStore struct {
	websites []*api.Website
}

func (s *WebsiteStore) All(rel bool) ([]*api.Website, error) {
	return s.websites, nil
}

func (s *WebsiteStore) ByUser(id api.UserID) ([]*api.Website, error) {
	websites := []*api.Website{}
	for _, website := range s.websites {
		if website.UserID == id {
			websites = append(websites, website)
		}
	}
	return websites, nil
}

func (s *WebsiteStore) Get(id api.WebsiteID) (*api.Website, error) {
	for _, w := range s.websites {
		if w.ID == api.WebsiteID(id) {
			return w, nil
		}
	}
	return nil, api.ErrWebsiteNotFound
}

func (s *WebsiteStore) Save(target *api.Website) error {
	if target.ID == 0 {
		var last int64
		for _, n := range s.websites {
			if int64(n.ID) > last {
				last = int64(n.ID)
			}
		}
		target.ID = api.WebsiteID(last + 1)
	}

	for i, website := range s.websites {
		if website.ID == target.ID {
			s.websites = append(s.websites[:i], s.websites[i+1:]...)
		}
	}
	s.websites = append(s.websites, target)

	return nil
}

func (s *WebsiteStore) Delete(id api.WebsiteID) error {
	for i, website := range s.websites {
		if website.ID == id {
			s.websites = append(s.websites[:i], s.websites[i+1:]...)
		}
	}
	return nil
}

var websites = []*api.Website{
	&api.Website{
		ID:     1,
		UserID: 2,
		URL:    "https://1a.lv",
	},
	&api.Website{
		ID:     2,
		UserID: 2,
		URL:    "https://220.lv",
	},
	&api.Website{
		ID:     3,
		UserID: 2,
		URL:    "https://nuko.lv",
	},
	&api.Website{
		ID:     4,
		UserID: 2,
		URL:    "https://230.lv",
	},
	&api.Website{
		ID:     5,
		UserID: 2,
		URL:    "https://maxtraffic.com",
	},
	&api.Website{
		ID:     6,
		UserID: 2,
		URL:    "https://given.lv",
	},
}

func NewWebsiteStore() *WebsiteStore {
	return &WebsiteStore{
		websites: websites,
	}
}
