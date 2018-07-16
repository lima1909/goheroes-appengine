package score

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/lima1909/goheroes-appengine/service"
)

// CreateClientFunc create a http.Client (in the cloud urlfetch.Client)
type CreateClientFunc func(c context.Context) *http.Client

// defaultClientFunc is the default, if not run in the cloud
var defaultClientFunc = func(c context.Context) *http.Client {
	return http.DefaultClient
}

// Score ...
type Score struct {
	client CreateClientFunc
}

// New instace of Score
func New(client CreateClientFunc) Score {
	return Score{client: client}
}

// Default instance of Score
func Default() Score {
	return Score{client: defaultClientFunc}
}

// Scores impl from ScoreService, get Scores by HeroService
func (s Score) Scores(c context.Context, svc service.HeroService) (map[int64]int, error) {
	heroes, err := svc.List(c, "")
	if err != nil {
		return nil, err
	}
	return s.ScoresByList(c, heroes)
}

// ScoresByList get Scores by Hero-List
func (s Score) ScoresByList(c context.Context, heroes []service.Hero) (map[int64]int, error) {
	type hscore struct {
		id    int64
		score int
		err   error
	}

	scores := make(chan hscore, 5)

	wg := sync.WaitGroup{}
	wg.Add(len(heroes))

	go func() {
		wg.Wait()
		close(scores)
	}()

	for _, h := range heroes {
		go func(c context.Context, h service.Hero) {
			defer wg.Done()

			hs := hscore{id: h.ID, score: 0}
			if h.ScoreData.Name != "" {
				hs.score, hs.err = s.Get(c, h)
			}
			scores <- hs
		}(c, h)
	}

	// convert channel to map
	scoreMap := map[int64]int{}
	for hs := range scores {
		if hs.err != nil {
			return nil, hs.err
		}
		scoreMap[hs.id] = hs.score
	}

	return scoreMap, nil
}

// Get the Score from a Hero
func (s Score) Get(c context.Context, h service.Hero) (int, error) {
	url := fmt.Sprintf("https://www.8a.nu/%s/scorecard/ranking/?City=%s", h.ScoreData.Country, h.ScoreData.City)
	pageContent, err := getBodyContent(url, s.client(c))
	if err != nil {
		return 0, err
	}

	score, err := getScore(pageContent, h.ScoreData.Name)
	if err != nil {
		return 0, err
	}

	return score, nil
}

func getBodyContent(url string, client *http.Client) (string, error) {
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("err by GET with URL: %s %v", url, err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	pageContent := string(dataInBytes)
	if pageContent == "" {
		return "", service.ErrNoContent
	}

	return pageContent, nil
}

func getScore(pageContent, name string) (int, error) {

	// Find a substr
	startIndex := strings.Index(pageContent, name)
	if startIndex == -1 {
		return 0, fmt.Errorf("Can not find %v ", name)
	}

	subString := pageContent[startIndex:(startIndex + 200)]

	// Find score
	indexStart := strings.Index(subString, "\">")
	indexEnd := strings.Index(subString, "</a>")

	if indexStart == -1 || indexEnd == -1 {
		return 0, fmt.Errorf("Can not find score for %v " + name)
	}

	return convertToNumber(subString[(indexStart + 2):indexEnd]), nil
}

func convertToNumber(s string) int {
	re := regexp.MustCompile("[0-9]+")
	scoreNumberArray := re.FindAllString(s, -1)

	scoreNumberString := ""
	for _, c := range scoreNumberArray {
		scoreNumberString = scoreNumberString + c
	}

	nb, err := strconv.Atoi(scoreNumberString)
	if err != nil {
		return 0
	}

	return nb
}
