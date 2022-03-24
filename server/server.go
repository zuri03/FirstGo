package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zuri03/FirstGo/spotify"
)

var (
	logger *log.Logger
)

func returnErr(writer http.ResponseWriter, statusCode int, res interface{}, errMsg string, err error) {
	if err != nil {
		logger.Printf("ERROR: %s : %s\n", errMsg, err.Error())
	} else {
		logger.Printf("ERROR: %s : %s\n", errMsg, "")
	}
	writer.WriteHeader(500)
	json, _ := json.Marshal(res)
	writer.Write(json)
}

//Handlers
type analysisHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *analysisHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var resp struct {
		Data  []string
		Error string
	}

	writer.Header().Set("Content-Type", "application/json")

	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Spotify authorization endpoint returned error", nil)
		return
	}

	code, ok := values["code"]
	if !ok {
		resp.Error = "Not Authorized"
		returnErr(writer, 400, resp, "Missing authoriztion code", nil)
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Analysis",
		code[0])
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "GetClientAccessToken returned error", nil)
		return
	}

	var token spotifyAccessToken
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling spotify access token", nil)
		return
	}

	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	var responses []userInfo
	for offset := 0; ; offset += 50 {
		body, err := h.Spotify.GetUserTopItems(offset, 50, token.AccessToken, expires)
		if err != nil {
			resp.Error = "Internal Server Error"
			returnErr(writer, 500, resp, "GetUserTopItems returned error", nil)
			return
		}

		var obj userInfo
		if err := json.Unmarshal(body, &obj); err != nil {
			resp.Error = "Internal Server Error"
			returnErr(writer, 500, resp, "Error unmarshaling json", nil)
			return
		}

		responses = append(responses, obj)
		if obj.Offset+obj.Total >= obj.Total {
			break
		}
	}

	genStats, err := generalStats(responses)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "generalStats returned error", nil)
		return
	}

	res, _ := json.Marshal(genStats)
	writer.Write([]byte(res))
	return
}

//Handlers
type playlistHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *playlistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var resp struct {
		Data  [10]string
		Error string
	}

	writer.Header().Set("Content-Type", "application/json")

	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		resp.Error = "Server error has occured"
		returnErr(writer, 500, resp, "Spotify authorization api returned error", nil)
		return
	}

	code, ok := values["code"]
	if !ok {
		resp.Error = "Not Authorized"
		returnErr(writer, 400, resp, "Request was missing authorization code", nil)
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Playlist",
		code[0])
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "GetClientAccessToken returned error", err)
		return
	}

	//Parse access token
	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling access token", err)
		return
	}
	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	//Get users saved tracks to use as parameters for recommendations
	body, err := h.Spotify.GetSavedTracks(token.AccessToken, expires)
	var obj savedTracksResponse
	if err := json.Unmarshal(body, &obj); err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling json", err)
		return
	}

	//Gather query params for song recomendations
	artistId := obj.Items[0].Track.Artists[0].Id
	trackId := obj.Items[0].Track.Id
	body, err = h.Spotify.GetItemFromId("artist", token.AccessToken, expires, artistId)
	var topArtist artist
	if err := json.Unmarshal(body, &topArtist); err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling json", err)
		return
	}
	genres := topArtist.Genres

	//Get song recommendations from spotify client
	body, err = h.Spotify.GetRecommendations([]string{artistId}, genres, []string{trackId}, token.AccessToken, expires)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "GetRecommendations returned error", err)
		return
	}

	//Parse a recommendation result object checking for errors
	var recommendations recommendationsResult
	if err := json.Unmarshal(body, &recommendations); err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling json", err)
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
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}

type getItemHandler struct{ Spotify *spotify.SpotifyApiClient }

func (h *getItemHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var resp struct {
		Data  artistSearchResult
		Error string
	}

	writer.Header().Set("Content-Type", "application/json")

	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok || searchQuery[0] == "" {
		resp.Error = "Missing parameter: search"
		returnErr(writer, 400, resp, "User missing search query param", nil)
		return
	}

	bytes, err := h.Spotify.GetClientAccessToken("clientCredentials", "", "")
	var token spotifyAccessToken
	err = json.Unmarshal([]byte(bytes), &token)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling access token", err)
		return
	}
	expires := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	item, err := h.Spotify.GetItemFromName(searchQuery[0], "artist", token.AccessToken, expires)
	if err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "GetItemFromName returned error", err)
		return
	}

	var obj artistResponseWrapper
	if err := json.Unmarshal(item, &obj); err != nil {
		resp.Error = "Internal Server Error"
		returnErr(writer, 500, resp, "Error unmarshaling artist response json", err)
		return
	}

	resp.Data = obj.Artist
	resp.Error = ""
	json, _ := json.Marshal(resp)

	writer.Write(json)
	return
}

func InitializeServer(f *os.File) {

	logger = log.New(f, "http: ", log.LstdFlags)

	s := spotify.NewClient()
	http.Handle("/Analysis", &analysisHandler{Spotify: s})
	http.Handle("/Info", &getItemHandler{Spotify: s})
	http.Handle("/Playlist", &playlistHandler{Spotify: s})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			logger.Fatal(err)
			return
		}
	}()
	logger.Printf("now listening on port 8080 \n")
}
