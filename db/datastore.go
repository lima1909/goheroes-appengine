package db

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

// ListHeroes all Heroes, there are saved in datastore
func ListHeroes(ctx context.Context) ([]Hero, error) {
	heroes := []Hero{}
	ctx, err := appengine.Namespace(ctx, NAMESPACE)
	if err != nil {
		return heroes, fmt.Errorf("Err by create CTX: %v", err)
	}

	_, err = datastore.NewQuery(KIND).Order("ID").GetAll(ctx, &heroes)
	if err != nil {
		return heroes, fmt.Errorf("Err by datastore.GetAll: %v", err)
	}

	return heroes, nil
}

// AddHero add a Hero to datastore
func AddHero(ctx context.Context, h Hero) (Hero, error) {
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
