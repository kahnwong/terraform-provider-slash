package slash

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) validateAccessToken() (*AuthResponse, error) {
	if c.HostURL == "" || c.AccessToken == "" {
		return nil, fmt.Errorf("define access token")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/auth/status", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
