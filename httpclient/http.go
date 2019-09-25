package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultHTTPTimeout     = 80 * time.Second
	maxNetworkRetriesDelay = 5000 * time.Millisecond
	minNetworkRetriesDelay = 500 * time.Millisecond
)

// Config represents the service config
type Config struct {
	APIURL            string
	HTTPClient        *http.Client
	MaxNetworkRetries int
}

// Client represent http client service
type Client struct {
	APIURL            string
	HTTPClient        *http.Client
	MaxNetworkRetries int
	useNormalSleep    bool
}

func (c *Client) shouldRetry(err error, res *http.Response, retry int) bool {
	if retry >= c.MaxNetworkRetries {
		return false
	}

	if err != nil {
		return true
	}

	return false
}

func (c *Client) sleepTime(numRetries int) time.Duration {
	if c.useNormalSleep {
		return 0
	}

	// exponentially backoff by 2^numOfRetries
	delay := minNetworkRetriesDelay + minNetworkRetriesDelay*time.Duration(1<<uint(numRetries))
	if delay > maxNetworkRetriesDelay {
		delay = maxNetworkRetriesDelay
	}

	// generate random jitter to prevent thundering herd problem
	jitter := rand.Int63n(int64(delay / 4))
	delay -= time.Duration(jitter)

	if delay < minNetworkRetriesDelay {
		delay = minNetworkRetriesDelay
	}

	return delay
}

// Do calls the api http request and parse the response into v
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error

	for retry := 0; ; {
		res, err = c.HTTPClient.Do(req)
		if !c.shouldRetry(err, res, retry) {
			break
		}

		sleepDuration := c.sleepTime(retry)
		retry++
		time.Sleep(sleepDuration)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// BuildResponse .
func (c *Client) BuildResponse(url string, statusCode int, jsonBody, resBody []byte, data interface{}) error {

	if statusCode != http.StatusOK {
		msg := fmt.Sprintf("Code : %d  Path: %s - Request: %s - Response : %s", statusCode, url, string(jsonBody), string(resBody))
		return errors.New(msg)
	}
	errDecode := json.Unmarshal(resBody, data)
	if errDecode != nil {
		return errDecode
	}
	return nil
}

// NewHTTPClient ...
func NewHTTPClient(config Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: defaultHTTPTimeout,
		}
	}

	return &Client{
		APIURL:            config.APIURL,
		HTTPClient:        config.HTTPClient,
		MaxNetworkRetries: config.MaxNetworkRetries,
	}
}
