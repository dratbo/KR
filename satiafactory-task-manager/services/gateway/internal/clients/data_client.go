package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DataClient struct {
	baseURL string
	client  *http.Client
}

func NewDataClient(baseURL string) *DataClient {
	if baseURL == "" {
		baseURL = os.Getenv("DATA_SERVICE_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:8083"
	}
	return &DataClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type Recipe struct {
	ClassName   string       `json:"class_name"`
	DisplayName string       `json:"display_name"`
	Ingredients []Ingredient `json:"ingredients"`
	Products    []Product    `json:"products"`
	Duration    float64      `json:"duration"`
}

type Ingredient struct {
	ItemClassName string  `json:"item_class_name"`
	Amount        float64 `json:"amount"`
}

type Product struct {
	ItemClassName string  `json:"item_class_name"`
	Amount        float64 `json:"amount"`
}

type Item struct {
	ClassName   string `json:"class_name"`
	DisplayName string `json:"display_name"`
}

func (c *DataClient) SearchRecipes(query string) ([]Recipe, error) {
	u := c.baseURL + "/api/recipes/search?q=" + url.QueryEscape(query)
	resp, err := c.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search recipes: %s", resp.Status)
	}
	var recipes []Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (c *DataClient) GetRecipe(className string) (*Recipe, error) {
	u := c.baseURL + "/api/recipes/" + url.PathEscape(className)
	resp, err := c.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get recipe: %s", resp.Status)
	}
	var recipe Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
		return nil, err
	}
	return &recipe, nil
}

// ItemIconURL builds an image URL from SatisfactoryTools (greeny) by item class name.
func ItemIconURL(className string) string {
	if className == "" {
		return ""
	}
	slug := strings.TrimSuffix(className, "_C")
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, "_", "-")
	return "https://raw.githubusercontent.com/greeny/SatisfactoryTools/dev/images/" + slug + ".png"
}
