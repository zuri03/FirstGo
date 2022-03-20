package server

type userInfo struct {
	Items    []item `json:"items"`
	Total    int    `json:"total"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Href     string `json:"href"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}
type savedTracksResponse struct {
	Items    []savedTrackItem `json:"items"`
	Total    int              `json:"total"`
	Limit    int              `json:"limit"`
	Offset   int              `json:"offset"`
	Href     string           `json:"href"`
	Previous string           `json:"previous"`
	Next     string           `json:"next"`
}

type savedTrackItem struct {
	AddedAt string `json:"added_at"`
	Track   track  `json:"track"`
}

type item struct {
	Album            album       `json:"album"`
	Artists          []artist    `json:"artists"`
	AvailableMarkets []string    `json:"available_markets"`
	DiscNumber       int         `json:"disc_number"`
	Duration         int         `json:"duration_ms"`
	Explicit         bool        `json:"explicit"`
	ExternalIds      externalId  `json:"external_ids"`
	ExternalUrl      externalUrl `json:"external_urls"`
	Href             string      `json:"href"`
	Id               string      `json:"id"`
	IsLocal          bool        `json:"is_local"`
	Name             string      `json:"name"`
	Popularity       int         `json:"popularity"`
	PreviewUrl       string      `json:"preview_url"`
	TrackNumber      int         `json:"track_number"`
	Type             string      `json:"type"`
	Uri              string      `json:"uri"`
}

type album struct {
	AlbumType            string      `json:"album_type"`
	Artists              []artist    `json:"artists"`
	AvailableMarkets     []string    `json:"available_markets"`
	ExternalUrls         externalUrl `json:"external_urls"`
	Href                 string      `json:"href"`
	Id                   string      `json:"id"`
	ReleaseDate          string      `json:"release_date"`
	ReleaseDatePercision string      `json:"release_date_percision"`
	TotalTracks          int         `json:"total_tracks"`
	Type                 string      `json:"album"`
	Uri                  string      `json:"uri"`
}

type track struct {
	Album            album       `json:"album"`
	Artists          []artist    `json:"artists"`
	AvailableMarkets []string    `json:"available_markets"`
	DiscNumber       int         `json:"disc_number"`
	DurationMs       int         `json:"duration_ms"`
	Explicit         bool        `json:"explicit"`
	ExternalId       externalId  `json:"external_ids"`
	ExternalUrl      externalUrl `json:"external_urls"`
	Href             string      `json:"href"`
	Id               string      `json:"id"`
	IsPlayable       bool        `json:"is_playable"`
	Restriction      restriction `json:"restrictions"`
	Name             string      `json:"name"`
	Popularity       int         `json:"popularity"`
	PreviewUrl       string      `json:"preview_url"`
	TrackNumber      int         `json:"track_number"`
	Type             string      `json:"type"`
	Uri              string      `json:"uri"`
	IsLocal          bool        `json:"is_local"`
}

type restriction struct {
	Reason string `json:"reason"`
}

type artistSearchResult struct {
	Artists []artist
}

type artist struct {
	ExternalUrl externalUrl `json:"external_url"`
	Href        string      `json:"href"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Id          string      `json:"id"`
	Uri         string      `json:"uri"`

	//May have to remove this line
	Followers  follower `json:"followers"`
	Genres     []string `json:"genres"`
	Image      []image  `json:"images"`
	Popularity int      `json:"popularity"`
}

type follower struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type spotifyAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type externalUrl struct {
	Spotify string `json:"spotify"`
}

type externalId struct {
	Isrc string `json:"isrc"`
}

type recommendationsResult struct {
	Tracks []track `json:"tracks"`
	Seeds  []seed  `json:"seeds"`
}

type seed struct {
	InitialPoolSize    int    `json:"initialPoolSize"`
	AfterFilteringSize int    `json:"afterFilteringSize"`
	AfterRelinkingSize int    `json:"afterRelinkingSize"`
	Id                 string `json:"id"`
	Type               string `json:"type"`
	Href               string `json:"href"`
}
