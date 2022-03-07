package spotifyClient

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
type spotifyApiClient struct {
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
)

func GetUserAnalysis(client *spotifyApiClient) (string, error) {
	if client.Token == nil || client.Token.AccessToken == "" {
		return "error", errors.New("error: client not authorized")
	}
	authorization := fmt.Sprintf("Bearer %s", client.Token.AccessToken)
	req, err := createRequest("GET", fmt.Sprintf("%s/me", baseApiUrl), nil, [2]string{"Authorization", authorization})
	if err != nil {
		return "error", err
	}
	resp, err := client.Client.Do(req)
	if err != nil {
		return "error", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json := string(body)

	return json, nil
}
func GetItemFromString(search string, itemType string, c *spotifyApiClient) (string, error) {

	if c.Token == nil || c.Token.AccessToken == "" {
		GetClientAccessToken("clientCredentials", "", c)
	}
	url := fmt.Sprintf("%s/search?q=%s&type=%s&limit=1", baseApiUrl, url.QueryEscape(search), url.QueryEscape(itemType))

	//Query search term endpoint
	authorizationValue := fmt.Sprintf("Bearer %s", c.Token.AccessToken)
	req, err := createRequest("GET", url, nil,
		[2]string{"Authorization", authorizationValue},
		[2]string{"Content-Type", "application/json"})

	resp, err := c.Client.Do(req)
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
func InitializeClient() *spotifyApiClient {
	return &spotifyApiClient{
		Client: &http.Client{},
		Token:  nil,
	}
}

//Make this function private after testing
func AuthorizationCode(redirectUri string) (string, error) {
	generateRandomState := func() string {
		const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		bytes := make([]byte, 16)
		for idx := range bytes {
			bytes[idx] = chars[rand.Intn(len(chars))]
		}
		return string(bytes)
	}
	url := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		url.QueryEscape(os.Getenv("CLIENT_ID")), url.QueryEscape(redirectUri), generateRandomState())
	return url, nil
}

func GetClientAccessToken(method string, code string, client *spotifyApiClient) (*spotifyApiClient, error) {

	form := url.Values{}
	var contentType string
	if method == "authorizationCode" {
		form.Add("code", code)
		form.Add("grant_type", "authorization_code")
		form.Add("redirect_uri", "http://localhost:8080/Analysis")
		contentType = "application/x-www-form-urlencoded"
	} else if method == "clientCredentials" {
		form.Add("grant_type", "client_credentials")
		contentType = "application/json"
	}

	c := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c)))

	req, err := createRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()),
		[2]string{"Authorization", authorization},
		[2]string{"Content-Type", contentType})
	if err != nil {
		fmt.Printf("error in request creation: %s \n", err)
		return client, err
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		fmt.Printf("error on token request: %s", err)
		return client, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error on reading body: %s", err)
		return client, err
	}
	body := string(bytes)

	var token spotifyAccessToken

	err = json.Unmarshal([]byte(body), &token)
	if err != nil {
		fmt.Printf("error formatting json: %s \n", err)
		return client, err
	}
	client.Token = &token
	return client, nil
}
