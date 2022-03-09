package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	FirstGo "github.com/zuri03/FirstGo/pkg/server"
	Spotify "github.com/zuri03/FirstGo/pkg/spotifyClient"
)

func main() {

	cleanVariables := func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}

	data, err := os.ReadFile("../.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	credentials := strings.Split(string(data), "\n")
	clientId := strings.Split(credentials[1], "=")
	clientSecret := strings.Split(credentials[0], "=")

	secret := strings.Map(cleanVariables, clientSecret[1])
	id := strings.Map(cleanVariables, clientId[1])
	os.Setenv(clientId[0], id)
	os.Setenv(clientSecret[0], secret)

	fmt.Println("====== INITIALIZING SERVER =====")
	FirstGo.SetSpotifyClient(Spotify.NewClient())
	FirstGo.InitializeServer()
	url, err := FirstGo.SendTest()
	if err != nil {
		fmt.Printf("error: %s \n", err)
		return
	}
	fmt.Printf("visit: %s \n", url)
	for {

	}
}
