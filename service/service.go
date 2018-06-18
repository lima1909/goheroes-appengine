package service

import (
	"context"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	// NAMESPACE  where are the Heroes are saved
	NAMESPACE = "heroes"
	// KIND of datastore
	KIND = "Hero"
)

type Hero struct {
	ID   int64  `json:"id" datastore:"-"`
	Name string `json:"name"`
}

type HeroService struct{}

func (HeroService) List(c context.Context) ([]Hero, error) {
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return nil, fmt.Errorf("Err by create CTX: %v", err)
	}

	heroes := []Hero{}
	_, err = datastore.NewQuery(KIND).Order("ID").GetAll(c, &heroes)
	if err != nil {
		return heroes, fmt.Errorf("Err by datastore.GetAll: %v", err)
	}

	return heroes, nil
}
