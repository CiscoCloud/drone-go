package drone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	token      string
	url        string
	isServer04 bool

	Commits    *CommitService
	Repos      *RepoService
	Users      *UserService
	HttpClient *http.Client
}

func NewClient(token, url string, client *http.Client) *Client {
	c := Client{
		token: token,
		url:   url,
	}

	c.Commits = &CommitService{&c}
	c.Repos = &RepoService{&c}
	c.Users = &UserService{&c}
	if client == nil {
		client = http.DefaultClient
	}
	c.HttpClient = client
	return &c
}

func NewClient04(token string, url string, client *http.Client) *Client {
	c := NewClient(token, url, client)
	c.isServer04 = true
	return c
}

var (
	ErrNotFound       = errors.New("Not Found")
	ErrForbidden      = errors.New("Forbidden")
	ErrBadRequest     = errors.New("Bad Request")
	ErrNotAuthorized  = errors.New("Unauthorized")
	ErrInternalServer = errors.New("Internal Server Error")
)

// runs an http.Request and parses the JSON-encoded http.Response,
// storing the result in the value pointed to by v.
func (c *Client) run(method, path string, in, out interface{}) error {

	// create the URI
	uri, err := url.Parse(c.url + path)
	if err != nil {
		return err
	}

	if len(uri.Scheme) == 0 {
		uri.Scheme = "http"
	}

	if len(c.token) > 0 {
		params := uri.Query()
		params.Add("access_token", c.token)
		uri.RawQuery = params.Encode()
	}

	// create the request
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return err
	}
	req.ProtoAtLeast(1, 1)
	req.Close = true
	req.ContentLength = 0

	// if data input is provided, serialize to JSON
	if in != nil {
		formIn, ok := in.(map[string]string)
		var buf *bytes.Buffer
		var contentType string
		if c.isServer04 && ok {
			contentType = "application/x-www-form-urlencoded"
			data := url.Values{}
			for key, val := range formIn {
				data.Set(key, val)
			}
			buf = bytes.NewBufferString(data.Encode())
		} else if bytesIn, ok := in.([]byte); c.isServer04 && ok {
			contentType = "text/plain"
			buf = bytes.NewBufferString(string(bytesIn[:]))
		} else {
			contentType = "application/json"
			inJson, err := json.Marshal(in)
			if err != nil {
				return err
			}
			buf = bytes.NewBuffer(inJson)
		}

		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(buf.Len())
		req.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
		req.Header.Set("Content-Type", contentType)
	}

	// make the request using the default http client
	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return err
	}

	// make sure we defer close the body
	defer resp.Body.Close()

	// Check for an http error status (ie not 200 StatusOK)
	switch resp.StatusCode {
	case 404:
		return ErrNotFound
	case 403:
		return ErrForbidden
	case 401:
		return ErrNotAuthorized
	case 400:
		return ErrBadRequest
	case 500:
		return ErrInternalServer
	}

	// Decode the JSON response
	if out != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading response body: %s", err))
		}
		err = json.Marshal(respBody, out)
		if err != nil {
			if outStr, ok := out.(*string); ok {
				*outStr = string(contents[:])
			} else {
				return err
			}
		}
	}

	return nil
}

// do makes an http.Request and returns the response
func (c *Client) do(method, path string) (*http.Response, error) {

	// create the URI
	uri, err := url.Parse(c.url + path)
	if err != nil {
		return nil, err
	}

	if len(uri.Scheme) == 0 {
		uri.Scheme = "http"
	}

	if len(c.token) > 0 {
		params := uri.Query()
		params.Add("access_token", c.token)
		uri.RawQuery = params.Encode()
	}

	// create the request
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.ProtoAtLeast(1, 1)
	req.Close = true
	req.ContentLength = 0

	// make the request using the default http client
	resp, err := c.HttpClient.Do(req)

	return resp, err
}
