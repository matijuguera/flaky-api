package httpretry

import (
	"errors"
	apierror "flaky-api/apierror"
	"fmt"
	"net/http"
	"time"
)

const DEFAULT_MAX_RETRIES = 3

var (
	ErrUnexpectedMethod          = errors.New("unexpected client method, must be Get")
	ErrMaxAmoutOfRetries         = errors.New("max amount of retries reached, no response was found")
	ErrMaxRetriesGreaterThanZero = errors.New("max retries must be greater than 0")
)

type BackoffStrategy func(retry int) time.Duration
type Client struct {
	HttpClient HTTPClient
	MaxRetries int
	Backoff    BackoffStrategy
}

// params represents all the params needed to run http client calls
type params struct {
	method string
	url    string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ExponentialBackoff returns ever increasing backoffs by a power of 2
func ExponentialBackoff(i int) time.Duration {
	return time.Duration(1<<uint(i)) * time.Second
}

// LinearBackoff returns increasing durations, each a second longer than the last
func LinearBackoff(i int) time.Duration {
	return time.Duration(i) * time.Second
}

// DefaultBackoff always returns 1 second
func DefaultBackoff(_ int) time.Duration {
	return 1 * time.Second
}

// New constructs a new DefaultClient with sensible default values
func New(client HTTPClient) *Client {
	return &Client{
		MaxRetries: DEFAULT_MAX_RETRIES,
		Backoff:    DefaultBackoff,
		HttpClient: client,
	}
}

// Get provides the same functionality as http.Client.Get and creates its own constructor
func Get(url string) (resp *http.Response, err error) {
	c := New(&http.Client{})
	return c.Get(url)
}

// Get provides the same functionality as http.Client.Get
func (c *Client) Get(url string) (resp *http.Response, err error) {
	return c.doWithRetry(params{method: http.MethodGet, url: url})
}

// getRequest returns the valids request with the given method that doWithRetry supports
func getRequest(p params) (*http.Request, error) {
	switch p.method {
	case http.MethodGet:
		return http.NewRequest(p.method, p.url, nil)
	default:
		return nil, ErrUnexpectedMethod
	}

}

// doWithRetry provides a generic way to do the request with the given params
func (c *Client) doWithRetry(p params) (*http.Response, error) {

	request, err := getRequest(p)
	if err != nil {
		return nil, err
	}

	if c.MaxRetries <= 0 {
		return nil, ErrMaxRetriesGreaterThanZero
	}

	for i := 0; i <= c.MaxRetries; i++ {
		resp, err := c.HttpClient.Do(request)
		if err != nil {
			ErrMaxAmoutOfRetries = fmt.Errorf("%w \n -%v", ErrMaxAmoutOfRetries, err)
			continue
		}

		if resp != nil && resp.StatusCode == 200 {
			return resp, nil
		} else {
			ErrMaxAmoutOfRetries = fmt.Errorf("%w \n -%v", ErrMaxAmoutOfRetries, apierror.NewAPIError(resp.StatusCode, "http retry status error", p.url, resp.Status))
		}

		//avoids sleep at the last iteration
		if i != c.MaxRetries {
			time.Sleep(c.Backoff(i))
		}
	}

	return nil, ErrMaxAmoutOfRetries
}
