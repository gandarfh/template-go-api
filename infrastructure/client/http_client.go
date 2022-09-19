package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type header struct {
	key   string
	value string
}

type params struct {
	key   string
	value string
}

type client struct {
	url     string
	method  string
	body    io.Reader
	params  []*params
	headers []*header
	cookies []*http.Cookie
}

func (c *client) Header(key string, value string) *client {
	c.headers = append(c.headers, &header{
		key:   key,
		value: value,
	})

	return c
}

func (c *client) AddCookie(cookie *http.Cookie) *client {
	c.cookies = append(c.cookies, cookie)

	return c
}

func (c *client) Params(key string, value string) *client {
	c.params = append(c.params, &params{key: key, value: value})

	queryParams := []string{"?"}

	last := len(c.params) - 1

	for i, item := range c.params {
		if i == last {
			queryParams = append(queryParams, item.key, "=", item.value)
			continue
		}

		queryParams = append(queryParams, item.key, "=", item.value, "&")
	}

	c.url = c.url + strings.Join(queryParams, "")

	return c
}

func (c *client) Body(b interface{}) *client {
	body, _ := json.Marshal(b)

	return &client{
		body: bytes.NewReader(body),
	}
}

func (c *client) Decode(decode any) (*http.Response, error) {
	request, err := http.NewRequest(c.method, c.url, c.body)

	if err != nil {
		return nil, err
	}

	for _, item := range c.headers {
		request.Header.Add(item.key, item.value)
	}

	for _, item := range c.cookies {
		request.AddCookie(item)
	}

	SERVER_READ_TIMEOUT, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	client := http.Client{
		Timeout: time.Second * time.Duration(SERVER_READ_TIMEOUT),
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)

	if decode != nil {
		if err := json.Unmarshal(data, decode); err != nil {
			return nil, err
		}
	}

	return response, nil
}

func Post(url string) *client {
	return &client{
		url:    url,
		method: http.MethodPost,
	}
}

func Get(url string) *client {
	return &client{
		url:    url,
		method: http.MethodGet,
	}
}

func Patch(url string) *client {
	return &client{
		url:    url,
		method: http.MethodPatch,
	}
}

func Delete(url string) *client {
	return &client{
		url:    url,
		method: http.MethodDelete,
	}
}

func Put(url string) *client {
	return &client{
		url:    url,
		method: http.MethodPut,
	}
}
