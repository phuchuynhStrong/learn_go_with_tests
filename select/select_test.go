package select_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Racer(a, b string) string {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration > bDuration {
		return b
	}
	return a
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

func makeDelayServer(d time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(d)
		w.WriteHeader(http.StatusOK)
	}))
}

func TestRacer(t *testing.T) {
	slowServer := makeDelayServer(20 * time.Millisecond)
	fastServer := makeDelayServer(0 * time.Millisecond)

	defer func() {
		slowServer.Close()
		fastServer.Close()
	}()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
