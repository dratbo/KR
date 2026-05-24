package clients

import (
	"encoding/json"
	"net/http"
)

type DataClient struct {
	baseURL string
	client  *http.Client
}

func NewDataClient(baseURL string) *DataClient {
	return &DataClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type Item struct {
	ClassName   string `json:"class_name"`
	DisplayName string `json:"display_name"`
}

func (c *DataClient) GetItems() ([]Item, error) {
	resp, err := c.client.Get(c.baseURL + "/api/items")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
