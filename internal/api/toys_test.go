package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestToyStore(items []Toy) *toyStore {
	store := &toyStore{toys: make(map[int]Toy)}
	for _, toy := range items {
		store.toys[toy.ID] = toy
		if toy.ID > store.nextID {
			store.nextID = toy.ID
		}
	}
	return store
}

func swapToyStore(t *testing.T, store *toyStore) {
	t.Helper()
	original := toys
	toys = store
	t.Cleanup(func() {
		toys = original
	})
}

func TestToysCollectionHandler_PostCreatesToy(t *testing.T) {
	store := newTestToyStore(nil)
	swapToyStore(t, store)

	payload := `{"name":"Test Rocket","category":"STEM","price":15.99,"in_stock":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/toys", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	toysCollectionHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var body Toy
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}

	if body.ID != 1 {
		t.Fatalf("expected created toy to have ID 1, got %d", body.ID)
	}
	if body.Name != "Test Rocket" {
		t.Fatalf("expected toy name 'Test Rocket', got %q", body.Name)
	}

	if _, ok := store.getToy(body.ID); !ok {
		t.Fatalf("expected toy to be persisted in store")
	}
}

func TestToysCollectionHandler_PostValidationError(t *testing.T) {
	store := newTestToyStore(nil)
	swapToyStore(t, store)

	payload := `{"name":"","category":"STEM","price":12.0}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/toys", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	toysCollectionHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var body errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}

	if body.Error == "" {
		t.Fatalf("expected error message to be populated")
	}
}

func TestToysItemHandler_GetNotFound(t *testing.T) {
	store := newTestToyStore(nil)
	swapToyStore(t, store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/toys/42", nil)
	rr := httptest.NewRecorder()

	toysItemHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestParseToyID(t *testing.T) {
	tests := map[string]struct {
		path string
		want int
		err  bool
	}{
		"valid":   {path: "/api/v1/toys/10", want: 10},
		"missing": {path: "/api/v1/toys", err: true},
		"invalid": {path: "/api/v1/toys/foo", err: true},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got, err := parseToyID(tc.path)
			if tc.err && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.err && got != tc.want {
				t.Fatalf("expected id %d, got %d", tc.want, got)
			}
		})
	}
}
