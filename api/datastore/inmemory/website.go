package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type WebsiteStore struct {
	websites []*api.Website
}

func (s *WebsiteStore) All(rel bool) ([]*api.Website, error) {
	return s.websites, nil
}

func (s *WebsiteStore) ByUser(id uuid.UUID) ([]*api.Website, error) {
	websites := []*api.Website{}
	for _, website := range s.websites {
		if website.UserID == id {
			websites = append(websites, website)
		}
	}
	return websites, nil
}

func (s *WebsiteStore) Get(ID uuid.UUID) (*api.Website, error) {
	for _, website := range s.websites {
		if website.ID == ID {
			return website, nil
		}
	}
	return nil, api.ErrWebsiteNotFound
}

func (s *WebsiteStore) Save(target *api.Website) error {
	if target.ID.String() == "00000000-0000-0000-0000-000000000000" {
		target.ID = uuid.NewV4()
	}

	for i, website := range s.websites {
		if website.ID == target.ID {
			s.websites = append(s.websites[:i], s.websites[i+1:]...)
		}
	}
	s.websites = append(s.websites, target)

	return nil
}

func (s *WebsiteStore) Delete(ID uuid.UUID) error {
	for i, website := range s.websites {
		if website.ID == ID {
			s.websites = append(s.websites[:i], s.websites[i+1:]...)
		}
	}
	return nil
}

var websites = []*api.Website{
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98")),
		URL:      "https://1a.lv",
		Services: nil,
	},
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98")),
		URL:      "https://220.lv",
		Services: nil,
	},
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98")),
		URL:      "https://nuko.lv",
		Services: nil,
	},
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
		URL:      "https://230.lv",
		Services: nil,
	},
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98")),
		URL:      "https://maxtraffic.com",
		Services: nil,
	},
	&api.Website{
		ID:       uuid.Must(uuid.NewV4(), nil),
		UserID:   uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98")),
		URL:      "https://given.lv",
		Services: nil,
	},
}

func NewWebsitesStore() *WebsiteStore {
	return &WebsiteStore{
		websites: websites,
	}
}
