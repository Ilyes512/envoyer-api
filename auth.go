package envoyerapi

import "net/http"

type Auth interface {
	Apply(req *http.Request)
}

type BearerToken string

func (b BearerToken) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+string(b))
}
