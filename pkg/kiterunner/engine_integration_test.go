//go:build integration
// +build integration

package kiterunner

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	krhttp "github.com/assetnote/kiterunner/pkg/http"
)

func TestEngineRunAgainstHTTPTestServer(t *testing.T) {
	var requests int32
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(stdhttp.StatusOK)
	}))
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(serverURL.Port())
	if err != nil {
		t.Fatal(err)
	}
	target := &krhttp.Target{IP: serverURL.Hostname(), Hostname: serverURL.Hostname(), Port: port}
	target.ParseHostHeader()
	target.SetContext(context.Background())
	target.HTTPClient(1, time.Second)

	config := NewDefaultConfig()
	config.RequestValidators = nil
	results := make(chan *Result, 1)
	job := &job{t: target, routes: []*krhttp.Route{{Path: []byte("/"), Method: krhttp.GET}}, client: target.HTTPClient(1, time.Second)}
	if err := handleRequest(context.Background(), job, results, config); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&requests) == 0 {
		t.Fatal("test server did not receive a request")
	}
}
