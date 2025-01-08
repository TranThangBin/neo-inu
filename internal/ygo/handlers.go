package ygo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const cardQueryUrl = "https://db.ygoprodeck.com/api/v7/cardinfo.php"
const randomCardUrl = "https://db.ygoprodeck.com/api/v7/randomcard.php."

type RequestQueryBuilder struct {
	queries map[string]string
}

func (r *RequestQueryBuilder) SetNum(num int) *RequestQueryBuilder {
	r.queries["num"] = strconv.Itoa(num)
	return r
}

func (r *RequestQueryBuilder) SetOffset(offset int) *RequestQueryBuilder {
	r.queries["offset"] = strconv.Itoa(offset)
	return r
}

func (r *RequestQueryBuilder) SetSort(sort string) *RequestQueryBuilder {
	r.queries["sort"] = sort
	return r
}

func (r *RequestQueryBuilder) SetCacheBust(cacheBust string) *RequestQueryBuilder {
	r.queries["cachebust"] = cacheBust
	return r
}

func (r *RequestQueryBuilder) SetFName(fname string) *RequestQueryBuilder {
	r.queries["fname"] = fname
	return r
}

func (r *RequestQueryBuilder) PushQueries(queries map[string]string) *RequestQueryBuilder {
	for k, v := range queries {
		if _, exists := r.queries[k]; !exists {
			r.queries[k] = v
		}
	}
	return r
}

func (r *RequestQueryBuilder) BuildUrlString() (string, error) {
	_url, err := url.Parse(cardQueryUrl)
	if err != nil {
		return "", err
	}
	q := _url.Query()
	for k, v := range r.queries {
		q.Add(k, v)
	}
	_url.RawQuery = q.Encode()
	return _url.String(), nil
}

func NewRequestQueryBuilder() *RequestQueryBuilder {
	return &RequestQueryBuilder{
		queries: make(map[string]string),
	}
}

func SearchRandomCard() (*Response, error) {
	resp, err := http.Get(randomCardUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ygoResponse := &Response{}
	err = json.Unmarshal(data, ygoResponse)
	if err != nil {
		return nil, err
	}
	return ygoResponse, err
}

func SearchCard(queries map[string]string) (*Response, error) {
	urlStr, err := NewRequestQueryBuilder().
		PushQueries(queries).
		BuildUrlString()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ygoResponse := &Response{}
	err = json.Unmarshal(data, ygoResponse)
	if err != nil {
		return nil, err
	}
	return ygoResponse, err
}
