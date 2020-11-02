package coinmarketcap

import (
	"bytes"
	"encoding/json"
	"github.com/AranScope/the-crypto-puppet/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultTimeout    = time.Second * 10
	productionBaseURL = "https://pro-api.coinmarketcap.com"
)

var errorPrefixByStatusCode = map[int]errors.ErrorPrefix{
	http.StatusBadRequest:          errors.ErrBadRequest,
	http.StatusUnauthorized:        errors.ErrUnauthorized,
	http.StatusForbidden:           errors.ErrForbidden,
	http.StatusNotFound:            errors.ErrNotFound,
	http.StatusInternalServerError: errors.ErrInternalServerError,
}

type Client struct {
	accessSecret string
	client       *http.Client
	baseURL      string
}

func NewClient(accessSecret string) *Client {
	return &Client{
		accessSecret: accessSecret,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: productionBaseURL,
	}
}

type ResponseFuture struct {
	err  error
	body io.ReadCloser
}

func (c *Client) Request(method string, path string, body interface{}) *ResponseFuture {
	var marshaledBody []byte
	if body != nil {
		var err error
		marshaledBody, err = json.Marshal(body)
		if err != nil {
			return &ResponseFuture{err: err}
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(marshaledBody))
	if err != nil {
		return &ResponseFuture{err: err}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req.Header.Set("X-CMC_PRO_API_KEY", c.accessSecret)

	rsp, err := c.client.Do(req)
	if err != nil {
		return &ResponseFuture{err: err}
	}

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		prefix, ok := errorPrefixByStatusCode[rsp.StatusCode]
		if !ok {
			prefix = errors.ErrUnknown
		}
		return &ResponseFuture{err: errors.New(prefix, ""), body: rsp.Body}
	}

	return &ResponseFuture{body: rsp.Body}
}

func (c *ResponseFuture) DecodeResponse(v interface{}) error {
	if c.body == nil {
		return c.err
	}

	if c.err != nil {
		return c.err
	}

	bytes, err := ioutil.ReadAll(c.body)
	if err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(bytes, v)
}
