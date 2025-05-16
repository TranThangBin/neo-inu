package ygo

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	cardInfoUrlStr   = "https://db.ygoprodeck.com/api/v7/cardinfo.php"
	randomCardUrlStr = "https://db.ygoprodeck.com/api/v7/randomcard.php."
)

type ygoprodeckClient struct{}

func (yg *ygoprodeckClient) SearchRandomCard() (*YgoproDeckResponse, error) {
	return yg.get(randomCardUrlStr)
}

func (yg *ygoprodeckClient) SearchCard(queries map[string]string) (*YgoproDeckResponse, error) {
	cardQueryUrl, err := url.Parse(cardInfoUrlStr)
	if err != nil {
		return nil, err
	}
	query := cardQueryUrl.Query()
	for k, v := range queries {
		query.Add(k, v)
	}
	cardQueryUrl.RawQuery = query.Encode()
	return yg.get(cardQueryUrl.String())
}

func (yg *ygoprodeckClient) get(url string) (*YgoproDeckResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ygoResponse := &YgoproDeckResponse{}
	err = json.NewDecoder(resp.Body).Decode(ygoResponse)
	if err != nil {
		return nil, err
	}
	return ygoResponse, err
}
