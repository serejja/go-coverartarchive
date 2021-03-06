package coverartarchive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	WSUrl        string
	MaxRedirects int
}

func NewClientConfig() ClientConfig {
	return ClientConfig{
		WSUrl:        "http://coverartarchive.org",
		MaxRedirects: 30,
	}
}

type Client struct {
	client    *http.Client
	wsRootURL *url.URL
}

func NewClient(config ClientConfig) (*Client, error) {
	client := &http.Client{}
	c := &Client{
		client: client,
	}

	var err error
	c.wsRootURL, err = url.Parse(config.WSUrl)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) ReleaseCoverArt(mbid string) (*CoverArtResponse, error) {
	response := new(CoverArtResponse)
	return response, c.getBody(fmt.Sprintf("release/%s", mbid), response)
}

func (c *Client) ReleaseGroupCoverArt(mbid string) (*CoverArtResponse, error) {
	response := new(CoverArtResponse)
	return response, c.getBody(fmt.Sprintf("release-group/%s", mbid), response)
}

func (c *Client) ReleaseFrontCoverArtURL(mbid string) (string, error) {
	return c.getRedirectLocation(fmt.Sprintf("release/%s/front", mbid))
}

func (c *Client) ReleaseBackCoverArtURL(mbid string) (string, error) {
	return c.getRedirectLocation(fmt.Sprintf("release/%s/back", mbid))
}

func (c *Client) getBody(path string, data interface{}) error {
	response, err := c.get(path, true)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}

func (c *Client) getRedirectLocation(path string) (string, error) {
	response, err := c.get(path, false)
	if err != nil {
		return "", err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	locationHeader := response.Header.Get("Location")
	if locationHeader == "" {
		return "", ErrNotFound
	}

	return locationHeader, nil
}

func (c *Client) get(path string, followRedirects bool) (*http.Response, error) {
	successCode := 200
	if followRedirects {
		c.client.CheckRedirect = nil
	} else {
		c.client.CheckRedirect = noFollowRedirects
		successCode = 307
	}

	reqUrl := *c.wsRootURL
	reqUrl.Path = path

	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != successCode {
		defer resp.Body.Close()
		return nil, c.handleError(resp)
	}

	return resp, nil
}

func (c *Client) handleError(response *http.Response) error {
	switch response.StatusCode {
	case http.StatusBadRequest:
		return ErrInvalidMBID
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusServiceUnavailable:
		return ErrRateLimitReached
	default:
		errorResponse, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("Received error response with status %d", response.StatusCode)
		}
		return fmt.Errorf("Received error response with status %d and message %s", response.StatusCode, string(errorResponse))
	}
}

func noFollowRedirects(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}
