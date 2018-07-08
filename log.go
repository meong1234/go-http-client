package http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	LoggedTransport struct {
		transport http.RoundTripper
		getLogger GetLoggerContextFunc
	}
)

func NewLoggedTransport(transport http.RoundTripper, getLogger GetLoggerContextFunc) *(LoggedTransport) {
	return &LoggedTransport{
		transport: transport,
		getLogger: getLogger,
	}
}

func (t *LoggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := t.getLogger(req.Context())

	t.logRequest(req, logger)

	resp, err := t.getTransport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp, logger)

	return resp, err
}

func (t *LoggedTransport) logRequest(req *http.Request, logger *logrus.Entry) {
	buf := GetBodyBytes(req)
	logger.WithField("URL", req.URL.String()).
		WithField("Header", fmt.Sprintf("%v", req.Header)).
		WithField("Body", string(buf)).
		Info("CLIENT-REQUEST")
}

func (t *LoggedTransport) logResponse(resp *http.Response, logger *logrus.Entry) {
	buf := GetResponseBodyBytes(resp)
	logger.WithField("Status", resp.StatusCode).
		WithField("Header", fmt.Sprintf("%v", resp.Header)).
		WithField("Body", string(buf)).
		Info("CLIENT-RESPONSE")
}

func (t *LoggedTransport) getTransport() http.RoundTripper {
	if t.transport != nil {
		return t.transport
	}

	return http.DefaultTransport
}
