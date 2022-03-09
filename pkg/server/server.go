package server

/*
*TODO:
	- Check accept header of each incoming request
	- Set content type for each response
	- Add tests to server tests file
	- Refactor spotify client into an interface
*/
import (
	"fmt"
	"net/http"
)

type SpotifyClient interface {
	GenerateAuthorizationCodeUrl(redirectUri string) string
	GetClientAccessToken(method string, redirectUri string, code string) error
	GetItemFromString(search string, itemType string) (string, error)
	GetUserAnalysis() (string, error)
}

var (
	spotify SpotifyClient
)

//Handlers
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func SetSpotifyClient(client SpotifyClient) {
	spotify = client
}

func analyzeUserProfile(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		writer.WriteHeader(400)
		writer.Write([]byte("Error with spotify api"))
		return
	}

	code, ok := values["code"]
	if !ok {
		url := spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/Analysis")
		http.Redirect(writer, req, url, 301)
		return
	}

	err := spotify.GetClientAccessToken("authorizationCode",
		"http://localhost:8080/Analysis",
		code[0])

	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Error with spotify api"))
		return
	}

	json, err := spotify.GetUserAnalysis()
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Error with spotify api"))
		return
	}
	fmt.Printf("json: %s \n", json)
}

func getItemInfo(writer http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok || searchQuery[0] == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("error: Missing query parameter - search"))
		return
	}

	json, err := spotify.GetItemFromString(searchQuery[0], "artist")
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("error: Encountered server error"))
		return
	}
	writer.Write([]byte(json))
}

func InitializeServer() {
	http.HandleFunc("/Analysis", analyzeUserProfile)
	http.HandleFunc("/Info", getItemInfo)
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}

func SendTest() (string, error) {
	url := spotify.GenerateAuthorizationCodeUrl("http://localhost:8080/Analysis")
	return url, nil
}
