package veritashttp

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartServer(t *testing.T) {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt)

	go func() {
		require.NoError(t, StartServer(context.Background(), shutdownSignal))
	}()

	time.Sleep(2 * time.Second)

	t.Run("root", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/")
		require.NoError(t, err)

		defer resp.Body.Close()
		rb, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Welcome to Veritas 🤠\n", string(rb))
	})

	t.Run("GET scan success", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/scan?image=cgr.dev/chainguard/static:latest")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-GET scan failure", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/scan?image=cgr.dev/chainguard/static:latest", "application/json", nil)
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

	t.Run("shutdown", func(t *testing.T) {
		// send SIGINT signal to shutdown server
		shutdownSignal <- os.Interrupt

		// allow server time to shutdown
		time.Sleep(1 * time.Second)

		// ensure server has shutdown
		resp, err := http.Get("http://localhost:8080/")
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
