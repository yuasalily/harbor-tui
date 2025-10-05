package docker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/containerd/errdefs"
)

func TestInfo_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/_ping"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		case strings.HasSuffix(r.URL.Path, "/version"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"Version":    "27.2.0",
				"Os":         "linux",
				"ApiVersion": "1.47",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	sdk := newSDKClient(t, srv)
	cli := &Client{cli: sdk}

	info, err := cli.Info(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Version != "27.2.0" || info.OS != "linux" || info.APIVersion != "1.47" {
		t.Fatalf("unexpected info: %#v", info)
	}
}

func TestInfo_PingError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_ping") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer srv.Close()

	sdk := newSDKClient(t, srv)
	cli := &Client{cli: sdk}
	_, err := cli.Info(context.Background())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errdefs.IsInternal(err) {
		t.Fatalf("exptected Internal Server Error, got %v", err)
	}
}
