package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type (
	GetLoggerContextFunc = func(context context.Context) *logrus.Entry

	GetTransactionContextFunc = func(context context.Context) newrelic.Transaction

	ClientConfig struct {
		ClientName                      string
		ServerCertificateFile           string
		HTTPClientTimeoutInMillis       int
		HTTPDialTimeoutInMillis         int
		HTTPTLSHandshakeTimeoutInMillis int
		ClientIdleConnTimeoutInMillis   int
		MaxIdleConns                    int
		MaxIdleConnsPerHost             int
		CbEnabled                       bool
		CbSetting                       *hystrix.CommandConfig
		getLogger                       GetLoggerContextFunc
		getTransaction                  GetTransactionContextFunc
	}
)

func (hc ClientConfig) Create() *http.Client {
	baseTransport := hc.getTransport(nil)

	if hc.ServerCertificateFile != "" {
		caCert, err := ioutil.ReadFile(hc.ServerCertificateFile)
		if err != nil {
			log.Panicf("can't read server certificate %v", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		baseTransport = hc.getTransport(tlsConfig)
	}

	var transport http.RoundTripper
	if hc.CbEnabled {
		cbTransport := NewCbTransport(baseTransport, hc.ClientName, hc.CbSetting)
		instrumentedTransport := NewIntrumentedTransport(cbTransport, hc.getTransaction)
		transport = NewLoggedTransport(instrumentedTransport, hc.getLogger)
	} else {
		instrumentedTransport := NewIntrumentedTransport(baseTransport, hc.getTransaction)
		transport = NewLoggedTransport(instrumentedTransport, hc.getLogger)
	}

	client := &http.Client{
		Timeout:   time.Duration(hc.HTTPClientTimeoutInMillis) * time.Millisecond,
		Transport: transport,
	}

	return client
}

func (hc ClientConfig) getTransport(tlsConfig *tls.Config) *http.Transport {
	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(hc.HTTPDialTimeoutInMillis) * time.Millisecond,
		}).Dial,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: time.Duration(hc.HTTPTLSHandshakeTimeoutInMillis) * time.Millisecond,
		MaxIdleConns:        hc.MaxIdleConns,
		MaxIdleConnsPerHost: hc.MaxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(hc.ClientIdleConnTimeoutInMillis) * time.Millisecond,
	}
}
