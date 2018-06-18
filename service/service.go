// Package service is the interface to interact with Heroes
package service

import (
	"context"
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
}
