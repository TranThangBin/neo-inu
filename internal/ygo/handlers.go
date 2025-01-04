package ygo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const apiUrl = "https://db.ygoprodeck.com/api/v7/cardinfo.php"

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

func (r *RequestQueryBuilder) BuildUrlString() (string, error) {
	_url, err := url.Parse(apiUrl)
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

func RandomCard() (*Response, error) {
	urlStr, err := NewRequestQueryBuilder().
		SetNum(1).
		SetOffset(0).
		SetSort("random").
		SetCacheBust("").
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
