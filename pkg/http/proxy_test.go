package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClientWithProxy(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	defer target.Close()

	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			http.Error(w, "expected CONNECT", http.StatusMethodNotAllowed)
			return
		}
		client, upstream := w.(http.Hijacker)
		conn, _, err := client.Hijack()
		if err != nil {
			t.Errorf("hijack: %v", err)
			return
		}
		remote, err := net.Dial("tcp", r.Host)
		if err != nil {
			conn.Close()
			t.Errorf("dial target: %v", err)
			return
		}
		fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		go func() { _, _ = io.Copy(remote, bufio.NewReader(conn)); remote.Close() }()
		go func() { _, _ = io.Copy(conn, remote); conn.Close() }()
		_ = upstream
	}))
	defer proxy.Close()

	client := NewHTTPClientWithProxy(target.Listener.Addr().String(), false, proxy.URL)
	client.ReadTimeout = 2e9
	client.WriteTimeout = 2e9
	addr := target.Listener.Addr().(*net.TCPAddr)
	req := Request{Route: &Route{Path: []byte("/")}, Target: &Target{Hostname: addr.IP.String(), Port: addr.Port, HostHeader: []byte(addr.String())}}
	_, err := DoClient(client, req, &Config{})
	if err != nil {
		t.Fatalf("request through proxy failed: %v", err)
	}
}
