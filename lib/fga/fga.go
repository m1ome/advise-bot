package fga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Client struct {
		http *http.Client
	}

	advise struct {
		Id    int64  `json:"id"`
		Text  string `json:"text"`
		Sound string `json:"sound"`
	}
)

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c Client) Random() (string, error) {
	req, err := http.NewRequest("GET", "http://fucking-great-advice.ru/api/random", nil)
	if err != nil {
		return "", err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("error getting advise [%d]: %s", res.StatusCode, string(body))
	}

	var a advise
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		return "", err
	}

	return a.Text, nil
}
