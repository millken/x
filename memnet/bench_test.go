package memnet_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/millken/x/memnet"
)

func BenchmarkMemnetHTTP(b *testing.B) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, world.")
		}))
	benchmark(b, mux)
}

func benchmark(b *testing.B, h http.Handler) {
	ln, err := memnet.Listen("memu", "MyNamedNetwork")
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
	// serverStopCh := make(chan struct{})
	go func() {
		// serverLn := net.Listener(ln)
		if err := http.Serve(ln, h); err != http.ErrServerClosed {
			b.Errorf("unexpected error in server: %v", err)
		}
		// close(serverStopCh)
	}()
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				return memnet.DialContext(context.Background(), "memu", "MyNamedNetwork")
			},
			DialContext: func(
				ctx context.Context, _, _ string) (net.Conn, error) {
				return memnet.DialContext(ctx, "memu", "MyNamedNetwork")
			},
			DisableKeepAlives:  true,
			DisableCompression: true,
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		runRequests(b, pb, client)
	})
	// ln.Close()
	// <-serverStopCh
}

func runRequests(b *testing.B, pb *testing.PB, c *http.Client) {
	req, _ := http.NewRequest("GET", "http://foo.bar/baz", nil)
	var resp *http.Response
	var err error
	for pb.Next() {
		if resp, err = c.Do(req); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected status code: %d. Expecting %d", resp.StatusCode, http.StatusOK)
		}
	}
}
