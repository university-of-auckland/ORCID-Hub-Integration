// RESTful API clients supporting both APIKEY and Client Credential authentication
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
)

// RESTClient ...
type RESTClient interface {
	GetAccessToken(url string) error
	Get(url string, resp interface{}) error
	Post(url string, body interface{}, resp interface{}) error
}

// Client - RESTfull service implementation
type Client struct {
	http.Client
	accessToken, baseURL, apiKey, clientID, clientSecret string
	jsonBody                                             []byte
}

var accessTokenIsOnTheWay sync.Mutex

func setupAPIClients() {
	if api.apiKey == "" {
		api.apiKey = os.Getenv("API_KEY")
		api.baseURL = APIBaseURL
	}

	// Ensure that two guys don't try both to get a token (data race)
	accessTokenIsOnTheWay.Lock()
	if oh.accessToken == "" {
		gotAccessTokenWG.Add(1)
		go func() {
			oh.clientID = os.Getenv("CLIENT_ID")
			oh.clientSecret = os.Getenv("CLIENT_SECRET")
			oh.baseURL = OHBaseURL
			err := oh.getAccessToken("oauth/token")
			if err != nil || oh.accessToken == "" {
				log.Fatal("filed to authorize with the client credentials", zap.Error(err))
			}
			accessTokenIsOnTheWay.Unlock()
			gotAccessTokenWG.Done()
		}()
	} else {
		accessTokenIsOnTheWay.Unlock()
	}
}

func (c *Client) getAccessToken(url string) error {
	var token struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	url = c.baseURL + "/" + url
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=client_credentials", c.clientID, c.clientSecret))))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = c.execute(req, &token)
	if err != nil {
		return err
	}
	c.accessToken = token.AccessToken
	return nil
}

func (c *Client) execute(req *http.Request, resp interface{}) error {

	if c.apiKey != "" {
		req.Header.Set("apikey", c.apiKey)
	} else if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	req.Header.Set("Accept", "application/json")
	if req.Method != "GET" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	log.Debug("*****************")
	log.Debugf("%s %q %d %q", req.Method, req.URL.RequestURI(), r.StatusCode, r.Status)
	if resp != nil && r.StatusCode == http.StatusOK {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, resp)
		if err == nil {
			output, _ := json.MarshalIndent(resp, "", "    ")
			log.Debug(string(output))
		} else {
			log.Debug(string(body))
		}
		log.Debug("*****************")
		return err
	}
	return nil
}

func (c *Client) get(url string, resp interface{}) error {
	url = c.baseURL + "/" + url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return c.execute(req, resp)
}

func (c *Client) prepare(method, url string, body interface{}) (req *http.Request, err error) {
	url = c.baseURL + "/" + url
	if body == nil {
		return http.NewRequest(method, url, nil)
	}
	switch body.(type) {
	case string:
		c.jsonBody = []byte(body.(string))
	default:
		c.jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	return http.NewRequest(method, url, bytes.NewBuffer(c.jsonBody))
}

func (c *Client) do(method, url string, body interface{}, resp interface{}) error {
	req, err := c.prepare(method, url, body)
	if err != nil {
		return err
	}
	return c.execute(req, resp)
}

func (c *Client) put(url string, body interface{}, resp interface{}) error {
	return c.do("PUT", url, body, resp)
}

func (c *Client) post(url string, body interface{}, resp interface{}) error {
	return c.do("POST", url, body, resp)
}

func (c *Client) patch(url string, body interface{}, resp interface{}) error {
	return c.do("PATCH", url, body, resp)
}
