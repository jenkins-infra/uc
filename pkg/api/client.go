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
	timeout = time.Second * 5
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

// AddHeader turns a RoundTripper into one that adds a request header.
func AddHeader(name, value string) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &funcTripper{roundTrip: func(req *http.Request) (*http.Response, error) {
			host := req.Host
			logrus.Debugf("sending request to host '%s'", host)
			if name == "Authorization" {
				if host == "api.github.com" {
					logrus.Debugf("Adding Authorization Header %s=%s", name, value)
					req.Header.Add(name, value)
				}
			} else {
				logrus.Debugf("Adding Header %s=%s", name, value)
				req.Header.Add(name, value)
			}

			return tr.RoundTrip(req)
		}}
	}
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

type funcTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (tr funcTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.roundTrip(req)
}

// Client facilitates making HTTP requests to the GitHub API.
type Client struct {
	http *http.Client
}

// REST performs a REST request and parses the response.
func (c Client) GET(version string, data interface{}) error {
	logrus.Debugf("checking url: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
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

func BasicClient() *Client {
	var opts []ClientOption

	// for testing purposes one can enable tracing of API calls.
	httpTrace := os.Getenv("HTTP_TRACE")

	if httpTrace == "1" || httpTrace == "on" || httpTrace == "true" {
		opts = append(opts, enableTracing())
	}

	return NewClient(opts...)
}
