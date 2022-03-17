package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	do         doFunc
}

func NewClient(apiKey, secretKey string, baseURL string) *Client {
	return &Client{
		UseAuth:    false,
		APIKey:     apiKey,
		SecretKey:  secretKey,
		UserAgent:  "convmic/golang",
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "convmic-golang ", log.LstdFlags),
	}
}

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
