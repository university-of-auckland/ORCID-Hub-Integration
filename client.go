// RESTful API clietns supporting both APIKEY and Client Credential authenticaiton
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RESTClient ...
type RESTClient interface {
	GetAccessToken(url string) error
	Get(url string, resp interface{}) error
	Post(url string, body interface{}, resp interface{}) error
}

// Service - RESTfull service implementation
type Client struct {
	http.Client
	AccessToken, BaseURL, ApiKey, ClientID, ClientSecret string
}

func (c *Client) GetAccessToken(url string) error {
	var token struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	url = c.BaseURL + "/" + url
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=client_credentials", c.ClientID, c.ClientSecret))))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = c.execute(req, &token)
	if err != nil {
		return err
	}
	c.AccessToken = token.AccessToken
	return nil
}

func (c *Client) execute(req *http.Request, resp interface{}) error {

	if c.ApiKey != "" {
		req.Header.Set("apikey", c.ApiKey)
	} else if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}

	req.Header.Set("Accept", "application/json")
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Println("*****************")
	log.Println("URL:", req.URL, "/", req.URL.RequestURI())
	log.Println(string(body))
	log.Println("*****************")
	err = json.Unmarshal(body, resp)
	return err
}

func (c *Client) Get(url string, resp interface{}) error {
	url = c.BaseURL + "/" + url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return c.execute(req, resp)
}

func (c *Client) Post(url string, body interface{}, resp interface{}) error {
	url = c.BaseURL + "/" + url
	jb, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jb))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.execute(req, resp)
}

func (c *Client) Put(url string, body interface{}, resp interface{}) (err error) {
	url = c.BaseURL + "/" + url
	var jb []byte
	switch body.(type) {
	case string:
		jb = []byte(body.(string))
	default:
		jb, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jb))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.execute(req, resp)
}
