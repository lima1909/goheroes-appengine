package gcloud

import (
	"context"
	"fmt"

	"github.com/lima1909/goheroes-appengine/service"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	// NAMESPACE  where are the Protocol are saved
	NAMESPACE = "heroes"
	// KIND of datastore
	KIND = "Protocol"
)

// ProtocolsFromDatastore List all Protocol, there are saved in datastore
func ProtocolsFromDatastore(c context.Context) ([]service.Protocol, error) {
	c = setNamespace(c)

	q := datastore.NewQuery(KIND).Order("-Time")

	p := []service.Protocol{}
	_, err := q.GetAll(c, &p)
	if err != nil {
		return p, fmt.Errorf("Err by datastore.GetAll: %v", err)
	}

	return p, nil
}

// GetByID get Protocol by the ID
func GetByID(c context.Context, id int64) (*service.Protocol, error) {
	c = setNamespace(c)

	p := []service.Protocol{}
	ks, err := datastore.NewQuery(KIND).Filter("ID =", id).GetAll(c, &p)
	if err != nil {
		return nil, fmt.Errorf("No Protocol with ID: %v found in datastore: %v", id, err)
	}

	if len(p) > 0 && len(ks) > 0 {
		return &p[0], nil
	}

	return nil, fmt.Errorf("No Protocol found with ID: %v", id)
}

// Add a Protocol to datastore
func Add(c context.Context, p service.Protocol) error {
	c = setNamespace(c)

	k := datastore.NewIncompleteKey(c, KIND, nil)
	_, err := datastore.Put(c, k, &p)
	if err != nil {
		return fmt.Errorf("Err by datastore.Put: %v", err)
	}

	return nil
}

// Update an Protocol
// func (ds DatastoreService) Update(c context.Context, p Protocol) (*Protocol, error) {
// 	c = setNamespace(c)

// 	prot, err := ds.GetByID(c, p.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("Err by datastore.Update (GetByID: %v", err)
// 	}

// 	k := datastore.NewKey(c, KIND, "", prot.ID, nil)
// 	_, err = datastore.Put(c, k, &h)
// 	if err != nil {
// 		return nil, fmt.Errorf("Err  by datastore.Update Hero: %v with err: %v", h, err)
// 	}
// 	h.ID = k.IntID()

// 	return &h, nil
// }

// Delete a Protocol from datastore
func Delete(c context.Context, id int64) (*service.Protocol, error) {
	c = setNamespace(c)

	p, err := GetByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Delete (GetByID: %v", err)
	}

	k := datastore.NewKey(c, KIND, "", id, nil)
	err = datastore.Delete(c, k)
	if err != nil {
		return nil, fmt.Errorf("Err by datastore.Delete: %v", err)
	}

	return p, nil
}

func setNamespace(c context.Context) context.Context {
	c, err := appengine.Namespace(c, NAMESPACE)
	if err != nil {
		log.Errorf(c, fmt.Sprintf("Err by set Namespace: %v", err))
	}
	return c
}
