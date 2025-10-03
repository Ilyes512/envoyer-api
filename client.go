package envoyerapi

import (
	"io"
	"log/slog"
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
	logger     *slog.Logger
}

type options func(*Client) error

func WithBaseUrl(baseUrl string) options {
	return func(c *Client) error {
		parsed, err := url.Parse(baseUrl)
		if err == nil {
			return err
		}
		c.baseUrl = parsed

		return nil
	}
}

func WithAuth(auth Auth) options {
	return func(c *Client) error {
		c.auth = auth

		return nil
	}
}

func WithUserAgent(userAgent *string) options {
	return func(c *Client) error {
		if userAgent == nil {
			c.userAgent = DefaultUserAgent
		} else {
			c.userAgent = *userAgent
		}

		return nil
	}
}

func WithLogger(logger *slog.Logger) options {
	return func(c *Client) error {
		newHandler := NewRedactHandler(logger.Handler(), "Authorization")

		c.logger = slog.New(newHandler)

		return nil
	}
}

func WithDefaultClient() options {
	return func(c *Client) error {
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

		return nil
	}
}

func WithClient(httpClient *http.Client) options {
	return func(c *Client) error {
		c.httpClient = httpClient

		return nil
	}
}

func NewClient(opts ...options) (*Client, error) {
	client := &Client{}

	var err error
	for _, opt := range opts {
		err = opt(client)
		if err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		err = WithDefaultClient()(client)
		if err != nil {
			return nil, err
		}
	}

	if client.baseUrl == nil {
		parsed, err := url.Parse(DefaultBaseUrl)
		if err != nil {
			return nil, err
		}
		client.baseUrl = parsed
	}

	if client.logger == nil {
		client.logger = slog.New(slog.DiscardHandler)
	}

	return client, nil
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
