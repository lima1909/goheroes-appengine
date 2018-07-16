package score

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lima1909/goheroes-appengine/db"
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

		pageContent, err := getBodyContent(urlString, Default().client(context.TODO()))
		if err != nil {
			t.Errorf("no err expected: %v", err)
		}
		score, err := getScore(pageContent, name)
		if err != nil {
			t.Errorf("no err expected: %v", err)
		}
		if score == 0 {
			t.Errorf("expected a string, got %v ", score)
		}
	}
}

func TestGetScoreWrongURL(t *testing.T) {
	_, err := getBodyContent("https://wrongURL.com", Default().client(context.TODO()))
	if err == nil {
		t.Errorf("Should throw an Error because of wrong URL")
	}
}
func TestGetScoreWrongName(t *testing.T) {
	if runTestsAgainst8anu {
		urlString := "https://www.8a.nu/de/scorecard/ranking/?City=Nuremberg"
		name := "jasmin-test"

		pageContent, err := getBodyContent(urlString, Default().client(context.TODO()))
		if err != nil {
			t.Errorf("no err expected: %v", err)
		}
		_, err = getScore(pageContent, name)
		if err == nil {
			t.Errorf("Should throw an Error because of wrong name")
		}
	}
}
func TestCreateScoreMap(t *testing.T) {
	// would be nice to mock the return of getScore, so that I don't have to call 8a.nu!!

	if runTestsAgainst8anu {

		scores, _ := Default().Scores(context.TODO(), db.NewMemService())

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
