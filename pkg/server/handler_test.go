package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		userInput      string
	}{
		{
			name:           "valid input",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			userInput:      "cgr.dev/chainguard/static:latest",
		},
		{
			name:           "invalid input",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
			userInput:      "not-a-valid-image-reference........",
		},
		{
			name:           "bad request",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			userInput:      "",
		},
		{
			name:           "non-GET request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://example.com/", nil)
			queryParams := req.URL.Query()
			queryParams.Set("image", tt.userInput)
			req.URL.RawQuery = queryParams.Encode()

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(scanReportHandler)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestStripSlashes(t *testing.T) {
	type ctxKey string

	const expectedPath ctxKey = "expectedPath"

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "no trailing slash",
			in:   "/test",
			out:  "/test",
		},
		{
			name: "single trailing slash",
			in:   "/test/",
			out:  "/test",
		},
		{
			name: "double trailing slashes",
			in:   "/test//",
			out:  "/test/",
		},
		{
			name: "multiple slashes",
			in:   "/test//example/",
			out:  "/test//example",
		},
		{
			name: "root path",
			in:   "/",
			out:  "/",
		},
		{
			name: "empty path",
			in:   "",
			out:  "",
		},
	}

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Context().Value(expectedPath), r.URL.Path)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com"+tt.in, nil)
			req = req.WithContext(context.WithValue(req.Context(), expectedPath, tt.out))
			rr := httptest.NewRecorder()

			handler := stripSlashes(mockHandler)
			handler.ServeHTTP(rr, req)
		})
	}
}
