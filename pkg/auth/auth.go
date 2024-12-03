package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type AuthHandlers struct {
	ClientID string
	Callback string
}

const STATE_COOKIE = "oauth_state"
const VERIFIER_COOKIE = "oauth_verifier"

// PreAuth is the landing page that the user arrives at when they first attempt
// to use the agent while unauthorized.  You can do anything you want here,
// including making sure the user has an account on your side.  At some point,
// you'll probably want to make a call to the authorize endpoint to authorize
// the app.
func (h *AuthHandlers) PreAuth(w http.ResponseWriter, r *http.Request) {
	// In our example, we're not doing anything except going through the
	// authorization flow.  This is standard Oauth2.  While the flow is
	// provided verbosely here, there's probably a library in your language
	// that can do this more concisely!  In go, you can use
	// https://pkg.go.dev/golang.org/x/oauth2
	state := uuid.New()
	code := uuid.New()

	code_sha := sha256.Sum256([]byte(code.String()))
	code_challenge := base64.URLEncoding.EncodeToString(code_sha[:])

	authURL, _ := url.Parse("https://github.com/login/oauth/authorize")
	query := url.Values{}
	query.Set("state", state.String())
	query.Set("redirect_uri", h.Callback)
	query.Set("client_id", h.ClientID)
	query.Set("response_type", "code")
	query.Set("code_challenge", code_challenge)
	query.Set("code_challenge_method", "S256")
	authURL.RawQuery = query.Encode()

	// The state and code verifier should be securely saved so that they
	// can be accessed on the other side.
	stateCookie := &http.Cookie{
		Name:     STATE_COOKIE,
		Value:    state.String(),
		MaxAge:   10 * 60, // 10 minutes in seconds
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	verifierCookie := &http.Cookie{
		Name:     VERIFIER_COOKIE,
		Value:    code.String(),
		MaxAge:   10 * 60, // 10 minutes in seconds
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, stateCookie)
	http.SetCookie(w, verifierCookie)
	w.Header().Set("location", authURL.String())
	w.WriteHeader(http.StatusFound)
}

// PostAuth is the landing page where the user lads after authorizing.  As
// above, you can do anything you want here.  A common thing you might do is
// get the user information and then perform some sort of account linking in
// your database.
func (h *AuthHandlers) PostAuth(w http.ResponseWriter, r *http.Request) {
	// This is standard oauth2, just provided verbosely as an example.
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	stateCookie, err := r.Cookie(STATE_COOKIE)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("state cookie not found"))
		return
	}

	verifierCookie, err := r.Cookie(VERIFIER_COOKIE)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("code cookie not found"))
		return
	}

	// Important:  Compare the state!  This prevents CSRF attacks
	if state != stateCookie.Value {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid state"))
		return
	}

	token_body := struct {
		grant_type    string
		code          string
		redirect_uri  string
		client_id     string
		code_verifier string
	}{
		grant_type:    "authorization_code",
		code:          code,
		redirect_uri:  h.Callback,
		client_id:     h.ClientID,
		code_verifier: verifierCookie.Value,
	}

	token_body_json, err := json.Marshal(&token_body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling token request"))
		return
	}

	fmt.Println("Making token call")
	_, err = http.Post("https://github.com/login/oauth/token", "application/json", bytes.NewBuffer(token_body_json))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error error getting authorization token"))
		return
	}

	// Response contains an access token, now the world is your oyster.  Get user information and perform account linking, or do whatever you want from here.

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All done!  Please return to the app"))

}
