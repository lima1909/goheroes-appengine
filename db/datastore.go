package db

import (
	"context"
	"fmt"

	"github.com/lima1909/goheroes-appengine/service"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	// NAMESPACE  where are the Heroes are saved
	NAMESPACE = "heroes"
	// KIND of datastore
	KIND = "Hero"
)

// DatastoreService is a Impl from service.HeroService
type DatastoreService struct{}

// List all Heroes, there are saved in datastore
func (DatastoreService) List(c context.Context, name string) ([]service.Hero, error) {
	heroes := []service.Hero{}
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return heroes, fmt.Errorf("Err by create CTX: %v", err)
	}

	q := datastore.NewQuery(KIND)
	if name != "" {
		q = q.Filter("Name =", name)
		log.Infof(c, "With Filter: Name=%s", name)
	}

	keys, err := q.Order("ID").GetAll(c, &heroes)
	if err != nil {
		return heroes, fmt.Errorf("Err by datastore.GetAll: %v", err)
	}

	log.Infof(c, "----- Find %v Keys for Heroes", len(keys))
	for i, k := range keys {
		heroes[i].Key = k.IntID()
	}

	return heroes, nil
}

// GetByID get Hero by the ID
func (DatastoreService) GetByID(c context.Context, id int64) (service.Hero, error) {
	hero := []service.Hero{}
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return service.Hero{}, fmt.Errorf("Err by create CTX: %v", err)
	}

	ks, err := datastore.NewQuery(KIND).Filter("ID =", id).GetAll(c, &hero)
	if err != nil {
		return service.Hero{}, fmt.Errorf("No Hero with ID: %v found in datastore: %v", id, err)
	}

	if len(hero) > 0 {
		hero[0].Key = ks[0].IntID()
		return hero[0], nil
	}
	return service.Hero{}, service.HeroNotFoundErr
}

// Add add a Hero to datastore
func (DatastoreService) Add(c context.Context, h service.Hero) (service.Hero, error) {
	hero := service.Hero{}
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return hero, fmt.Errorf("Err by create CTX: %v", err)
	}

	k := datastore.NewIncompleteKey(c, KIND, nil)
	k, err = datastore.Put(c, k, &h)
	if err != nil {
		return hero, fmt.Errorf("Err by datastore.Put: %v", err)
	}

	hero.Key = k.IntID()
	return hero, nil
}
