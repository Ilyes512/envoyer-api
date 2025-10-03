package envoyerapi

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseUrl   string = "https://envoyer.io/api/"
	DefaultUserAgent string = "envoyer-api-go/0.0.0"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
	auth       Auth
	userAgent  string
}

type options func(*Client)

func WithBaseUrl(baseUrl string) options {
	return func(c *Client) {
		parsed, err := url.Parse(baseUrl)
		if err == nil {
			panic(err)
		}
		c.baseUrl = parsed
	}
}

func WithAuth(auth Auth) options {
	return func(c *Client) {
		c.auth = auth
	}
}

func WithUserAgent(userAgent *string) options {
	return func(c *Client) {
		if userAgent == nil {
			c.userAgent = DefaultUserAgent
		} else {
			c.userAgent = *userAgent
		}
	}
}

func WithDefaultClient() options {
	return func(c *Client) {
		c.httpClient = &http.Client{
			Timeout: 3 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   2 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:           100,
				MaxIdleConnsPerHost:    10,
				IdleConnTimeout:        60 * time.Second,
				TLSHandshakeTimeout:    3 * time.Second,
				ExpectContinueTimeout:  0,
				ForceAttemptHTTP2:      true,
				MaxResponseHeaderBytes: 1 << 20, // 1 MiB
			},
		}
	}
}

func WithClient(httpClient *http.Client) options {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func NewClient(opts ...options) *Client {
	client := &Client{}

	for _, opt := range opts {
		opt(client)
	}

	if client.httpClient == nil {
		WithDefaultClient()(client)
	}

	if client.baseUrl == nil {
		parsed, err := url.Parse(DefaultBaseUrl)
		if err != nil {
			panic(err)
		}
		client.baseUrl = parsed
	}

	return client
}

func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	fullUrl, err := c.baseUrl.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullUrl.String(), body)
	if err != nil {
		return nil, err
	}

	c.auth.Apply(req)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) BaseUrl() *url.URL {
	return c.baseUrl
}

func (c *Client) UserAgent() string {
	return c.userAgent
}

func (c *Client) Projects() *ProjectsResource {
	return &ProjectsResource{
		Client: *c,
	}
}
