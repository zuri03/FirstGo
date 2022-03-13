package server

/*
*TODO:
	- Check accept header of each incoming request
	- Set content type for each response
	- Add tests to server tests file
	- Refactor spotify client into an interface
*/
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpotifyClient interface {
	GenerateAuthorizationCodeUrl(redirectUri string, scopes ...string) string
	GetClientAccessToken(method string, redirectUri string, code string) error
	GetItemFromString(search string, itemType string) (string, error)
	GetUserAnalysis(offset int) ([]byte, error)
}

//Handlers
type analysisHandler struct{ spotify SpotifyClient }

func (h *analysisHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		writer.WriteHeader(400)
		writer.Write([]byte("Error with spotify api"))
		return
	}

	code, ok := values["code"]
	if !ok {
		url := h.spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/Analysis")
		http.Redirect(writer, req, url, 301)
		return
	}

	err := h.spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Analysis",
		code[0])
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var responses []userInfo
	for offset := 0; ; offset += 50 {
		body, err := h.spotify.GetUserAnalysis(offset)
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		var obj userInfo
		if err := json.Unmarshal(body, &obj); err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		fmt.Printf("offset, total => %d, %d \n", obj.Offset, obj.Total)
		responses = append(responses, obj)
		if obj.Offset+obj.Total >= obj.Total {
			fmt.Println("BREAKING")
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
}

type getItemHandler struct{ spotify SpotifyClient }

func (h *getItemHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok || searchQuery[0] == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("error: Missing query parameter - search"))
		return
	}

	json, err := h.spotify.GetItemFromString(searchQuery[0], "artist")
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error: Encountered server error"))
		return
	}
	writer.Write([]byte(json))
}

func InitializeServer(spotifyService SpotifyClient) {

	http.Handle("/Analysis", &analysisHandler{spotify: spotifyService})
	http.Handle("/Info", &getItemHandler{spotify: spotifyService})
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}
