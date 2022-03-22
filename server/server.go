package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zuri03/FirstGo/spotify"
)

//Handlers
type analysisHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *analysisHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("GOT ANALYSIS REQUEST")
	var resp struct {
		data  []string
		error string
	}

	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		resp.error = "Server error has occured"
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
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
		resp.error = "Server error has occured"
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	var token spotifyAccessToken
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		resp.error = "Server error has occured"
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	var responses []userInfo
	for offset := 0; ; offset += 50 {
		body, err := h.Spotify.GetUserTopItems(offset, 50, token.AccessToken, expires)
		if err != nil {
			resp.error = "Server error has occured"
			writer.WriteHeader(500)
			json, _ := json.Marshal(resp)
			writer.Write(json)
			return
		}
		fmt.Printf("json: \n %s \n", string(body))
		var obj userInfo
		if err := json.Unmarshal(body, &obj); err != nil {
			resp.error = "Server error has occured"
			writer.WriteHeader(500)
			json, _ := json.Marshal(resp)
			writer.Write(json)
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
type playlistHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *playlistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("GOT PLAYLIST REQUEST")
	var resp struct {
		Data  [10]string
		Error string
	}

	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		resp.Error = "Server error has occured"
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	code, ok := values["code"]
	if !ok {
		url := h.Spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/Playlist", "user-top-read", "playlist-modify-public", "user-library-read")
		http.Redirect(writer, req, url, 301)
		return
	} else {
		fmt.Printf("CODE FOUND \n")
	}

	bytes, err := h.Spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Playlist",
		code[0])
	if err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	//Parse access token
	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}
	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	//Get users saved tracks to use as parameters for recommendations
	body, err := h.Spotify.GetSavedTracks(token.AccessToken, expires)
	var obj savedTracksResponse
	if err := json.Unmarshal(body, &obj); err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	//Gather query params for song recomendations
	artistId := obj.Items[0].Track.Artists[0].Id
	trackId := obj.Items[0].Track.Id
	body, err = h.Spotify.GetItemFromId("artist", token.AccessToken, expires, artistId)
	var topArtist artist
	if err := json.Unmarshal(body, &topArtist); err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}
	genres := topArtist.Genres

	//Get song recommendations from spotify client
	body, err = h.Spotify.GetRecommendations([]string{artistId}, genres, []string{trackId}, token.AccessToken, expires)
	if err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	//Parse a recommendation result object checking for errors
	var recommendations recommendationsResult
	if err := json.Unmarshal(body, &recommendations); err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	//For now make an array of the first 10 recommendations' names
	var names [10]string
	for i := 0; i < 10; i++ {
		names[i] = recommendations.Tracks[i].Name
	}

	//Create the response
	resp.Data = names
	resp.Error = ""
	json, _ := json.Marshal(resp)
	fmt.Printf("json => %s\n", json)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
	return
}

type getItemHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *getItemHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("got req => %s\n", string(req.URL.Path))
	var resp struct {
		Data  artistSearchResult
		Error string
	}
	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok || searchQuery[0] == "" {
		resp.Error = "error: Missing query parameter - search"
		writer.WriteHeader(400)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("clientCredentials", "", "")
	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}
	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	item, err := h.Spotify.GetItemFromName(searchQuery[0], "artist", token.AccessToken, expires)
	if err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	var obj artistResponseWrapper
	if err := json.Unmarshal(item, &obj); err != nil {
		resp.Error = err.Error()
		writer.WriteHeader(500)
		json, _ := json.Marshal(resp)
		writer.Write(json)
		return
	}

	resp.Data = obj.Artist
	resp.Error = ""
	json, _ := json.Marshal(resp)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
	return
}

func InitializeServer() {
	s := spotify.NewClient()
	http.Handle("/Analysis", &analysisHandler{Spotify: s})
	http.Handle("/Info", &getItemHandler{Spotify: s})
	http.Handle("/Playlist", &playlistHandler{Spotify: s})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
	fmt.Printf("now listening on port 8080")
}
