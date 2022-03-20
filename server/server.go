/*
*
TODO:
	- Check accept header of each incoming request
	- Set content type for each response
	- Determine way to store and match access tokens for each request
*/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpotifyClient interface {
	GetUserTopItems(offset int, limit int, accessToken string) ([]byte, error)
	GetItemFromId(itemType string, accessToken string, ids ...string) ([]byte, error)
	GetItemFromName(search string, itemType string, accessToken string) ([]byte, error)
	GetClientAccessToken(method string, redirectUri string, code string) ([]byte, error)
	GenerateAuthorizationCodeUrl(redirectUri string, scopes ...string) string
	GetRelatedArtist(artistId string, accessToken string) ([]byte, error)
	GetTracksFromArtist(artistId string, accessToken string) ([]byte, error)
	GetSavedTracks(accessToken string) ([]byte, error)
	GetRecommendations(artistId []string, genres []string, trackIds []string, accessToken string) ([]byte, error)
}

//Handlers
type analysisHandler struct{ Spotify SpotifyClient }

func (h *analysisHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		writer.WriteHeader(400)
		writer.Write([]byte("Error with spotify api"))
		return
	}

	code, ok := values["code"]
	if !ok {
		url := h.Spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/Analysis", "user-top-read")
		http.Redirect(writer, req, url, 301)
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Analysis",
		code[0])
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var token spotifyAccessToken
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var responses []userInfo
	for offset := 0; ; offset += 50 {
		body, err := h.Spotify.GetUserTopItems(offset, 50, token.AccessToken)
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("json: \n %s \n", string(body))
		var obj userInfo
		if err := json.Unmarshal(body, &obj); err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		responses = append(responses, obj)
		if obj.Offset+obj.Total >= obj.Total {
			break
		}
	}

	genStats, err := generalStats(responses)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("Top artists => %s \n", genStats)
	writer.Write([]byte(genStats))
}

//Handlers
type playlistHandler struct{ Spotify SpotifyClient }

func (h *playlistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		writer.WriteHeader(500)
		writer.Write([]byte("Error with spotify api"))
		return
	}

	code, ok := values["code"]
	if !ok {
		url := h.Spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/MakePlaylist", "user-top-read", "playlist-modify-public", "user-library-read")
		http.Redirect(writer, req, url, 301)
		return
	} else {
		fmt.Printf("CODE FOUND \n")
	}

	bytes, err := h.Spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/MakePlaylist",
		code[0])
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("scopes => %s\n", token.Scope)
	body, err := h.Spotify.GetSavedTracks(token.AccessToken)
	var obj savedTracksResponse
	if err := json.Unmarshal(body, &obj); err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	artistId := obj.Items[0].Track.Artists[0].Id
	trackId := obj.Items[0].Track.Id
	body, err = h.Spotify.GetItemFromId("artist", token.AccessToken, artistId)
	var topArtist artist
	if err := json.Unmarshal(body, &topArtist); err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}
	genres := topArtist.Genres

	if len(genres) > 3 {
		genres = genres[:3]
	}
	body, err = h.Spotify.GetRecommendations([]string{artistId}, genres, []string{trackId}, token.AccessToken)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	var recommendations recommendationsResult
	if err := json.Unmarshal(body, &recommendations); err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	names := make([]string, 10)
	for i := 0; i < 10; i++ {
		fmt.Printf("rec %d => %s\n", i, recommendations.Tracks[i].Name)
		names = append(names, recommendations.Tracks[i].Name)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(""))
	return
}

type getItemHandler struct{ Spotify SpotifyClient }

func (h *getItemHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok || searchQuery[0] == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("error: Missing query parameter - search"))
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("clientCredentials", "", "")
	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}
	json, err := h.Spotify.GetItemFromName(searchQuery[0], "artist", token.AccessToken)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error: Encountered server error"))
		return
	}
	writer.Write([]byte(json))
}

func InitializeServer(spotifyService SpotifyClient) {

	http.Handle("/Analysis", &analysisHandler{Spotify: spotifyService})
	http.Handle("/Info", &getItemHandler{Spotify: spotifyService})
	http.Handle("/MakePlaylist", &playlistHandler{Spotify: spotifyService})
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}
