package server

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestStart(t *testing.T) {
	var server *http.Server
	var err error
	var g errgroup.Group
	g.Go(func() error {
		server, err = StartServer("../../dashboard")
		require.NoError(t, err)
		return nil
	})

	t.Cleanup(func() {
		err := server.Shutdown(context.Background())
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://localhost:8080/")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("dashboard", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	})

	t.Run("GET scan success", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/scan/report?image=cgr.dev/chainguard/static:latest")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	})

	t.Run("non-GET scan failure", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/scan/report?image=cgr.dev/chainguard/static:latest", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		rb, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		require.Equal(t, "405 Method Not Allowed\n", string(rb))
	})

	t.Run("invalid path", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/invalid")
		require.NoError(t, err)
		defer resp.Body.Close()

		rb, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		require.Equal(t, "404 page not found\n", string(rb))
	})
}
