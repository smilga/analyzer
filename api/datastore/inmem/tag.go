package inmem

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type TagStore struct {
	tags []*api.Tag
}

func (s *TagStore) Save(target *api.Tag) error {
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
		ID:    api.TagID(uuid.Must(uuid.FromString("db36d210-9ca4-47f7-ab99-94e817e19daa"))),
		Value: "PushNotification",
	},
	&api.Tag{
		ID:    api.TagID(uuid.Must(uuid.FromString("43115aa3-3c9a-4cdd-b0c0-c74c9aeeaed6"))),
		Value: "ActivePush",
	},
	&api.Tag{
		ID:    api.TagID(uuid.Must(uuid.FromString("f19a8845-68ff-4165-a2b4-fc6d57ab0d8c"))),
		Value: "Maxtraffic",
	},
	&api.Tag{
		ID:    api.TagID(uuid.Must(uuid.FromString("c7d2b8b5-33e8-4593-849a-43c23ddf78f3"))),
		Value: "Google analytics",
	},
}
