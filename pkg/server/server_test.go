package server

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	server, err := Start()
	require.NoError(t, err)

	// give the server some time to come up
	time.Sleep(1 * time.Second)

	t.Cleanup(func() {
		err := server.Shutdown(context.Background())
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://localhost:8080/")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	tests := []struct {
		name                string
		method              string
		url                 string
		expectedStatusCode  int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "GET /scan/report success",
			method:              http.MethodGet,
			url:                 "http://localhost:8080/scan/report?image=cgr.dev/chainguard/static:latest",
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
		{
			name:               "/scan/report method not allowed",
			method:             http.MethodPost,
			url:                "http://localhost:8080/scan/report?image=cgr.dev/chainguard/static:latest",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "405 Method Not Allowed\n",
		},
		{
			name:               "/scan/report internal server error",
			method:             http.MethodGet,
			url:                "http://localhost:8080/scan/report?image=fake",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "unable to scan \"fake\". please ensure the image exists, is publicly accessible, and does not require authentication\n",
		},
		{
			name:               "/scan/report bad request",
			method:             http.MethodGet,
			url:                "http://localhost:8080/scan/report?image=",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Missing 'image' query parameter\n",
		},
		{
			name:                "GET /scan/details success",
			method:              http.MethodGet,
			url:                 "http://localhost:8080/scan/details?image=cgr.dev/chainguard/static:latest",
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
		{
			name:               "/scan/details method not allowed",
			method:             http.MethodPost,
			url:                "http://localhost:8080/scan/details?image=cgr.dev/chainguard/static:latest",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "405 Method Not Allowed\n",
		},
		{
			name:               "/scan/details internal server error",
			method:             http.MethodGet,
			url:                "http://localhost:8080/scan/details?image=fake",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "unable to scan \"fake\". please ensure the image exists, is publicly accessible, and does not require authentication\n",
		},
		{
			name:               "/scan/details bad request",
			method:             http.MethodGet,
			url:                "http://localhost:8080/scan/details?image=",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Missing 'image' query parameter\n",
		},
		{
			name:               "invalid path",
			method:             http.MethodGet,
			url:                "http://localhost:8080/invalid",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(tt.method, tt.url, nil)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			rb, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedContentType != "" {
				require.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))
			}

			if tt.expectedBody != "" {
				require.Equal(t, tt.expectedBody, string(rb))
			}
		})
	}
}
