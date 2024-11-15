package geocoding

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// most code taken from https://gist.github.com/MelchiSalins/27c11566184116ec1629a0726e0f9af5

type httpClientWithRateLimit struct {
	client      *http.Client
	RateLimiter *rate.Limiter
	TimeOut     float32
	UserAgent   string
}

// Do dispatches the HTTP request to the network
func (c *httpClientWithRateLimit) Do(req *http.Request) (*http.Response, error) {
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	// rate limiting
	if c.RateLimiter != nil {
		ctx := context.Background()
		err := c.RateLimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
		if err != nil {
			return nil, err
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.TimeOut != 0 {
		time.Sleep(time.Millisecond * time.Duration(c.TimeOut*1000))
	}
	return resp, nil
}

// Returns http client with a ratelimiter
// optionally pass a rate.Limiter object OR a timeout value in seconds (after each request, the process will be blocked for the duration of the timeout)
// pass a user agent string using the third param
func NewHttpClient(rl *rate.Limiter, timeout float32, userAgent string) *httpClientWithRateLimit {
	var tr = &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	var client = &http.Client{Transport: tr}
	c := &httpClientWithRateLimit{
		client:      client,
		RateLimiter: rl,
		TimeOut:     timeout,
		UserAgent:   userAgent,
	}
	return c
}
