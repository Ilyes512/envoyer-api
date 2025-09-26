package envoyerapi

import "net/http"

const DefaultBaseUrl = "https://envoyer.io/api/"

type Client struct {
	httpClient *http.Client
	baseUrl    string
	apiToken   string
}

type options func(*Client)

func WithBaseUrl(baseUrl string) options {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

func NewClient(apiToken string, opts ...options) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}

	for _, opt := range opts {
		opt(client)
	}

	if client.baseUrl == "" {
		client.baseUrl = DefaultBaseUrl
	}

	return client
}

func (c *Client) HelloWorld() string {
	return "Hello, World from envoyer-api!"
}
