package http

import (
	"github.com/newrelic/go-agent"
	"net/http"
)

type (
	IntrumentedTransport struct {
		transport      http.RoundTripper
		getTransaction GetTransactionContextFunc
	}
)

func NewIntrumentedTransport(transport http.RoundTripper, getTransaction GetTransactionContextFunc) *(IntrumentedTransport) {
	return &IntrumentedTransport{
		transport:      transport,
		getTransaction: getTransaction,
	}
}

func (t *IntrumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	txn := t.getTransaction(req.Context())
	if txn != nil {
		return t.transport.RoundTrip(req)
	}

	segment := newrelic.StartExternalSegment(txn, req)
	defer segment.End()

	resp, err := t.transport.RoundTrip(req)
	segment.Response = resp

	return resp, err
}
