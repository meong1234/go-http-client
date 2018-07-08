package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func GetBodyBytes(r *http.Request) []byte {
	if r.Body == nil {
		return make([]byte, 0)
	}
	requestBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
	return requestBytes
}

func GetResponseBodyBytes(r *http.Response) []byte {
	if r.Body == nil {
		return make([]byte, 0)
	}
	requestBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
	return requestBytes
}

func CopyHeader(source http.Header, target http.Header) {
	for header, values := range source {
		for _, value := range values {
			target.Set(header, value)
		}
	}
}
