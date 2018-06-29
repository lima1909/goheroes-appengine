package gcloud

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
	ProjectID = "goheros-207118"
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
						"Name":  "Foo Bar ÜÜÜ",
					},
					Data: base64.StdEncoding.EncodeToString([]byte("My Message ÖÖÖ ßßß")),
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

	return hs.hs.List(c, name)
}

// Message mapping for PubsubMessage
type Message struct {
	Attributes map[string]string `json:"attributes,omitempty"`
	Data       string            `json:"data,omitempty"`
	MessageID  string            `json:"messageId,omitempty"`
}

// Subscription subscripe ang get all published messages
//
// https://developers.google.com/apis-explorer/#p/pubsub/v1/pubsub.projects.subscriptions.pull
//
// to read the message-queue:
// gcloud.cmd pubsub subscriptions pull HERO_SUB
// with: --auto-ack you can clear the queue
func Subscription(c context.Context) ([]Message, error) {
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

	// resultA, err := svc.Projects.Subscriptions.Acknowledge(
	// 	"projects/goheros-207118/subscriptions/HERO_SUB",
	// 	&pubsub.AcknowledgeRequest{},
	// ).Do()

	result, err := svc.Projects.Subscriptions.Pull(
		"projects/goheros-207118/subscriptions/HERO_SUB",
		&pubsub.PullRequest{
			MaxMessages:       5,
			ReturnImmediately: true,
		},
	).Do()
	if err != nil {
		e := fmt.Errorf("Publish error: %v", err)
		log.Errorf(c, "%v", e)
		return nil, e
	}

	msgs := make([]Message, len(result.ReceivedMessages))
	for i, m := range result.ReceivedMessages {
		b, _ := base64.StdEncoding.DecodeString(m.Message.Data)
		msgs[i] = Message{
			Attributes: m.Message.Attributes,
			Data:       string(b),
			MessageID:  m.Message.MessageId,
		}
	}

	return msgs, nil
}
