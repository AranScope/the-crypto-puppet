package coinbase

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-make-me-rich/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultTimeout    = time.Second * 10
	productionBaseURL = "https://api.pro.coinbase.com"
	sandboxBaseURL    = "https://api-public.sandbox.pro.coinbase.com"
)

var errorPrefixByStatusCode = map[int]errors.ErrorPrefix{
	http.StatusBadRequest:          errors.ErrBadRequest,
	http.StatusUnauthorized:        errors.ErrUnauthorized,
	http.StatusForbidden:           errors.ErrForbidden,
	http.StatusNotFound:            errors.ErrNotFound,
	http.StatusInternalServerError: errors.ErrInternalServerError,
}

type Client struct {
	accessKey        string
	accessSecret     string
	accessPassPhrase string
	client           *http.Client
	baseURL          string
}

type requestError struct {
	Message string `json:"message"`
}

func NewClient(accessKey, accessSecret, accessPassPhrase string) *Client {
	return &Client{
		accessKey:        accessKey,
		accessSecret:     accessSecret,
		accessPassPhrase: accessPassPhrase,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: productionBaseURL,
	}
}

func NewSandboxClient(accessKey, accessSecret, accessPassPhrase string) *Client {
	return &Client{
		accessKey:        accessKey,
		accessSecret:     accessSecret,
		accessPassPhrase: accessPassPhrase,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: sandboxBaseURL,
	}
}

type ResponseFuture struct {
	err  error
	body io.ReadCloser
}

func (c *Client) authHeaders(method, path string, body []byte, t time.Time) (map[string]string, error) {
	accessTimestamp := fmt.Sprintf("%d", t.Unix())
	preHashString := accessTimestamp + method + path + string(body)
	decodedKey, err := base64.StdEncoding.DecodeString(c.accessSecret)
	if err != nil {
		return nil, err
	}

	// create a sha256 mac with the secret
	h := hmac.New(sha256.New, decodedKey)

	// sign the prehash with the hmac and base64 encode the result
	_, err = h.Write([]byte(preHashString))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"CB-ACCESS-KEY":        c.accessKey,
		"CB-ACCESS-SIGN":       base64.StdEncoding.EncodeToString(h.Sum(nil)),
		"CB-ACCESS-TIMESTAMP":  accessTimestamp,
		"CB-ACCESS-PASSPHRASE": c.accessPassPhrase,
	}, nil
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

	headers, err := c.authHeaders(method, path, marshaledBody, time.Now().UTC())
	if err != nil {
		return &ResponseFuture{err: err}
	}

	for k, v := range headers {
		req.Header[k] = []string{v}
	}

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

	bytes, err := ioutil.ReadAll(c.body)
	if err != nil {
		return err
	}

	if c.err != nil {
		v := &requestError{}
		err := json.Unmarshal(bytes, v)
		if err != nil {
			return err
		}

		return errors.New("", v.Message)
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(bytes, v)
}
