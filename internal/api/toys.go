package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Toy struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	InStock  bool    `json:"in_stock"`
}

type toyInput struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	InStock  *bool   `json:"in_stock"`
}

type toyStore struct {
	mu     sync.RWMutex
	toys   map[int]Toy
	nextID int
}

func newToyStore() *toyStore {
	store := &toyStore{
		toys:   make(map[int]Toy),
		nextID: 0,
	}

	store.bootstrap()
	return store
}

func (s *toyStore) bootstrap() {
	s.createToy(Toy{Name: "Rocket Racer", Category: "Vehicles", Price: 19.99, InStock: true})
	s.createToy(Toy{Name: "Galactic Puzzle", Category: "Puzzles", Price: 12.50, InStock: true})
	s.createToy(Toy{Name: "Coding Bot", Category: "STEM", Price: 34.75, InStock: false})
}

func (s *toyStore) listToys() []Toy {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]int, 0, len(s.toys))
	for id := range s.toys {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	toys := make([]Toy, 0, len(ids))
	for _, id := range ids {
		toys = append(toys, s.toys[id])
	}
	return toys
}

func (s *toyStore) getToy(id int) (Toy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	toy, ok := s.toys[id]
	return toy, ok
}

func (s *toyStore) createToy(input Toy) Toy {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	input.ID = s.nextID
	s.toys[input.ID] = input
	return input
}

func (s *toyStore) updateToy(id int, input Toy) (Toy, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.toys[id]; !ok {
		return Toy{}, false
	}

	input.ID = id
	s.toys[id] = input
	return input, true
}

func (s *toyStore) deleteToy(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.toys[id]; !ok {
		return false
	}

	delete(s.toys, id)
	return true
}

var toys = newToyStore()

func toysCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		respondJSON(w, http.StatusOK, toys.listToys())
	case http.MethodPost:
		toy, err := decodeToyInput(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		created := toys.createToy(toy)
		respondJSON(w, http.StatusCreated, created)
	default:
		respondMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func toysItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseToyID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		toy, ok := toys.getToy(id)
		if !ok {
			respondError(w, http.StatusNotFound, "toy not found")
			return
		}
		respondJSON(w, http.StatusOK, toy)
	case http.MethodPut:
		toy, err := decodeToyInput(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, ok := toys.updateToy(id, toy)
		if !ok {
			respondError(w, http.StatusNotFound, "toy not found")
			return
		}
		respondJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if ok := toys.deleteToy(id); !ok {
			respondError(w, http.StatusNotFound, "toy not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		respondMethodNotAllowed(w, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func decodeToyInput(r *http.Request) (Toy, error) {
	defer r.Body.Close()

	var payload toyInput
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return Toy{}, errors.New("invalid JSON payload")
	}

	if strings.TrimSpace(payload.Name) == "" {
		return Toy{}, errors.New("name is required")
	}
	if strings.TrimSpace(payload.Category) == "" {
		return Toy{}, errors.New("category is required")
	}
	if payload.InStock == nil {
		return Toy{}, errors.New("in_stock is required")
	}
	if payload.Price < 0 {
		return Toy{}, errors.New("price must be zero or positive")
	}

	return Toy{
		Name:     strings.TrimSpace(payload.Name),
		Category: strings.TrimSpace(payload.Category),
		Price:    payload.Price,
		InStock:  *payload.InStock,
	}, nil
}

func parseToyID(path string) (int, error) {
	const prefix = "/api/v1/toys/"
	if !strings.HasPrefix(path, prefix) {
		return 0, errors.New("missing toy id")
	}

	idPart := strings.TrimPrefix(path, prefix)
	if idPart == "" {
		return 0, errors.New("missing toy id")
	}
	if strings.Contains(idPart, "/") {
		return 0, errors.New("invalid toy id")
	}

	id, err := strconv.Atoi(idPart)
	if err != nil {
		return 0, errors.New("invalid toy id")
	}

	return id, nil
}

func respondMethodNotAllowed(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	respondError(w, http.StatusMethodNotAllowed, "method not allowed")
}
