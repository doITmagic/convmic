package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type doFunc func(req *http.Request) (*http.Response, error)

type Client struct {
	UseAuth    bool
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
}

//create nre httpc client
func NewClient(apiKey, secretKey string, baseURL string) *Client {
	return &Client{
		UseAuth:    false,
		APIKey:     apiKey,
		SecretKey:  secretKey,
		UserAgent:  "convmic/golang",
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

//get the request and return it like []byte
func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// MakeReq HTTP request helper
func (c *Client) MakeReq(apiUrl string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?%s", c.BaseURL, apiUrl, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := doReq(req, c.HTTPClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}
