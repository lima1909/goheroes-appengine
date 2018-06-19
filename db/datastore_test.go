package db

import (
	"context"
	"testing"
)

func TestGetByID(t *testing.T) {
	if !testing.Short() {
		ds := DatastoreService{}
		hero, err := ds.GetByID(context.TODO(), 1)
		if err != nil {
			t.Errorf("no err expected: %v", err)
		}
		if 1 != hero.ID {
			t.Errorf("%v != %v", 1, hero.ID)
		}
	}
}
