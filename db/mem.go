package db

import (
	"context"
	"log"

	"github.com/lima1909/goheroes-appengine/service"
)

// MemService is a Impl from service.HeroService
type MemService struct {
	heroes []service.Hero
}

// NewMemService create a new instance of MemService
func NewMemService() *MemService {
	heroes := []service.Hero{
		service.Hero{ID: 1, Name: "Jasmin"},
		service.Hero{ID: 2, Name: "Mario"},
		service.Hero{ID: 3, Name: "Alex M"},
		service.Hero{ID: 4, Name: "Adam O"},
		service.Hero{ID: 5, Name: "Shauna C"},
		service.Hero{ID: 6, Name: "Lena H"},
		service.Hero{ID: 7, Name: "Chris S"},
	}
	return &MemService{heroes}
}

// List all Heroes, there are saved in the heroes array
func (m MemService) List(c context.Context, name string) ([]service.Hero, error) {
	return findHeroByName(m.heroes, name), nil
}

// GetByID get Hero by the ID
func (m MemService) GetByID(c context.Context, id int64) (*service.Hero, error) {
	for _, h := range m.heroes {
		if h.ID == id {
			return &h, nil
		}
	}
	return nil, service.ErrHeroNotFound
}

// Add an Hero
func (m *MemService) Add(c context.Context, h service.Hero) (*service.Hero, error) {
	m.heroes = append(m.heroes, h)
	log.Printf("add hero: %v\n", h)
	return &h, nil
}

// Update an Hero
func (m *MemService) Update(c context.Context, h service.Hero) (*service.Hero, error) {

	for i, hero := range m.heroes {
		if hero.ID == h.ID {
			m.heroes[i] = h
			log.Printf("update hero from: %v to: %v\n", hero, h)
			return &m.heroes[i], nil
		}
	}

	return nil, service.ErrHeroNotFound
}

// Delete an Hero
func (m *MemService) Delete(c context.Context, id int64) (*service.Hero, error) {
	index := -1
	for i, h := range m.heroes {
		if h.ID == id {
			index = i
			break
		}
	}

	if index != -1 {
		h := &m.heroes[index]
		log.Printf("delete hero: %v\n", h)
		m.heroes = append(m.heroes[:index], m.heroes[index+1:]...)

		return h, nil
	}

	return nil, service.ErrHeroNotFound
}

func findHeroByName(heroes []service.Hero, name string) []service.Hero {
	if name == "" {
		return heroes
	}

	hs := make([]service.Hero, 0)
	for _, h := range heroes {
		if h.Name == name {
			hs = append(hs, h)
			log.Printf("find hero: %v\n", h)
		}
	}
	return hs
}
