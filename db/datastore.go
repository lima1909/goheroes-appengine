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

// DataStoreService is a Impl from service.HeroService
type DataStoreService struct{}

// List all Heroes, there are saved in datastore
func (DataStoreService) List(c context.Context, name string) ([]service.Hero, error) {
	heroes := []service.Hero{}
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return heroes, fmt.Errorf("Err by create CTX: %v", err)
	}

	q := datastore.NewQuery(KIND)
	if name != "" {
		q = q.Filter("Name = ", name)
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

// AddHero add a Hero to datastore
func AddHero(ctx context.Context, h service.Hero) (service.Hero, error) {
	ctx, err := appengine.Namespace(ctx, NAMESPACE)
	if err != nil {
		return h, fmt.Errorf("Err by create CTX: %v", err)
	}

	k := datastore.NewIncompleteKey(ctx, KIND, nil)
	_, err = datastore.Put(ctx, k, &h)
	if err != nil {
		return h, fmt.Errorf("Err by datastore.Put: %v", err)
	}

	return h, nil
}
