package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lima1909/goheroes-appengine/service"
)

/** only run real tests agains 8a.nu from time to time and not automatically
- you'll need a network connections
- too much calls to 8a.nu can get us into trouble */
var runTestsAgainst8anu bool

func init() {
	// to run tests against 8a.nu, you must set a Environment Variable: NU to any value
	runTestsAgainst8anu = len(os.Getenv("NU")) > 0
	if runTestsAgainst8anu {
		fmt.Printf("\n-----> runTestsAgainst8anu is enabled: %v <-----\n\n", runTestsAgainst8anu)
	}
}

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
	m.Add(context.TODO(), "Alex M")

	fh, _ := m.List(context.TODO(), "Alex M")

	if 2 != len(fh) {
		t.Errorf("%v != %v", 2, len(fh))
	}
}

func TestDelete(t *testing.T) {
	m := NewMemService()
	size := len(m.heroes)

	m.Delete(context.TODO(), 1)
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
	m.Delete(context.TODO(), newID)
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
	m.Add(context.TODO(), "Abc")
	m.Add(context.TODO(), "Abcd")
	m.Add(context.TODO(), "Abcde")

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

func TestConvertToNumber(t *testing.T) {
	s1 := "1 234"
	s2 := "1A234"
	s3 := "1A 234"
	s4 := "1A 2A 3A4"
	notANb := "ABC"

	res := convertToNumber(s1)
	if res != 1234 {
		t.Errorf("%v can not be converted to number; Result %v; Expected %v ", s1, res, 1234)
	}
	res = convertToNumber(s2)
	if res != 1234 {
		t.Errorf("%v can not be converted to number; Result %v; Expected %v ", s2, res, 1234)
	}
	res = convertToNumber(s3)
	if res != 1234 {
		t.Errorf("%v can not be converted to number; Result %v; Expected %v ", s3, res, 1234)
	}
	res = convertToNumber(s4)
	if res != 1234 {
		t.Errorf("%v can not be converted to number; Result %v; Expected %v ", s4, res, 1234)
	}
	res = convertToNumber(notANb)
	if res != 0 {
		t.Errorf("%v can not be converted to number; Result %v; Expected %v ", notANb, res, 0)
	}
}

func TestGetScore(t *testing.T) {
	if runTestsAgainst8anu {

		urlString := "https://www.8a.nu/de/scorecard/ranking/?City=Nuremberg"
		name := "jasmin-roeper"

		score, err := getScore(urlString, name)
		if err != nil {
			t.Errorf("no err expected: %v", err)
		}
		if score == 0 {
			t.Errorf("expected a string, got %v ", score)
		}
	}
}

func TestGetScoreWrongURL(t *testing.T) {
	urlString := "https://wrongURL.com"
	name := "jasmin-roeper"

	_, err := getScore(urlString, name)
	if err == nil {
		t.Errorf("Should throw an Error because of wrong URL")
	}
}
func TestGetScoreWrongName(t *testing.T) {
	if runTestsAgainst8anu {
		urlString := "https://www.8a.nu/de/scorecard/ranking/?City=Nuremberg"
		name := "jasmin-test"

		_, err := getScore(urlString, name)
		if err == nil {
			t.Errorf("Should throw an Error because of wrong name")
		}
	}
}
func TestCreateScoreMap(t *testing.T) {
	// would be nice to mock the return of getScore, so that I don't have to call 8a.nu!!

	m := NewMemService()

	if runTestsAgainst8anu {
		scores := m.CreateScoreMap(context.TODO())

		if scores[1] <= 0 {
			t.Errorf("Expected a score of over 0 for 1. Hero; got %v ", scores[1])
		}
		if scores[2] <= 0 {
			t.Errorf("Expected a score of over 0 for 2. Hero; got %v ", scores[2])
		}
		if scores[4] != 0 {
			t.Errorf("Expected a score of 0 for 4. Hero; got %v ", scores[4])
		}
		if scores[5] != 0 {
			t.Errorf("Expected a score of 0 for 5. Hero; got %v ", scores[5])
		}
	}
}
