package spotify

/*
*TODO: Convert all of these functions to methods
 */
//client should be reused on every request not recreated
import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//Client acess token
type SpotifyApiClient struct {
	Client *http.Client
	Token  *spotifyAccessToken
}
type spotifyAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const (
	baseApiUrl = "https://api.spotify.com/v1"
	authUrl    = "https://accounts.spotify.com"
)

func (s *SpotifyApiClient) GetUserAnalysis(offset int) ([]byte, error) {
	if s.Token == nil || s.Token.AccessToken == "" {
		return nil, errors.New("error: client not authorized")
	}
	authorization := fmt.Sprintf("Bearer %s", s.Token.AccessToken)
	url := fmt.Sprintf("%s/me/top/tracks?limit=50&offset=%d&time_range=long_term", baseApiUrl, offset)
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

func (s *SpotifyApiClient) GetItemFromString(search string, itemType string) (string, error) {

	if s.Token == nil || s.Token.AccessToken == "" {
		s.GetClientAccessToken("clientCredentials", "", "")
	}
	url := fmt.Sprintf("%s/search?q=%s&type=%s&limit=1", baseApiUrl, url.QueryEscape(search), url.QueryEscape(itemType))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", s.Token.AccessToken)
	req, err := createRequest("GET", url, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	resp, err := s.Client.Do(req)
	if err != nil {
		return "error", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json := string(body)

	return json, nil
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
		Token:  nil,
	}
}

//Make this function private after testing
func (s *SpotifyApiClient) GenerateAuthorizationCodeUrl(redirectUri string, scopes ...string) string {
	generateRandomState := func() string {
		const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		bytes := make([]byte, 16)
		for idx := range bytes {
			bytes[idx] = chars[rand.Intn(len(chars))]
		}
		return string(bytes)
	}

	uri := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		authUrl,
		url.QueryEscape(os.Getenv("CLIENT_ID")),
		url.QueryEscape(redirectUri),
		generateRandomState())

	var scopeStr string
	if len(scopes) > 0 {
		scopeStr = ""

		for _, scope := range scopes {
			scopeStr += fmt.Sprintf("%s ", scope)
		}

		uri = fmt.Sprintf("%s&scope=%s", uri, url.QueryEscape(scopeStr))
	}

	return uri
}

func (s *SpotifyApiClient) GetClientAccessToken(method string, redirectUri string, code string) error {

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
		return err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Printf("error on token request: %s", err)
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error on reading body: %s", err)
		return err
	}
	body := string(bytes)

	var token spotifyAccessToken

	err = json.Unmarshal([]byte(body), &token)
	if err != nil {
		fmt.Printf("error formatting json: %s \n", err)
		return err
	}
	s.Token = &token
	return nil
}
