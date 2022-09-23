package main

// This module is not used, because apparently blizzard API does not provide gem sockets on items.
// Just in case this is useful later it is being kept for now.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type BlizzardAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// There are more fields, these are just the ones we care about
type BlizzardItemResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Quality struct {
		Type string `json:"type"`
	} `json:"quality"`

	ItemClass struct {
		Id int `json:"id"`
	} `json:"item_class"`

	ItemSubclass struct {
		Id int `json:"id"`
	} `json:"item_subclass"`

	InventoryType struct {
		Type string `json:"type"`
	} `json:"inventory_type"`

	PreviewItem struct {
		Armor struct {
			Value int `json:"value"`
		} `json:"armor"`

		Stats []struct {
			Type struct {
				Type string `json:"type"`
			} `json:"type"`
			Value int `json:"value"`
		} `json:"stats"`
	} `json:"preview_item"`
}

func (item BlizzardItemResponse) GetStatValue(statType string) int {
	for _, stat := range item.PreviewItem.Stats {
		if stat.Type.Type == statType {
			return stat.Value
		}
	}
	return 0
}

func getAccessToken(clientId string, clientSecret string) string {
	url := "https://us.battle.net/oauth/token"

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	postData := []byte(`grant_type=client_credentials`)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal(err)
	}

	request.SetBasicAuth(clientId, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	result, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer result.Body.Close()

	resultBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}

	authResponse := BlizzardAuthResponse{}
	err = json.Unmarshal(resultBody, &authResponse)
	if err != nil {
		log.Fatal(err)
	}

	return authResponse.AccessToken
}

func getItemData(itemId int, accessToken string) BlizzardItemResponse {
	url := fmt.Sprintf("https://us.api.blizzard.com/data/wow/item/%d?namespace=static-classic-us&locale=en_US&access_token=%s", itemId, accessToken)

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer result.Body.Close()

	resultBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(string(resultBody))
	itemResponse := BlizzardItemResponse{}
	err = json.Unmarshal(resultBody, &itemResponse)
	if err != nil {
		log.Fatal(err)
	}

	return itemResponse
}
