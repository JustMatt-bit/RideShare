package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"main/core"
	"main/db"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authModule struct {
	db     *sql.DB
	oauth2 oauth2.Config
}

var authTypeGoogle = "google"

type googleToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type googleClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type googleUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *googleUser) ToUser() core.User {
	return core.User{
		Name:  u.Name,
		Email: u.Email,
	}
}

func NewGoogleAuthModule(googleAuthConfig core.GoogleAuthConfig, db *sql.DB) *authModule {
	m := &authModule{
		db: db,
		oauth2: oauth2.Config{
			ClientID:     googleAuthConfig.ClientID,
			ClientSecret: googleAuthConfig.ClientSecret,
			RedirectURL:  googleAuthConfig.CallbackURL,
			Scopes:       []string{"openid", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}

	return m
}

func (m authModule) ApplyRoutes(r *mux.Router) {
	logMiddleware := func(handler func(http.ResponseWriter, *http.Request, *logrus.Logger)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value(core.CtxLog).(*logrus.Logger)
			handler(w, r, log)
		}
	}

	r.HandleFunc("/auth/google", logMiddleware(m.googleAuth)).Methods("GET")
	r.HandleFunc("/auth/google/callback", logMiddleware(m.googleAuthCallback)).Methods("GET")
}

func (m *authModule) googleAuth(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	state := randToken()
	u := m.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, u, http.StatusFound)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (m *authModule) googleAuthCallback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	data, err := m.getUserDataFromGoogle(log, r.FormValue("code"))
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func (m *authModule) getUserDataFromGoogle(log *logrus.Logger, code string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	v := url.Values{}
	v.Set("code", code)
	v.Set("client_id", m.oauth2.ClientID)
	v.Set("client_secret", m.oauth2.ClientSecret)
	v.Set("redirect_uri", m.oauth2.RedirectURL)
	v.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", m.oauth2.Endpoint.TokenURL, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "RideShare-Go")

	req.URL.RawQuery = v.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var token googleToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	googleUser, err := parseGoogleIDToken(log, token.IDToken)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var userID int64
	_, err = db.GetUserAuthByToken(m.db, googleUser.ID, authTypeGoogle)
	switch {
	case err == nil:
		// user exists
	case err == sql.ErrNoRows:
		// user does not exist
		userID, err = db.CreateUser(m.db, googleUser.ToUser())
		if err != nil {
			log.Error(err)
			return nil, err
		}

		userAuth := core.UserAuthRecord{
			UserID:  userID,
			Service: authTypeGoogle,
			Token:   googleUser.ID,
		}

		err = db.CreateUserAuth(m.db, userAuth)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	default:
		log.Error(err)
		return nil, err
	}

	user, err := db.GetUserByID(m.db, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	jwt, err := createToken(user, time.Now())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return []byte(jwt), nil
}

func parseGoogleIDToken(log *logrus.Logger, idToken string) (*googleUser, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		log.Error("invalid idToken")
		return nil, fmt.Errorf("invalid idToken")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var claims googleClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		log.Error(err)
		return nil, err
	}

	user := &googleUser{
		ID:    claims.Sub,
		Email: claims.Email,
		Name:  claims.Name,
	}

	return user, nil
}
