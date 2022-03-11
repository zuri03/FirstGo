package spotify

type reponse struct {
	Items    []item `json:"items"`
	Total    int    `json:"total"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Href     string `json:"href"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type item struct {
	Album            album       `json:"album"`
	Artists          artist      `json:"artists"`
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

type artist struct {
	ExternalUrl externalUrl `json:"external_url"`
	Href        string      `json:"href"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Uri         string      `json:"uri"`
}

type image struct {
	Url    string `json:"url"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

type externalUrl struct {
	Spotify string `json:"spotify"`
}

type externalId struct {
	Isrc string `json:"isrc"`
}
