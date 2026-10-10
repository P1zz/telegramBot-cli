package cmd

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (fn roundTripper) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestDownloadFileKeepsTokenOutOfOutput(t *testing.T) {
	previous := downloadHTTPClient
	downloadHTTPClient = &http.Client{Transport: roundTripper(func(request *http.Request) (*http.Response, error) {
		if !strings.Contains(request.URL.Path, "secret-token") {
			t.Fatalf("download URL missing token internally: %s", request.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Body: io.NopCloser(strings.NewReader("payload"))}, nil
	})}
	t.Cleanup(func() { downloadHTTPClient = previous })
	var output bytes.Buffer
	if err := downloadFile(context.Background(), "secret-token", "documents/file.txt", "-", &output); err != nil {
		t.Fatal(err)
	}
	if output.String() != "payload" || strings.Contains(output.String(), "secret-token") {
		t.Fatalf("unsafe output %q", output.String())
	}
}
