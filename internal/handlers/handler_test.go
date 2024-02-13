package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var testUser = int64(1)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, reqBody io.Reader) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	contentType := resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, string(respBody)
}
