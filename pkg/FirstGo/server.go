package FirstGo

import (
	"fmt"
	"net/http"

	"github.com/zuri03/FirstGo/pkg/spotifyClient"
)

//Handlers
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func analyzeUserProfile(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("======== ANALYSIS ENDPOINT HAS BEEN CALLED ============= \n")
	values := req.URL.Query()
	if _, ok := values["error"]; ok {
		fmt.Printf("Found error in body \n")
		return
	}
	code, ok := values["code"]
	if !ok {
		url, err := spotifyClient.AuthorizationCode("http://localhost:8080/Analysis")
		if err != nil {
			//return error in response
		}
		fmt.Printf("follow this link: %s \n", url)
		http.Redirect(writer, req, url, 301)
		return
	}
	fmt.Printf("Got code : %s \n", code[0])
	fmt.Printf("Getting Client access token \n")
	client, err := spotifyClient.GetClientAccessToken("authorizationCode", code[0], spotifyClient.InitializeClient())
	if err != nil {
		//return error in response
		fmt.Printf("error getting access token: %s \n", err)
		return
	}
	fmt.Printf("Got access token: %s \n", client.Token.AccessToken)
	json, err := spotifyClient.GetUserAnalysis(client)
	if err != nil {
		fmt.Printf("error getting user data: %s \n", err)
		return
	}
	fmt.Printf("json: %s \n", json)
}

func getItemInfo(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("===== ITEM INFO ENDPOINT CALLED ================")
	values := req.URL.Query()
	searchQuery, ok := values["search"]
	if !ok {
		writer.WriteHeader(400)
		writer.Write([]byte("Missing query parameter: search"))
		return
	}
	json, err := spotifyClient.GetItemFromString(searchQuery[0], "artist", spotifyClient.InitializeClient())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Encountered server error"))
		return
	}
	writer.Write([]byte(json))
}

func InitializeServer() {
	/*
		http.HandleFunc("/Album", getAlbum)
		http.HandleFunc("/Artist", getArtist)*/
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
	url, err := spotifyClient.AuthorizationCode("http://localhost:8080/Analysis")
	if err != nil {
		return "error", err
	}
	return url, nil
}
