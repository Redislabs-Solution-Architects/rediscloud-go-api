package rediscloud_api

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

/*
This file contains test utility functions needed for the cfn testing
*/

func clientFromTestServerV2(s *httptest.Server, apiKey string, secretKey string, opts ...Option) (*Client, error) {
	// return NewClientV2(LogRequests(true), BaseURL(s.URL), Auth(apiKey, secretKey), Transporter(s.Client().Transport), opts...)
	opts = append(opts, LogRequests(true))
	opts = append(opts, BaseURL(s.URL))
	opts = append(opts, Auth(apiKey, secretKey))
	opts = append(opts, Transporter(s.Client().Transport))
	return NewClientV2(opts...)
}

// infiniteTestServer will return a server which will return the
// last response forever
func infiniteTestServer(apiKey, secretKey string, mockedResponses ...endpointRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "go-http-client") {
			w.WriteHeader(504)
			return
		}
		if r.Header.Get("X-Api-Key") != apiKey {
			w.WriteHeader(502)
			return
		}
		if r.Header.Get("X-Api-Secret-Key") != secretKey {
			w.WriteHeader(503)
			return
		}
		mockedResponse := mockedResponses[0]
		if !mockedResponse.matches(r) {
			w.WriteHeader(501)
			return
		}

		response := mockedResponse.response()
		if len(mockedResponses) > 1 {
			mockedResponses = mockedResponses[1:]
		}
		w.WriteHeader(mockedResponse.status)
		_, _ = w.Write([]byte(response))
	}
}
