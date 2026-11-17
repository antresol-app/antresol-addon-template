package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antresol-app/antresol/addonkit"
)

func serve(t *testing.T, ready *addonkit.Ready, r *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	if ready == nil {
		ready = &addonkit.Ready{}
	}
	rec := httptest.NewRecorder()
	newHandler(addonkit.Env{}, nil, ready).ServeHTTP(rec, r)
	return rec
}

func TestHealth(t *testing.T) {
	ready := &addonkit.Ready{}
	for _, want := range []bool{false, true} {
		rec := serve(t, ready, httptest.NewRequest("GET", "/health", nil))
		var got map[string]bool
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != http.StatusOK || got["ready"] != want {
			t.Fatalf("health: %d %s, want ready %v", rec.Code, rec.Body, want)
		}
		ready.Set()
	}
}

func TestHelloNamesTheUser(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/hello", nil)
	// the core escapes every non-ASCII character of the header's JSON
	req.Header.Set("X-Antresol-User", `{"id":1,"name":"Мама","isAdmin":true,"isKid":false,"lang":"ru"}`)
	rec := serve(t, nil, req)
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != http.StatusOK || got["name"] != "Мама" {
		t.Fatalf("hello: %d %s", rec.Code, rec.Body)
	}
}

func TestHelloWithoutUser(t *testing.T) {
	if rec := serve(t, nil, httptest.NewRequest("GET", "/api/hello", nil)); rec.Code != http.StatusUnauthorized {
		t.Fatalf("no user: %d", rec.Code)
	}
}

func TestAdminOnly(t *testing.T) {
	for user, want := range map[string]int{
		`{"id":1,"name":"A","isAdmin":true}`: http.StatusOK,
		`{"id":2,"name":"B"}`:                http.StatusForbidden,
		``:                                   http.StatusForbidden,
	} {
		req := httptest.NewRequest("GET", "/api/admin/hello", nil)
		if user != "" {
			req.Header.Set("X-Antresol-User", user)
		}
		if rec := serve(t, nil, req); rec.Code != want {
			t.Fatalf("%q: %d, want %d", user, rec.Code, want)
		}
	}
}

func TestUndeclaredHookIs404(t *testing.T) {
	if rec := serve(t, nil, httptest.NewRequest("POST", "/hooks/library.pending", nil)); rec.Code != http.StatusNotFound {
		t.Fatalf("hook: %d", rec.Code)
	}
}
