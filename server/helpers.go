package server

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
)

func generalStats(u []userInfo) (map[string]int, error) {
	//A map to map artist name to number of times said aritist occurs in user library
	artists := make(map[string]int)
	for _, info := range u {
		for _, item := range info.Items {
			for _, artist := range item.Artists {
				if _, ok := artists[artist.Name]; !ok {
					artists[artist.Name] = 1
				} else {
					artists[artist.Name]++
				}
			}
		}
	}
	topArtists := getTopArtists(artists)
	return topArtists, nil
}

func getTopArtists(artists map[string]int) map[string]int {
	var topFive = make(map[string]int)
	swapMin := func(currentOccurrenceCount int, currentArtist string) {
		loser := currentArtist
		currentMin := currentOccurrenceCount
		for k, v := range topFive {
			if v < currentMin {
				currentMin = v
				loser = k
			}
		}
		if loser != currentArtist {
			delete(topFive, loser)
			topFive[currentArtist] = currentOccurrenceCount
		}
	}
	//PLACEHOLDER
	for k, v := range artists {
		if len(topFive) < 5 {
			topFive[k] = v
		} else {
			swapMin(v, k)
		}
	}
	return topFive
}

func generateAuthorizationCodeUrl(redirectUri string, scopes ...string) string {
	generateRandomState := func() string {
		const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		bytes := make([]byte, 16)
		for idx := range bytes {
			bytes[idx] = chars[rand.Intn(len(chars))]
		}
		return string(bytes)
	}

	uri := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&edirect_uri=%s&state=%s",
		url.QueryEscape(os.Getenv("CLIENT_ID")),
		url.QueryEscape(redirectUri),
		generateRandomState())

	var scopeStr string
	if len(scopes) > 0 && scopes != nil {
		scopeStr = ""

		for _, scope := range scopes {
			fmt.Printf("scope found => %sn", scope)
			scopeStr += fmt.Sprintf("%s ", scope)
		}

		uri = fmt.Sprintf("%s&scope=%s", uri, url.QueryEscape(scopeStr))
	}

	return uri
}
