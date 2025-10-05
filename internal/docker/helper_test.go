package docker

import (
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/client"
)

func newSDKClient(t *testing.T, srv *httptest.Server) *client.Client {
	t.Helper()
	c, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(srv.Client()),
		client.WithVersion("1.43"),
	)
	if err != nil {
		t.Fatalf("failed to create sdk client: %v", err)
	}
	return c
}
