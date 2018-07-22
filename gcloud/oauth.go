package gcloud

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"google.golang.org/appengine"

	"github.com/gorilla/sessions"
	"github.com/lima1909/goheroes-appengine/com"
	"github.com/lima1909/goheroes-appengine/service"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	plus "google.golang.org/api/plus/v1"
	loga "google.golang.org/appengine/log"
)

// Doc:
// https://cloud.google.com/go/getting-started/authenticate-users
// https://github.com/GoogleCloudPlatform/golang-samples/tree/master/getting-started/bookshelf
// -----
const (
	// the ID, where I find the session from store
	sessionDefaultID = "Session-ID"
	// key to find the current user in the session
	sessionUserKey = "user"
	// key to find the current token in the session
	sessionOauthTokenKey = "oauthtoken"
)

var (
	oauth2Config *oauth2.Config
	sessionStore sessions.Store
	once         sync.Once
)

func init() {
	// Configure storage method for session-wide information.
	// Update "something-very-secret" with a hard to guess string or byte sequence.
	secret := uuid.Must(uuid.NewV4()).String()
	cookieStore := sessions.NewCookieStore([]byte(secret))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   1 * 60, // 1 minute
	}
	sessionStore = cookieStore

	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&service.User{})

}

// load ClientID and ClientKey from Environment or from Datastore
func getCliendIDandClientKey(ctx context.Context) (clientID, clientKey string, err error) {
	clientID = os.Getenv("CLIENT_ID")
	clientKey = os.Getenv("CLIENT_KEY")
	if clientID == "" || clientKey == "" {
		log.Println("Env-Var: CLIENT_ID and CLIENT_KEY are not set")
		log.Println("try to load from Datastore")
		clientID, clientKey, err = OAuth(ctx)
		if err != nil {
			loga.Errorf(ctx, "Can not load OAuth-Data (ClientID and ClientKey): %v", err)
			return
		}
	}
	return
}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/oauth2callback"
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint:     google.Endpoint,
	}
}

// LoginHandler initiates an OAuth flow to authenticate the user.
func LoginHandler(w http.ResponseWriter, r *http.Request) *com.Error {
	var err error
	once.Do(func() {
		var id, key string
		id, key, err = getCliendIDandClientKey(appengine.NewContext(r))
		if err != nil {
			return
		}
		oauth2Config = configureOAuthClient(id, key)
	})
	if err != nil {
		return com.Errorf(err, "could not create oauth config: %v", err)
	}

	state := uuid.Must(uuid.NewV4()).String()

	// create a new session, to save the state
	// this is importend, the callback will check this state
	sessn, err := sessionStore.New(r, state)
	if err != nil {
		return com.Errorf(err, "could not create oauth session: %v", err)
	}
	if err := sessn.Save(r, w); err != nil {
		return com.Errorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	url := oauth2Config.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)

	return nil
}

// OauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
func OauthCallbackHandler(w http.ResponseWriter, r *http.Request) *com.Error {
	// check, whether a session with the state from login (uuid) exist
	session, err := sessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return com.Errorf(err, "invalid state parameter. try logging in again.")
	}

	session, err = sessionStore.New(r, sessionDefaultID)
	if err != nil {
		return com.Errorf(err, "could not create new user session: %v", err)
	}

	ctx := appengine.NewContext(r)
	tok, err := oauth2Config.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		return com.Errorf(err, "could not get auth token: %v", err)
	}
	session.Values[sessionOauthTokenKey] = tok

	user, err := fetchUser(ctx, tok)
	if err != nil {
		return com.Errorf(err, "could not fetch Google profile: %v", err)
	}
	session.Values[sessionUserKey] = user

	if err := session.Save(r, w); err != nil {
		return com.Errorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func fetchUser(ctx context.Context, tok *oauth2.Token) (service.User, error) {
	client := oauth2.NewClient(ctx, oauth2Config.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return service.User{}, err
	}
	p, err := plusService.People.Get("me").Do()
	if err != nil {
		return service.User{}, err
	}
	return service.User{
		ID:          p.Id,
		DisplayName: p.DisplayName,
		ImageURL:    p.Image.Url,
	}, nil
}

// GetUser from the session, if exist, else nil
func GetUser(r *http.Request) (service.User, error) {
	s, err := sessionStore.Get(r, sessionDefaultID)
	if err != nil {
		return service.User{}, fmt.Errorf("err by GetUser: %v", err)
	}
	u, ok := s.Values[sessionUserKey]
	if !ok {
		return service.User{}, fmt.Errorf("no user found in session")

	}

	return *u.(*service.User), nil
}
