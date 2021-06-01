package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	timeout = time.Second * 10
	url     = "https://updates.jenkins.io/update-center.actual.json"
)

// ClientOption represents an argument to NewClient.
type ClientOption = func(http.RoundTripper) http.RoundTripper

// NewClient initializes a Client.
func NewClient(opts ...ClientOption) *Client {
	tr := http.DefaultTransport
	for _, opt := range opts {
		tr = opt(tr)
	}

	h := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	client := &Client{http: h}

	return client
}

func enableTracing() ClientOption {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &Tracer{RoundTripper: rt}
	}
}

// ReplaceTripper substitutes the underlying RoundTripper with a custom one.
func ReplaceTripper(tr http.RoundTripper) ClientOption {
	return func(http.RoundTripper) http.RoundTripper {
		return tr
	}
}

// Client facilitates making HTTP requests to the GitHub API.
type Client struct {
	http *http.Client
}

// GET performs a REST request and parses the response.
func (c Client) GET(version string, data interface{}) error {
	requestURL := url
	if version != "" {
		requestURL += "?version=" + version
	}
	logrus.Debugf("checking url: %s", requestURL)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		logrus.Debugf("failed with resp code %d", resp.StatusCode)
		return handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

// REST performs a REST request and parses the response.
func (c Client) REST(url string, body io.Reader, data interface{}) error {
	logrus.Debugf("calling url: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		logrus.Debugf("failed with resp code %d", resp.StatusCode)
		return handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

func handleHTTPError(resp *http.Response) error {
	var message string
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}
	message = string(body)
	return fmt.Errorf("http error, '%s' failed (%d): '%s'", resp.Request.URL, resp.StatusCode, message)
}

// BasicClient returns a basic implemenation of a http client.
func BasicClient() *Client {
	var opts []ClientOption

	// for testing purposes one can enable tracing of API calls.
	httpTrace := os.Getenv("HTTP_TRACE")

	if httpTrace == "1" || httpTrace == "on" || httpTrace == "true" {
		opts = append(opts, enableTracing())
	}

	return NewClient(opts...)
}
