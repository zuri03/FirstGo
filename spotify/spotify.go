package spotify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//Client acess token
type SpotifyApiClient struct {
	Client *http.Client
}

const (
	baseApiUrl = "https://api.spotify.com/v1"
	authUrl    = "https://accounts.spotify.com"
)

func (s *SpotifyApiClient) GetUserTopItems(offset int, limit int, token string, expires time.Time) ([]byte, error) {
	if token == "" {
		return nil, errors.New("error: client not authorized")
	}
	authorization := fmt.Sprintf("Bearer %s", token)
	url := fmt.Sprintf("%s/me/top/tracks?limit=%d&offset=%d&time_range=long_term", baseApiUrl, limit, offset)
	req, err := createRequest("GET", url, nil, [2]string{"Authorization", authorization})

	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	/*
		fmt.Printf("status => %s \n", resp.Status)
		fmt.Printf("json => %s \n", string(body))
	*/
	return body, err
}

func (s *SpotifyApiClient) GetItemFromId(itemType string, token string, expires time.Time, ids ...string) ([]byte, error) {

	if token == "" {
		bytes, err := s.GetClientAccessToken("clientCredentials", "", "")
		if err != nil {
			return nil, err
		}
		token = string(bytes)
	}

	idsParam := strings.Join(ids, ",")
	var endpoint string
	switch itemType {
	case "album":
		if len(ids) == 1 {
			endpoint = fmt.Sprintf("%s/albums/%s", baseApiUrl,
				url.QueryEscape(ids[0]))
		} else {
			endpoint = fmt.Sprintf("%s/albums/%s", baseApiUrl,
				url.QueryEscape(idsParam))
		}
	case "artist":
		if len(ids) == 1 {
			endpoint = fmt.Sprintf("%s/artists/%s", baseApiUrl, url.QueryEscape(ids[0]))
		} else {
			endpoint = fmt.Sprintf("%s/artists/%s", baseApiUrl, url.QueryEscape(idsParam))
		}
	default:
		return nil, errors.New("error: incorrect item type only \"album\" or \"artits\" allowed")
	}

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (s *SpotifyApiClient) GetItemFromName(search string, itemType string, token string, expires time.Time) ([]byte, error) {

	if token == "" {
		bytes, err := s.GetClientAccessToken("clientCredentials", "", "")
		if err != nil {
			return nil, err
		}
		token = string(bytes)
	}

	endpoint := fmt.Sprintf("%s/search?q=%s&type=%s&limit=1", baseApiUrl,
		url.QueryEscape(search),
		url.QueryEscape(itemType))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (s *SpotifyApiClient) GetClientAccessToken(method string, redirectUri string, code string) ([]byte, error) {

	form := url.Values{}
	if method == "authorizationCode" {
		form.Add("code", code)
		form.Add("grant_type", "authorization_code")
		form.Add("redirect_uri", redirectUri)
	} else if method == "clientCredentials" {
		form.Add("grant_type", "client_credentials")
	}

	c := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c)))

	url := fmt.Sprintf("%s/api/token", authUrl)
	req, err := createRequest("POST", url, strings.NewReader(form.Encode()),
		[2]string{"Authorization", authorization},
		[2]string{"Content-Type", "application/x-www-form-urlencoded"})

	if err != nil {
		fmt.Printf("error in request creation: %s \n", err)
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Printf("error on token request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error on reading body: %s", err)
		return nil, err
	}

	return bytes, nil
}

//Make this function private after testing

func (s *SpotifyApiClient) GetRelatedArtist(artistId string, token string, expires time.Time) ([]byte, error) {
	if token == "" {
		return nil, errors.New("error: client not authenticated")
	}

	endpoint := fmt.Sprintf("%s/artists/%s/related-artists", baseApiUrl,
		url.QueryEscape(artistId))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (s *SpotifyApiClient) GetTracksFromArtist(artistId string, token string, expires time.Time) ([]byte, error) {
	if token == "" {
		return nil, errors.New("error: client not authenticated")
	}

	endpoint := fmt.Sprintf("%s/artists/%s/top-tracks", baseApiUrl,
		url.QueryEscape(artistId))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (s *SpotifyApiClient) GetSavedTracks(token string, expires time.Time) ([]byte, error) {

	if token == "" {
		return nil, errors.New("error: client not authenticated")
	}

	endpoint := fmt.Sprintf("%s/me/tracks", baseApiUrl)

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (s *SpotifyApiClient) GetRecommendations(artistId []string, genres []string, trackIds []string, token string, expires time.Time) ([]byte, error) {
	if token == "" {
		return nil, errors.New("error: client not authenticated")
	}

	if len(genres) > 3 {
		genres = genres[:3]
	}

	if len(artistId) > 3 {
		artistId = artistId[:3]
	}

	if len(trackIds) > 3 {
		trackIds = trackIds[:3]
	}
	endpoint := fmt.Sprintf("%s/recommendations?seed_artists=%s&seed_genres=%s&seed_tracks=%s", baseApiUrl,
		url.QueryEscape(strings.Join(artistId, ",")),
		url.QueryEscape(strings.Join(genres, ",")),
		url.QueryEscape(strings.Join(trackIds, ",")))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", token)
	req, err := createRequest("GET", endpoint, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func createRequest(method string, url string, body io.Reader, headers ...[2]string) (*http.Request, error) {
	//For now it can only fetch resources
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		req.Header.Add(header[0], header[1])
	}

	return req, nil
}

//Client Authorization and Intialization Functions
func NewClient() *SpotifyApiClient {
	return &SpotifyApiClient{
		Client: &http.Client{},
	}
}
