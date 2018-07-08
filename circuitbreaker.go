package http

import (
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
)

type (
	CbTransport struct {
		Transport     http.RoundTripper
		transportName string
	}
)

func NewCbTransport(transport http.RoundTripper, transportName string, config *hystrix.CommandConfig) *(CbTransport) {
	hystrix.ConfigureCommand(transportName, *config)
	return &CbTransport{transport, transportName}
}

func (t *CbTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	output := make(chan *http.Response, 1)
	errors := hystrix.Go(t.transportName, func() error {
		resp, err := t.Transport.RoundTrip(req)
		if err != nil {
			return err
		}
		output <- resp

		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}
