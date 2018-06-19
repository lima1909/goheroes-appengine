package db

import (
	"context"

	"github.com/lima1909/goheroes-appengine/service"
)

var (
	heroes = []service.Hero{
		service.Hero{ID: 1, Name: "Jasmin"},
		service.Hero{ID: 2, Name: "Mario"},
		service.Hero{ID: 3, Name: "Alex M"},
		service.Hero{ID: 4, Name: "Adam O"},
		service.Hero{ID: 5, Name: "Shauna C"},
		service.Hero{ID: 6, Name: "Lena H"},
		service.Hero{ID: 7, Name: "Chris S"},
	}
)

// MemService is a Impl from service.HeroService
type MemService struct{}

// List all Heroes, there are saved in the heroes array
func (MemService) List(c context.Context, name string) ([]service.Hero, error) {
	return findHeroByName(name), nil
}

// GetByID get Hero by the ID
func (MemService) GetByID(c context.Context, id int64) (service.Hero, error) {
	for _, h := range heroes {
		if h.ID == id {
			return h, nil
		}
	}
	return service.Hero{}, service.HeroNotFoundErr
}

// Add an Hero
func (MemService) Add(c context.Context, h service.Hero) (service.Hero, error) {
	heroes = append(heroes, h)
	return h, nil
}

func findHeroByName(name string) []service.Hero {
	if name == "" {
		return heroes
	}

	hs := make([]service.Hero, 0)
	for _, h := range heroes {
		if h.Name == name {
			hs = append(hs, h)
		}
	}
	return hs
}
