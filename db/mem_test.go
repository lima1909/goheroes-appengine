package db

import (
	"context"
	"testing"

	"github.com/lima1909/goheroes-appengine/service"
)

func TestList(t *testing.T) {
	m := NewMemService()

	fh, _ := m.List(context.TODO(), "")

	if len(m.heroes) != len(fh) {
		t.Errorf("%v != %v", len(m.heroes), len(fh))
	}
}

func TestListFilter(t *testing.T) {
	m := NewMemService()

	fh, _ := m.List(context.TODO(), "Alex M")
	if 1 != len(fh) {
		t.Errorf("%v != %v", 1, len(fh))
	}
}

func TestListFilterNotFound(t *testing.T) {
	m := NewMemService()

	fh, _ := m.List(context.TODO(), "not available")
	if 0 != len(fh) {
		t.Errorf("%v != %v", 0, len(fh))
	}
}

func TestAdd(t *testing.T) {
	m := NewMemService()
	size := len(m.heroes)

	hero, err := m.Add(context.TODO(), "Test")
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}

	size = size + 1

	//check if list length is one item longer
	if len(m.heroes) != size {
		t.Errorf("Amount of Heroes: %v != %v", len(m.heroes), size)
	}

	//check if added hero has id = maxID
	if hero.ID != m.maxID {
		t.Errorf("ID of new Hero: %v != %v", hero.ID, m.maxID)
	}
}

func TestAddAndList(t *testing.T) {
	m := NewMemService()
	_, _ = m.Add(context.TODO(), "Alex M")

	fh, _ := m.List(context.TODO(), "Alex M")

	if 2 != len(fh) {
		t.Errorf("%v != %v", 2, len(fh))
	}
}

func TestDelete(t *testing.T) {
	m := NewMemService()
	size := len(m.heroes)

	_, _ = m.Delete(context.TODO(), 1)
	size = size - 1

	//check if size is reduced
	if len(m.heroes) != size {
		t.Errorf("%v != %v", len(m.heroes), size)
	}
}

func TestAddAndGetAndDelete(t *testing.T) {
	m := NewMemService()
	size := len(m.heroes)

	// Add
	newHero, err := m.Add(context.TODO(), "Test")
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}

	size = size + 1
	if len(m.heroes) != size {
		t.Errorf("%v != %v", len(m.heroes), size)
	}

	newID := newHero.ID

	// Get
	h, err := m.GetByID(context.TODO(), newID)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}
	if h.ID != newID {
		t.Errorf("%v != %v", newID, h.ID)
	}

	// Del
	_, _ = m.Delete(context.TODO(), newID)
	size = size - 1
	if len(m.heroes) != size {
		t.Errorf("%v != %v", len(m.heroes), size)
	}
}

func TestUpdate(t *testing.T) {
	m := NewMemService()
	h := service.Hero{ID: 7, Name: "Chris S"}

	hu, err := m.Update(context.TODO(), h)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}
	if hu.ID != h.ID {
		t.Errorf("%v != %v", hu.ID, h.ID)
	}
	if hu.Name != h.Name {
		t.Errorf("%v != %v", hu.Name, h.Name)
	}

}

func TestAddAndSearch(t *testing.T) {
	m := NewMemService()
	_, _ = m.Add(context.TODO(), "Abc")
	_, _ = m.Add(context.TODO(), "Abcd")
	_, _ = m.Add(context.TODO(), "Abcde")

	fh, _ := m.List(context.TODO(), "abc")
	if 3 != len(fh) {
		t.Errorf("%v != %v", 3, len(fh))
	}

	fh, _ = m.List(context.TODO(), "Abc")
	if 3 != len(fh) {
		t.Errorf("%v != %v", 3, len(fh))
	}

	fh, _ = m.List(context.TODO(), "cd")
	if 2 != len(fh) {
		t.Errorf("%v != %v", 2, len(fh))
	}
}

func TestSwitchPositions(t *testing.T) {
	m := NewMemService()
	h := service.Hero{ID: 7, Name: "Chris S"}

	_, err := m.UpdatePosition(context.TODO(), h, 6)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}

	heroAt6 := m.heroes[6]

	if heroAt6.ID != h.ID {
		t.Errorf("%v != %v", heroAt6.ID, h.ID)
	}

	_, err = m.UpdatePosition(context.TODO(), h, 5)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}

	heroAt5 := m.heroes[5]

	if heroAt5.ID != h.ID {
		t.Errorf("%v != %v", heroAt5.ID, h.ID)
	}

}
