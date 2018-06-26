// Package service is the interface to interact with Heroes
package service

import (
	"context"
	"errors"
)

var (
	// ErrHeroNotFound if no Hero was found
	ErrHeroNotFound = errors.New("Hero not Found")
	// ErrPosNotFound if new Position is out of range
	ErrPosNotFound = errors.New("Out of Range")
)

// Hero the struct
type Hero struct {
	Key  int64  `json:"key" datastore:"-"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// HeroService access to Heroes methods
type HeroService interface {
	List(c context.Context, name string) ([]Hero, error)
	GetByID(c context.Context, id int64) (*Hero, error)
	Add(c context.Context, n string) (*Hero, error)
	Update(c context.Context, h Hero) (*Hero, error)
	UpdatePosition(c context.Context, h Hero, pos int64) (*Hero, error)
	Delete(c context.Context, id int64) (*Hero, error)
}
