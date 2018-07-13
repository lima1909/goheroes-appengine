package gcloud

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/lima1909/goheroes-appengine/service"
	"google.golang.org/appengine/urlfetch"
)

// GetBodyContent http-get from the urlString and get the bytes from the body
// the http get is depend on run in the cloud or not
func GetBodyContent(c context.Context, urlString string) (string, error) {
	var response *http.Response
	var err error

	if service.RunInCloud() {
		client := urlfetch.Client(c)
		response, err = client.Get(urlString)
	} else {
		response, err = http.Get(urlString)
	}
	if err != nil {
		return "", err
	}
	// Make HTTP GET request
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
