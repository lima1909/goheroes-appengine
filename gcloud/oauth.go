package gcloud

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	plus "google.golang.org/api/plus/v1"
)

// Doc:
// https://cloud.google.com/go/getting-started/authenticate-users
// https://github.com/GoogleCloudPlatform/golang-samples/tree/master/getting-started/bookshelf
// -----
const (
	// This key is used in the OAuth flow session to store the URL to redirect the
	// user to after the OAuth flow is complete.
	oauthFlowRedirectKey = "redirect"
	defaultSessionID     = "default"
	// The following keys are used for the default session. For example:
	//  session, _ := bookshelf.SessionStore.New(r, defaultSessionID)
	//  session.Values[oauthTokenSessionKey]
	googleProfileSessionKey = "google_profile"
	oauthTokenSessionKey    = "oauthtoken"
)

var (
	oauth2Config *oauth2.Config
	sessionStore sessions.Store
)

func init() {
	// Configure storage method for session-wide information.
	// Update "something-very-secret" with a hard to guess string or byte sequence.
	cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	sessionStore = cookieStore

	oauth2Config = configureOAuthClient(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_KEY"))

	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&Profile{})

}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/oauth2callback", //"https://goheros-207118.appspot.com/oauth2callback",
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint:     google.Endpoint,
	}
}

// LoginHandler initiates an OAuth flow to authenticate the user.
func LoginHandler(w http.ResponseWriter, r *http.Request) *AppError {
	sessionID := uuid.Must(uuid.NewV4()).String()

	oauthFlowSession, err := sessionStore.New(r, sessionID)
	if err != nil {
		return AppErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		return AppErrorf(err, "invalid redirect URL: %v", err)
	}
	oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURL

	if err := oauthFlowSession.Save(r, w); err != nil {
		return AppErrorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	url := oauth2Config.AuthCodeURL(sessionID, oauth2.ApprovalForce,
		oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)

	return nil
}

// OauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
func OauthCallbackHandler(w http.ResponseWriter, r *http.Request) *AppError {
	oauthFlowSession, err := sessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return AppErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		return AppErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return AppErrorf(err, "could not get auth token: %v", err)
	}

	session, err := sessionStore.New(r, defaultSessionID)
	if err != nil {
		return AppErrorf(err, "could not get default session: %v", err)
	}

	ctx := context.Background()
	profile, err := fetchProfile(ctx, tok)
	if err != nil {
		return AppErrorf(err, "could not fetch Google profile: %v", err)
	}

	session.Values[oauthTokenSessionKey] = tok
	// Strip the profile to only the fields we need. Otherwise the struct is too big.
	session.Values[googleProfileSessionKey] = stripProfile(profile)
	if err := session.Save(r, w); err != nil {
		return AppErrorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}

// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func fetchProfile(ctx context.Context, tok *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, oauth2Config.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

// Profile ...
type Profile struct {
	ID, DisplayName, ImageURL string
}

// stripProfile returns a subset of a plus.Person.
func stripProfile(p *plus.Person) *Profile {
	return &Profile{
		ID:          p.Id,
		DisplayName: p.DisplayName,
		ImageURL:    p.Image.Url,
	}
}

// validateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must not be absolute")
	}
	return path, nil
}

// AppHandler http://blog.golang.org/error-handling-and-go
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

// AppError ...
type AppError struct {
	Error   error
	Message string
	Code    int
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

// AppErrorf ...
func AppErrorf(err error, format string, v ...interface{}) *AppError {
	return &AppError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
