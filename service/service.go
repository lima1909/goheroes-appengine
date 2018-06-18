/*
 Package service is the interface to interact with Heroes

 Golang Doc for App Engine:
	https://cloud.google.com/appengine/docs/standard/go/reference
	https://github.com/GoogleCloudPlatform/go-endpoints

 Cron:
	https://cloud.google.com/appengine/docs/standard/go/config/cron
 Taskqueue:
	https://cloud.google.com/appengine/docs/standard/go/taskqueue/
 PubSub:
 	https://cloud.google.com/pubsub/docs/overview
	https://godoc.org/cloud.google.com/go/pubsub


 Samples:
	https://github.com/GoogleCloudPlatform/golang-samples/
*/
package service

import (
	"context"
	"fmt"

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

// Hero the struct
type Hero struct {
	ID   int64  `json:"id" datastore:"-"`
	Name string `json:"name"`
}

// HeroService access to Heroes methods
type HeroService struct{}

// List get all Heros
func (HeroService) List(c context.Context) ([]Hero, error) {
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		return nil, fmt.Errorf("Err by create CTX: %v", err)
	}

	heroes := []Hero{}
	keys, err := datastore.NewQuery(KIND).Order("ID").GetAll(c, &heroes)
	if err != nil {
		log.Errorf(c, "----- Err by datastore.GetAll: %v", err)
		return heroes, fmt.Errorf("Err by datastore.GetAll: %v", err)
	}

	log.Infof(c, "----- Find %v Keys for Heroes", len(keys))
	for i, k := range keys {
		heroes[i].ID = k.IntID()
	}

	return heroes, nil
}
