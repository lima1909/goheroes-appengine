package ps

import (
	"context"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2/google"

	"github.com/lima1909/goheroes-appengine/service"
	pubsub "google.golang.org/api/pubsub/v1"
	"google.golang.org/appengine/log"
)

const (
// ProjectID from Cloud Project
// ProjectID = "goheros-207118"
)

// HeroService wrapeer to HeroService
// protocoll HeroService calls
type HeroService struct {
	service.HeroService
	hs service.HeroService
}

// NewHeroService create a new instance
func NewHeroService(hs service.HeroService) *HeroService {
	return &HeroService{hs: hs}
}

// List protocoll list call
// PubSub in the App Engine runs only with OLD impl!!!
// good example find here: https://github.com/d2g/dg-pubsubtest
//
// find current project dynamic: appengine.RequestID(c) (ProjectID = "goheros-207118")
//
// to read the message-queue: gcloud.cmd pubsub subscriptions pull HERO_SUB
// with: --auto-ack you can clear the queue
func (hs HeroService) List(c context.Context, name string) ([]service.Hero, error) {

	hc, err := google.DefaultClient(c, pubsub.PubsubScope)
	if err != nil {
		e := fmt.Errorf("can not create new default client: %v", err)
		log.Errorf(c, "%v", e)
		return nil, e
	}

	svc, err := pubsub.New(hc)
	if err != nil {
		e := fmt.Errorf("can not create new service: %v", err)
		log.Errorf(c, "%v", e)
		return nil, e
	}

	result, err := svc.Projects.Topics.Publish(
		"projects/goheros-207118/topics/HERO",
		&pubsub.PublishRequest{
			Messages: []*pubsub.PubsubMessage{
				{
					Attributes: map[string]string{
						"ATTR1": "Yes",
						"ATTR2": "true",
					},
					Data: base64.StdEncoding.EncodeToString([]byte("hello")),
				},
			},
		},
	).Do()
	if err != nil {
		e := fmt.Errorf("Publish error: %v", err)
		log.Errorf(c, "%v", e)
		return nil, e
	}

	log.Infof(c, "Publish result: %v ", result)

	hs.hs.Add(c, "Foo")
	return hs.hs.List(c, name)
}
