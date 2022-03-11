package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	Spotify "github.com/zuri03/FirstGo/pkg/spotify"
)

func TestItemIdEndpoint(t *testing.T) {

	client := Spotify.NewClient()
	SetSpotifyClient(client)

	t.Run("Get info returns error on null query parameters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/Info", nil)
		resp := httptest.NewRecorder()

		getItemInfo(resp, req)

		if resp.Code != 400 {
			t.Errorf("Expected error code 400 got: %d \n", resp.Code)
		}

		if result := resp.Body.String(); !strings.Contains(result, "error") {
			t.Errorf("Expected error message got: %s \n", result)
		}
	})

	t.Run("Get info returns valid response", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/Info?search=future", nil)
		resp := httptest.NewRecorder()

		getItemInfo(resp, req)

		if resp.Code != 200 {
			t.Errorf("Wanted success status code go: %d", resp.Code)
		}

		body := resp.Body.String()
		if body == "" {
			t.Errorf("Wanted json got: empty string")
		}

		var obj map[string]interface{}

		err := json.Unmarshal([]byte(body), &obj)
		if err != nil {
			t.Errorf("Wanted json format got: %s \n With Error: %s", body, err)
		}
	})
}
