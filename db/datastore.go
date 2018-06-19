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
	c = setNamespace(c)

	q := datastore.NewQuery(KIND)
	if name != "" {
		q = q.Filter("Name =", name)
		log.Infof(c, "With Filter: Name=%s", name)
	} else {
		q = q.Order("ID")
	}

	heroes := []service.Hero{}
	keys, err := q.GetAll(c, &heroes)
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
func (DatastoreService) GetByID(c context.Context, id int64) (*service.Hero, error) {
	c = setNamespace(c)

	hero := []service.Hero{}
	ks, err := datastore.NewQuery(KIND).Filter("ID =", id).GetAll(c, &hero)
	if err != nil {
		return nil, fmt.Errorf("No Hero with ID: %v found in datastore: %v", id, err)
	}

	if len(hero) > 0 && len(ks) > 0 {
		hero[0].Key = ks[0].IntID()
		return &hero[0], nil
	}

	return nil, service.ErrHeroNotFound
}

// Add a Hero to datastore
func (DatastoreService) Add(c context.Context, h service.Hero) (*service.Hero, error) {
	c = setNamespace(c)

	k := datastore.NewIncompleteKey(c, KIND, nil)
	k, err := datastore.Put(c, k, &h)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Put: %v", err)
	}

	h.Key = k.IntID()
	return &h, nil
}

// Update an Hero
func (ds DatastoreService) Update(c context.Context, h service.Hero) (*service.Hero, error) {
	c = setNamespace(c)

	hf, err := ds.GetByID(c, h.ID)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Update (GetByID: %v", err)
	}

	k := datastore.NewKey(c, KIND, "", hf.Key, nil)
	_, err = datastore.Put(c, k, &h)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Update Hero: %v with err: %v", h, err)
	}
	h.Key = k.IntID()

	return &h, nil
}

// Delete a Hero from datastore
func (ds DatastoreService) Delete(c context.Context, id int64) (*service.Hero, error) {
	c = setNamespace(c)

	hf, err := ds.GetByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Delete (GetByID: %v", err)
	}

	k := datastore.NewKey(c, KIND, "", hf.Key, nil)
	err = datastore.Delete(c, k)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Delete: %v", err)
	}

	return hf, nil
}

func setNamespace(c context.Context) context.Context {
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		log.Errorf(c, fmt.Sprintf("Err by set Namespace: %v", err))
	}
	return c
}
