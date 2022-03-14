package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	Spotify "github.com/zuri03/FirstGo/spotify"
)

func TestItemIdEndpoint(t *testing.T) {

	client := Spotify.NewClient()

	t.Run("Get info returns error on null query parameters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/Info", nil)
		resp := httptest.NewRecorder()

		handler := getItemHandler{spotify: client}
		handler.ServeHTTP(resp, req)

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

		handler := getItemHandler{spotify: client}
		handler.ServeHTTP(resp, req)

		if resp.Code != 200 {
			t.Errorf("Wanted success status code go: %d \n", resp.Code)
		}

		body := resp.Body.String()
		if body == "" {
			t.Errorf("Wanted json got: empty string \n")
		}

		var obj map[string]interface{}

		err := json.Unmarshal([]byte(body), &obj)
		if err != nil {
			t.Errorf("Wanted json format got: %s \n With Error: %s \n", body, err)
		}
	})
}

func TestAnalysisEndpoint(t *testing.T) {

	client := Spotify.NewClient()

	t.Run("Get user analysis returns valid response", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/Analysis", nil)
		resp := httptest.NewRecorder()

		handler := analysisHandler{spotify: client}
		handler.ServeHTTP(resp, req)

		res := resp.Result()

		if status := res.StatusCode; status != 301 {
			t.Errorf("Wanted status code 301 got: %d \n", status)
		}

		redirect, err := res.Location()
		if err != nil {
			t.Errorf("Wanted redirect url got: %s \n", err)
		}

		state := redirect.Query().Get("state")
		if state == "" {
			t.Errorf("State was not found in redirect url \n")
		}

		if len(state) != 16 {
			t.Errorf("Wanted state of length 16 got: %d => %s", len(state), state)
		}
	})
}
