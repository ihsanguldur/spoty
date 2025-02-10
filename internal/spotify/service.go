package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"spoty/configs"
	"spoty/internal/pkg/util"
	"strings"
)

const (
	stateKey               = "spotify_auth_state"
	SpotifyAccessTokenKey  = "spotify_access_token"
	SpotifyRefreshTokenKey = "spotify_refresh_token"
)

type Service struct {
	Config *configs.Config
}

func NewSpotifyService(config *configs.Config) *Service {
	return &Service{
		Config: config,
	}
}
func (s *Service) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	state, _ := util.FixedLengthRandomString(16)

	cookie := &http.Cookie{Name: stateKey, Value: state, HttpOnly: true}
	http.SetCookie(w, cookie)

	scope := "user-read-private user-read-email user-read-recently-played user-top-read playlist-read-private"

	http.Redirect(
		w,
		r,
		fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=%s&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
			"code",
			s.Config.SpotifyClientId,
			scope,
			s.Config.SpotifyRedirectUri,
			state,
		),
		http.StatusTemporaryRedirect)
}

func (s *Service) Callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	code := query.Get("code")
	state := query.Get("state")

	// TODO check state and code here

	storedState, err := r.Cookie(stateKey)
	if err != nil {
		// TODO error handler
		if errors.Is(err, http.ErrNoCookie) {
			log.Fatal("cookie not found")
		} else {
			log.Fatal(err)
		}
	}

	if state == "" || state != storedState.Value {
		// TODO error handler
		log.Fatal("state is not match")
		w.Write([]byte("state is not match"))
		return
	}

	var spotifyAuthorizationResponse SpotifyAuthorizationResponse

	cookie := &http.Cookie{Name: stateKey, Value: "", HttpOnly: true, MaxAge: -1}
	http.SetCookie(w, cookie)

	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Config.SpotifyClientId, s.Config.SpotifyClientSecret)))
	formData := url.Values{
		"code":         {code},
		"redirect_uri": {s.Config.SpotifyRedirectUri},
		"grant_type":   {"authorization_code"},
	}

	// TODO add error handler
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(formData.Encode()))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	// TODO add error handler
	authResponse, err := client.Do(req)
	if err == nil && authResponse.StatusCode == 200 {
		//byteThing, _ := io.ReadAll(authResponse.Body)
		//fmt.Println(string(byteThing))

		_ = json.NewDecoder(authResponse.Body).Decode(&spotifyAuthorizationResponse)

		accessTokenCookie := &http.Cookie{
			Name:     SpotifyAccessTokenKey,
			Value:    spotifyAuthorizationResponse.AccessToken,
			HttpOnly: true,
			MaxAge:   spotifyAuthorizationResponse.ExpiresIn,
			Path:     "/",
		}
		http.SetCookie(w, accessTokenCookie)

		//TODO get and set cookie util
		refreshTokenCookie := &http.Cookie{
			Name:     SpotifyRefreshTokenKey,
			Value:    spotifyAuthorizationResponse.RefreshToken,
			HttpOnly: true,
			MaxAge:   math.MaxInt,
			Path:     "/",
		}
		http.SetCookie(w, refreshTokenCookie)

		http.Redirect(
			w,
			r,
			"http://localhost:8080/home",
			http.StatusPermanentRedirect)
	} else {
		// TODO error handler
		w.Write([]byte("invalid token"))
	}
	defer authResponse.Body.Close()
}

// TODO this is just test for refresh token add logic from this
func (s *Service) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie(SpotifyRefreshTokenKey)
	if err != nil {
		// TODO error handler
		if errors.Is(err, http.ErrNoCookie) {
			log.Fatal("cookie not found")
		} else {
			log.Fatal(err)
		}
	}

	if refreshToken.Value == "" {
		// TODO error handler
		log.Fatal("refresh token is not found")
		w.Write([]byte("state is not match"))
		return
	}

	var spotifyAuthorizationResponse SpotifyAuthorizationResponse

	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Config.SpotifyClientId, s.Config.SpotifyClientSecret)))
	formData := url.Values{
		"refresh_token": {refreshToken.Value},
		"grant_type":    {"refresh_token"},
	}

	// TODO add error handler
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(formData.Encode()))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	// TODO add error handler
	refreshResponse, err := client.Do(req)
	if err == nil && refreshResponse.StatusCode == 200 {
		_ = json.NewDecoder(refreshResponse.Body).Decode(&spotifyAuthorizationResponse)

		accessTokenCookie := &http.Cookie{
			Name:     SpotifyAccessTokenKey,
			Value:    spotifyAuthorizationResponse.AccessToken,
			HttpOnly: true,
			MaxAge:   spotifyAuthorizationResponse.ExpiresIn,
			Path:     "/",
		}
		http.SetCookie(w, accessTokenCookie)

		//TODO get and set cookie util
		if spotifyAuthorizationResponse.RefreshToken != "" {
			refreshTokenCookie := &http.Cookie{
				Name:     SpotifyRefreshTokenKey,
				Value:    spotifyAuthorizationResponse.RefreshToken,
				HttpOnly: true,
				MaxAge:   math.MaxInt,
				Path:     "/",
			}
			http.SetCookie(w, refreshTokenCookie)
		}

		http.Redirect(
			w,
			r,
			"http://localhost:8080/home",
			http.StatusPermanentRedirect)
	} else {
		// TODO error handler
		w.Write([]byte("invalid refresh token"))
	}

	defer refreshResponse.Body.Close()
}

// TODO: Not done yet

// --------------------------------------

func (s *Service) GetTopArtists(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie(SpotifyAccessTokenKey)
	if err != nil {
		// TODO error handler
		if errors.Is(err, http.ErrNoCookie) {
			log.Fatal("cookie not found")
		} else {
			log.Fatal(err)
		}
	}

	query := r.URL.Query()

	// bunlarin default degeri olacak api da neyse o
	timeRange := util.SetDefaultValue(query.Get("time_range"), "medium_term")
	limit := util.SetDefaultValue(query.Get("limit"), "20")

	// TODO add error handler
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artists", nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Value))

	formData := url.Values{
		"time_range": {timeRange},
		"limit":      {limit},
		"offset":     {"0"},
	}

	req.URL.RawQuery = formData.Encode()

	client := &http.Client{}
	// TODO add error handler
	response, _ := client.Do(req)

	var topArtists SpotifyTopTracks

	if err == nil && response.StatusCode == 200 {
		out, _ := io.ReadAll(response.Body)

		_ = json.NewDecoder(response.Body).Decode(&topArtists)

		w.Write(out)
	} else {
		errBody, _ := io.ReadAll(response.Body)
		fmt.Println(string(errBody))

		// TODO error handler
		w.Write([]byte("something gone wrong"))
	}

	defer response.Body.Close()
}

func (s *Service) GetTopTracks(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie(SpotifyAccessTokenKey)
	if err != nil {
		// TODO error handler
		if errors.Is(err, http.ErrNoCookie) {
			log.Fatal("cookie not found")
		} else {
			log.Fatal(err)
		}
	}

	query := r.URL.Query()

	// bunlarin default degeri olacak api da neyse o
	timeRange := util.SetDefaultValue(query.Get("time_range"), "medium_term")
	limit := util.SetDefaultValue(query.Get("limit"), "20")

	// TODO add error handler
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Value))

	formData := url.Values{
		"time_range": {timeRange},
		"limit":      {limit},
		"offset":     {"0"},
	}

	req.URL.RawQuery = formData.Encode()

	client := &http.Client{}
	// TODO add error handler
	response, _ := client.Do(req)

	var topArtists SpotifyTopTracks

	if err == nil && response.StatusCode == 200 {
		out, _ := io.ReadAll(response.Body)

		_ = json.NewDecoder(response.Body).Decode(&topArtists)

		w.Write(out)
	} else {
		errBody, _ := io.ReadAll(response.Body)
		fmt.Println(string(errBody))

		// TODO error handler
		w.Write([]byte("something gone wrong"))
	}

	defer response.Body.Close()
}

// --------------------------------------
