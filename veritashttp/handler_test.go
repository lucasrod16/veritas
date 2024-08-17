package veritashttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anchore/grype/grype/db"
	"github.com/lucasrod16/veritas/vdb"
	"github.com/stretchr/testify/require"
)

func TestRootHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to Veritas 🤠\n",
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "405 Method Not Allowed\n",
		},
		{
			name:           "PUT request",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "405 Method Not Allowed\n",
		},
		{
			name:           "DELETE request",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "405 Method Not Allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://example.com/", nil)
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(rootHandler)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestScanHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
		cfg            db.Config
	}{
		{
			name:           "successful load",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "successfully loaded vulnerability database 🔐\n",
			cfg:            vdb.NewGrypeDBCfg(),
		},
		{
			name:           "failed load",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to load vulnerability database\n",
			cfg:            newTestDBCfg(),
		},
		{
			name:           "non-GET request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "405 Method Not Allowed\n",
			cfg:            vdb.NewGrypeDBCfg(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://example.com/", nil)
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(scanHandler(tt.cfg))
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedBody, rr.Body.String())
			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func newTestDBCfg() db.Config {
	return db.Config{
		DBRootDir:  "fake",
		ListingURL: "https://fake", // cause LoadVulnerabilityDB to fail
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
