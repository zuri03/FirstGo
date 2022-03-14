package server

import (
	"encoding/json"
)

func generalStats(u []userInfo) (string, error) {
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
	json, err := json.Marshal(topArtists)
	if err != nil {
		return "", err
	}
	return string(json), nil
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
