package slash

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const HostURL string = "http://localhost:5231"

type Client struct {
	HostURL     string
	HTTPClient  *http.Client
	AccessToken string
}

type AuthResponse struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewClient(host *string, accessToken *string) (*Client, error) {
	c := Client{
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
		HostURL:     HostURL,
		AccessToken: *accessToken,
	}

	if host != nil {
		c.HostURL = *host
	}

	_, err := c.validateAccessToken()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.AccessToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
