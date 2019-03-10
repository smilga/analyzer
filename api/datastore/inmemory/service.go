package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type ServiceStore struct {
	services []*api.Service
}

func (s *ServiceStore) All(rel bool) ([]*api.Service, error) {
	return s.services, nil
}

func (s *ServiceStore) ByUser(id uuid.UUID) ([]*api.Service, error) {
	services := []*api.Service{}
	for _, service := range s.services {
		if service.UserID == id {
			services = append(services, service)
		}
	}
	return services, nil
}

func (s *ServiceStore) ManyByUser(id uuid.UUID, ids []uuid.UUID) ([]*api.Service, error) {
	services := []*api.Service{}
	for _, service := range s.services {
		if service.UserID == id && inSlice(service.ID, ids) {
			services = append(services, service)
		}
	}
	return services, nil
}

func (s *ServiceStore) Get(ID uuid.UUID) (*api.Service, error) {
	for _, service := range s.services {
		if service.ID == ID {
			return service, nil
		}
	}
	return nil, api.ErrServiceNotFound
}

func (s *ServiceStore) Save(target *api.Service) error {
	for i, service := range s.services {
		if service.ID == target.ID {
			s.services = append(s.services[:i], s.services[i+1:]...)
		}
	}
	s.services = append(s.services, target)

	return nil
}

func (s *ServiceStore) Delete(ID uuid.UUID) error {
	for i, service := range s.services {
		if service.ID == ID {
			s.services = append(s.services[:i], s.services[i+1:]...)
		}
	}
	return nil
}

var services = []*api.Service{
	&api.Service{
		&api.ServiceIdentity{
			ID:      uuid.Must(uuid.NewV4(), nil),
			Name:    "maxtraffic",
			LogoURL: "https://assets.mxapis.com/img/maxtraffic-logo-new.png",
		},
		uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
		[]*api.Pattern{
			&api.Pattern{
				uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
				api.Resource,
				"*mt.js*",
				false,
			},
		},
	},
	&api.Service{
		&api.ServiceIdentity{
			ID:      uuid.Must(uuid.NewV4(), nil),
			Name:    "google analytics",
			LogoURL: "https://s3.amazonaws.com/ceblog/wp-content/uploads/2018/03/24172201/why-ga-inaccurate.jpg",
		},
		uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
		[]*api.Pattern{
			&api.Pattern{
				uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
				api.Resource,
				"https://www.google-analytics.com/analytics.js*",
				false,
			},
		},
	},
	&api.Service{
		&api.ServiceIdentity{
			ID:      uuid.Must(uuid.NewV4(), nil),
			Name:    "facebook analytics",
			LogoURL: "https://scontent.frix1-1.fna.fbcdn.net/v/t39.2365-6/28985538_1634015629967706_6299852360415969280_n.png?_nc_cat=100&_nc_ht=scontent.frix1-1.fna&oh=21e83fa842cc814070fa53b97956df2d&oe=5D24F6AB",
		},
		uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
		[]*api.Pattern{
			&api.Pattern{
				uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
				api.Resource,
				"*fbevents.js*",
				false,
			},
		},
	},
	&api.Service{
		&api.ServiceIdentity{
			ID:      uuid.Must(uuid.NewV4(), nil),
			Name:    "pushcrew",
			LogoURL: "https://images.g2crowd.com/uploads/product/image/social_landscape/social_landscape_1499941730/pushcrew.png",
		},
		uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6")),
		[]*api.Pattern{},
	},
}

func NewServiceStore() *ServiceStore {
	return &ServiceStore{
		services: services,
	}
}
