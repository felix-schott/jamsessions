package geocoding

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type httpClientWithRateLimit struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
	UserAgent   string
}

// Do dispatches the HTTP request to the network
func (c *httpClientWithRateLimit) Do(req *http.Request) (*http.Response, error) {
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	// Comment out the below 5 lines to turn off ratelimiting
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Returns http client with a ratelimiter
func NewHttpClient(rl *rate.Limiter, userAgent string) *httpClientWithRateLimit {
	var tr = &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	var client = &http.Client{Transport: tr}
	c := &httpClientWithRateLimit{
		client:      client,
		Ratelimiter: rl,
		UserAgent:   userAgent,
	}
	return c
}
