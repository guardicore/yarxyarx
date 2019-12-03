package http

import (
    "yarxyarx"
    "strings"
    "io"
    "net/http"
    "github.com/aws/aws-xray-sdk-go/xray"
    "net/url"
    "reflect"
)

type Client http.Client

var DefaultClient = &Client{}

func Get(url string) (resp *Response, err error) {
	return DefaultClient.Get(url)
}
func Post(url string, contentType string, body io.Reader) (resp *Response, err error) {
	return DefaultClient.Post(url, contentType, body)
}
func PostForm(url string, data url.Values) (resp *Response, err error) {
	return DefaultClient.PostForm(url, data)
}
func Head(url string) (resp *Response, err error) {
	return DefaultClient.Head(url)
}

func (c* Client) doXrayWrap() {
    *c = *(*Client)(xray.Client((*http.Client)(c)))
}

func (c* Client) wrapOnCall() {
    if c.Transport == nil {
        c.doXrayWrap()
    } else {
        transportType := reflect.TypeOf(c.Transport)
        // this is hideous but it's the quickest way I've found to make sure this client isn't already wrapped by xray.Client()
        if transportType.String() != "*xray.roundtripper" {
            c.doXrayWrap()
        }
    }
}

func (c *Client) Get(url string) (resp *Response, err error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
    c.wrapOnCall()

    f := (*http.Client).Do
    reqCtx := req.Context()
    wrappedCtx := yarxyarx.WithXrayContext(reqCtx)
    resp, err := f((*http.Client)(c), req.WithContext(wrappedCtx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil {
		select {
		case <-wrappedCtx.Done():
			err = wrappedCtx.Err()
		default:
		}
	}
    return resp, err
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *Response, err error) {
	req, err := NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *Client) Head(url string) (resp *Response, err error) {
	req, err := NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

