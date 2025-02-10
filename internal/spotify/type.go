package spotify

type SpotifyAuthorizationResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// TODO: Not done yet
// -------------------------------------------------------
type SpotifyTopArtists struct {
	Href     string               `json:"href"`
	Limit    int                  `json:"limit"`
	Next     string               `json:"next"`
	Offset   int                  `json:"offset"`
	Previous string               `json:"previous"`
	Total    int                  `json:"total"`
	Items    []SpotifyArtistsItem `json:"items"`
}

type SpotifyArtistsItem struct {
	ExternalUrls ExternalUrlsItem `json:"external_urls"`
	Followers    FollowersItem    `json:"followers"`
	Genres       []string         `json:"genres"`
	Href         string           `json:"href"`
	ID           string           `json:"id"`
	Images       []ImageItem      `json:"images"`
	Name         string           `json:"name"`
	Popularity   int              `json:"popularity"`
	Type         string           `json:"type"`
	URI          string           `json:"uri"`
}

type ExternalUrlsItem struct {
	Spotify string `json:"spotify"`
}

type FollowersItem struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type ImageItem struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type SpotifyTopTracks struct {
	Href     string              `json:"href"`
	Limit    int                 `json:"limit"`
	Next     string              `json:"next"`
	Offset   int                 `json:"offset"`
	Previous string              `json:"previous"`
	Total    int                 `json:"total"`
	Items    []SpotifyTracksItem `json:"items"`
}

type SpotifyTracksItem struct {
}

// -------------------------------------------------------
