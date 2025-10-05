package docker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/containerd/errdefs"
)

func TestImagesList_OK(t *testing.T) {
	created := time.Unix(1_720_000_000, 0).Unix()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/images/json"):
			_, _ = url.ParseQuery(r.URL.RawQuery)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{
					"Id":       "sha256:111",
					"RepoTags": []string{"alpine:latest"},
					"Size":     12345,
					"Created":  created,
				},
				{
					"Id":       "sha256:222",
					"RepoTags": []string{"busybox:1.36"},
					"Size":     67890,
					"Created":  created + 100,
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	sdk := newSDKClient(t, srv)
	cli := &Client{cli: sdk}

	imgs, err := cli.ImagesList(context.Background(), ImagesListOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(imgs) != 2 {
		t.Fatalf("expected 2 images, got %d", len(imgs))
	}
	if imgs[0].ID == "" || len(imgs[0].RepoTags) == 0 || imgs[0].Size == 0 {
		t.Fatalf("unexpected empty fields: %#v", imgs[0])
	}
}

func TestImagesList_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/images/json") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	sdk := newSDKClient(t, srv)
	cli := &Client{cli: sdk}
	_, err := cli.ImagesList(context.Background(), ImagesListOptions{All: true})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errdefs.IsInternal(err) {
		t.Fatalf("exptected Internal Server Error, got %v", err)
	}
}
