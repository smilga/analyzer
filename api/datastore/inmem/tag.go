package inmem

import (
	"github.com/smilga/analyzer/api"
)

type TagStore struct {
	tags []*api.Tag
}

func (s *TagStore) Save(target *api.Tag) error {
	if target.ID == 0 {
		var last int64
		for _, n := range s.tags {
			if int64(n.ID) > last {
				last = int64(n.ID)
			}
		}
		target.ID = api.TagID(last + 1)
	}

	for i, t := range s.tags {
		if t.ID == api.TagID(target.ID) {
			s.tags = append(s.tags[:i], s.tags[i+1:]...)
		}
	}
	s.tags = append(s.tags, target)

	return nil
}

func (s *TagStore) All() ([]*api.Tag, error) {
	return s.tags, nil
}

func (s *TagStore) Get(id api.TagID) (*api.Tag, error) {
	for _, t := range s.tags {
		if t.ID == api.TagID(id) {
			return t, nil
		}
	}
	return nil, api.ErrTagNotFound
}

func NewTagStore() *TagStore {
	return &TagStore{
		tags: tags,
	}
}

var tags = []*api.Tag{
	&api.Tag{
		ID:    1,
		Value: "PushNotification",
	},
	&api.Tag{
		ID:    2,
		Value: "ActivePush",
	},
	&api.Tag{
		ID:    3,
		Value: "Maxtraffic",
	},
	&api.Tag{
		ID:    4,
		Value: "Google analytics",
	},
}
