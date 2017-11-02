package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGetQuerySnippetGET(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	params := make(url.Values)
	q := "SELECT column FROM table"
	params.Set("query", q)
	req.URL.RawQuery = params.Encode()
	query := getQuerySnippet(req)
	if query != q {
		t.Fatalf("got: %q; expected: %q", query, q)
	}
}

func TestGetQuerySnippetPOST(t *testing.T) {
	q := "SELECT column FROM table"
	body := bytes.NewBufferString(q)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		panic(fmt.Sprintf("BUG: unexpected error: %s", err))
	}
	query := getQuerySnippet(req)
	if query != q {
		t.Fatalf("got: %q; expected: %q", query, q)
	}
}

func TestGetQuerySnippetGzipped(t *testing.T) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	q := makeQuery(1000)
	_, err := zw.Write([]byte(q))
	if err != nil {
		t.Fatal(err)
	}
	zw.Close()
	req, err := http.NewRequest("POST", "http://127.0.0.1:9090", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	if err != nil {
		t.Fatal(err)
	}
	query := getQuerySnippet(req)
	if query[:100] != string(q[:100]) {
		t.Fatalf("got: %q; expected: %q", query[:100], q[:100])
	}
}

func makeQuery(n int) []byte {
	q1 := "SELECT column "
	q2 := "WHERE Date=today()"

	var b []byte
	b = append(b, q1...)
	for i := 0; i < n; i++ {
		b = append(b, fmt.Sprintf("col%d, ", i)...)
	}
	b = append(b, q2...)
	return b
}