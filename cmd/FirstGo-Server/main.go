package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	Server "github.com/zuri03/FirstGo/server"
	Spotify "github.com/zuri03/FirstGo/spotify"
)

func main() {

	cleanVariables := func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}

	data, err := os.ReadFile("../../.env")
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
	c := Spotify.NewClient()
	Server.InitializeServer(c)
	url := c.GenerateAuthorizationCodeUrl("http://localhost:8080/Analysis", "user-top-read")
	if err != nil {
		fmt.Printf("error: %s \n", err)
		return
	}
	fmt.Printf("visit: %s \n", url)
	for {

	}
}
