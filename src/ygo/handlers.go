package ygo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const (
	cardInfoUrlStr   = "https://db.ygoprodeck.com/api/v7/cardinfo.php"
	randomCardUrlStr = "https://db.ygoprodeck.com/api/v7/randomcard.php."
)

func SearchRandomCard() (*Response, error) {
	resp, err := http.Get(randomCardUrlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ygoResponse := &Response{}
	err = json.NewDecoder(resp.Body).Decode(ygoResponse)
	if err != nil {
		return nil, err
	}
	return ygoResponse, err
}

func SearchCard(queries map[string]string) (*Response, error) {
	cardQueryUrl, err := url.Parse(cardInfoUrlStr)
	if err != nil {
		log.Printf("Card query url is malformed. {%v}\n", err)
	}
	query := cardQueryUrl.Query()
	for k, v := range queries {
		query.Add(k, v)
	}
	cardQueryUrl.RawQuery = query.Encode()
	resp, err := http.Get(cardQueryUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ygoResponse := &Response{}
	err = json.NewDecoder(resp.Body).Decode(ygoResponse)
	if err != nil {
		return nil, err
	}
	return ygoResponse, err
}
