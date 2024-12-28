package slash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetShortcut(shortcutID string) (*Shortcut, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/shortcuts/%s", c.HostURL, shortcutID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := Shortcut{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func (c *Client) CreateShortcut(shortcut Shortcut) (*Shortcut, error) {
	rb, err := json.Marshal(shortcut)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/shortcuts", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := Shortcut{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func (c *Client) UpdateShortcut(shortcutID string, shortcut Shortcut) (*Shortcut, error) {
	rb, err := json.Marshal(shortcut)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/shortcuts/%s?updateMask=name,link,title", c.HostURL, shortcutID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := Shortcut{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

//// DeleteOrder - Deletes an order
//func (c *Client) DeleteOrder(orderID string) error {
//	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/orders/%s", c.HostURL, orderID), nil)
//	if err != nil {
//		return err
//	}
//
//	body, err := c.doRequest(req)
//	if err != nil {
//		return err
//	}
//
//	if string(body) != "Deleted order" {
//		return errors.New(string(body))
//	}
//
//	return nil
//}
