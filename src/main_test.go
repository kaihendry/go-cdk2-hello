package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead} {
		for _, path := range []string{"/", "/nested/path?test=1"} {
			t.Run(method+path, func(t *testing.T) {
				req, err := http.NewRequest(method, server.URL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				res, err := server.Client().Do(req)
				if err != nil {
					t.Fatal(err)
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if res.StatusCode != http.StatusOK || res.Header.Get("X-Version") != version {
					t.Fatalf("status=%d, version=%q", res.StatusCode, res.Header.Get("X-Version"))
				}
				want := "CICD World ... " + version
				if method == http.MethodHead {
					want = ""
				}
				if string(body) != want {
					t.Fatalf("body=%q, want %q", body, want)
				}
			})
		}
	}
}
