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
	result, err := json.Marshal(artists)
	if err != nil {
		return "", nil
	}
	return string(result), nil
}

func getTopArtists(artists map[string]int) map[string]int {
	top := make(map[string]int)
	isTopThree := func(artist string, occurrences int) {
		if len(top) < 3 {
			top[artist] = occurrences
			return
		}
	}
	//PLACEHOLDER
	isTopThree("", 0)
	return nil
}
