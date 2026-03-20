package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	srv := httptest.NewServer(handler)
	c := &Client{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	}
	return srv, c
}

func TestList(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("page_size") != "10" {
			t.Errorf("expected page_size=10, got %s", r.URL.Query().Get("page_size"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":          []map[string]interface{}{{"id": "abc", "name": "test"}},
			"next_page_token": "tok123",
		})
	})
	defer srv.Close()

	params := url.Values{}
	params.Set("page_size", "10")
	items, token, err := c.List(context.Background(), "/things", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0]["id"] != "abc" {
		t.Errorf("expected id=abc, got %v", items[0]["id"])
	}
	if token != "tok123" {
		t.Errorf("expected token=tok123, got %s", token)
	}
}

func TestListEmpty(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":          nil,
			"next_page_token": "",
		})
	})
	defer srv.Close()

	items, token, err := c.List(context.Background(), "/things", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected empty items, got %d", len(items))
	}
	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}
}

func TestGet(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/things/abc" {
			t.Errorf("expected /things/abc, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "abc", "name": "test"})
	})
	defer srv.Close()

	result, err := c.Get(context.Background(), "/things/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "abc" {
		t.Errorf("expected id=abc, got %v", result["id"])
	}
}

func TestPost(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "new" {
			t.Errorf("expected name=new, got %v", body["name"])
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "xyz", "name": "new"})
	})
	defer srv.Close()

	result, err := c.Post(context.Background(), "/things", map[string]interface{}{"name": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "xyz" {
		t.Errorf("expected id=xyz, got %v", result["id"])
	}
}

func TestPostNilBody(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "" {
			t.Errorf("expected no Content-Type for nil body, got %s", ct)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})
	defer srv.Close()

	result, err := c.Post(context.Background(), "/things/abc/action", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", result["status"])
	}
}

func TestPostEmptyResponse(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	defer srv.Close()

	result, err := c.Post(context.Background(), "/things/abc/action", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result on success")
	}
}

func TestPut(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "abc", "name": "updated"})
	})
	defer srv.Close()

	result, err := c.Put(context.Background(), "/things/abc", map[string]interface{}{"name": "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "updated" {
		t.Errorf("expected name=updated, got %v", result["name"])
	}
}

func TestDelete(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "abc", "deleted": true})
	})
	defer srv.Close()

	result, err := c.Delete(context.Background(), "/things/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "abc" {
		t.Errorf("expected id=abc, got %v", result["id"])
	}
}

func TestAPIError(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"message":"forbidden"}`))
	})
	defer srv.Close()

	_, err := c.Get(context.Background(), "/things/abc")
	if err == nil {
		t.Fatal("expected error on 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should contain status code, got: %s", err.Error())
	}
}

func TestAccessKeyTransport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("accesskey") != "mykey123" {
			t.Errorf("expected accesskey=mykey123, got %s", r.URL.Query().Get("accesskey"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "me"})
	}))
	defer srv.Close()

	c := New(srv.URL, "mykey123")
	// Override the HTTP client to use the test server's client but keep our transport
	transport := c.HTTPClient.Transport
	c.HTTPClient = srv.Client()
	c.HTTPClient.Transport = transport

	result, err := c.Get(context.Background(), "/me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "me" {
		t.Errorf("expected id=me, got %v", result["id"])
	}
}

func TestRawGet(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("binary data"))
	})
	defer srv.Close()

	resp, err := c.RawGet(context.Background(), "/file/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		t.Errorf("expected octet-stream, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestListAPIError(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"internal error"}`))
	})
	defer srv.Close()

	_, _, err := c.List(context.Background(), "/things", nil)
	if err == nil {
		t.Fatal("expected error on 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code, got: %s", err.Error())
	}
}
