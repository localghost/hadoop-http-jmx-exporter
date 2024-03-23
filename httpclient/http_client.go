package httpclient

import (
	"net/http"

	"github.com/jcmturner/gokrb5/spnego"
)

type HttpClient interface {
    Get(url string) (*http.Response, error)
}

type HttpClientPure struct {
    *http.Client
}

type HttpClientWithSpnego struct {
    *spnego.Client
}

