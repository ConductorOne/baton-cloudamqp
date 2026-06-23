// Command test-server is an in-process mock of the CloudAMQP customer API,
// covering the surface baton-cloudamqp uses: GET /team, POST /team/invite,
// POST /team/remove, and PUT /team/{id}. It exists so CI (and local runs) can
// exercise sync + account provisioning end-to-end without a real CloudAMQP
// tenant or credentials.
//
// Point the connector at it with: --base-url http://localhost:8080/api
// (or BATON_BASE_URL). Any non-empty Basic auth is accepted.
//
// Invite semantics: the real CloudAMQP API only materializes a team member once
// the invitee accepts the emailed invitation. To make the full
// create -> find -> delete lifecycle testable without a human in the loop (the
// account-provisioning CI action creates an account and then resolves it by
// email), this mock AUTO-ACCEPTS: POST /team/invite immediately adds the member
// to /team with a generated id. A duplicate invite (already a member) returns
// 409 Conflict, which the connector classifies as already-exists. That
// auto-accept is the one intentional divergence from production.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const maxFormBytes = 1 << 20

type user struct {
	Id    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type store struct {
	mu     sync.Mutex
	users  []user
	nextID int
}

func newStore() *store {
	return &store{
		users: []user{
			{Id: "1", Email: "owner@example.com", Roles: []string{"owner"}},
			{Id: "2", Email: "dev@example.com", Roles: []string{"developer"}},
		},
		nextID: 3,
	}
}

func (s *store) list() []user {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]user, len(s.users))
	copy(out, s.users)
	return out
}

// findByEmail must be called with s.mu held.
func (s *store) findByEmail(email string) (int, bool) {
	for i := range s.users {
		if strings.EqualFold(s.users[i].Email, email) {
			return i, true
		}
	}
	return 0, false
}

// invite auto-accepts: it immediately adds the invitee as a team member so the
// account-provisioning lifecycle (create -> resolve-by-email -> delete) can run
// end-to-end. Returns 409 if the email is already a member — the connector's
// IsAlreadyExistsError path.
func (s *store) invite(email, role string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.findByEmail(email); ok {
		return http.StatusConflict
	}
	if role == "" {
		role = "member"
	}
	s.users = append(s.users, user{
		Id:    fmt.Sprintf("%d", s.nextID),
		Email: email,
		Roles: []string{role},
	})
	s.nextID++
	return http.StatusOK
}

// remove deletes a member by email. Returns 404 if no such member (the connector
// treats not-found-on-delete as success).
func (s *store) remove(email string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.findByEmail(email)
	if !ok {
		return http.StatusNotFound
	}
	s.users = append(s.users[:i], s.users[i+1:]...)
	return http.StatusOK
}

func (s *store) updateRole(id, role string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.users {
		if s.users[i].Id == id {
			s.users[i].Roles = []string{role}
			return http.StatusOK
		}
	}
	return http.StatusNotFound
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
		http.Error(w, "missing Basic auth", http.StatusUnauthorized)
		return false
	}
	return true
}

// parseForm bounds the request body before parsing (gosec G120).
func parseForm(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
	_ = r.ParseForm()
}

func newMux(s *store) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check (unauthenticated) — the CI start-test-server action polls this.
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// GET /api/team — list members.
	mux.HandleFunc("/api/team", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s.list())
	})

	mux.HandleFunc("/api/team/invite", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parseForm(w, r)
		if code := s.invite(r.PostForm.Get("email"), r.PostForm.Get("role")); code != http.StatusOK {
			http.Error(w, "invite failed", code)
		}
	})

	mux.HandleFunc("/api/team/remove", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parseForm(w, r)
		if code := s.remove(r.PostForm.Get("email")); code != http.StatusOK {
			http.Error(w, "remove failed", code)
		}
	})

	// PUT /api/team/{id} — update role. Registered on the trailing-slash prefix;
	// ServeMux longest-prefix match routes "/api/team/<id>" here while the exact
	// "/api/team", "/api/team/invite" and "/api/team/remove" routes win for theirs.
	mux.HandleFunc("/api/team/", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/team/")
		if id == "" || id == "invite" || id == "remove" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parseForm(w, r)
		if code := s.updateRole(id, r.PostForm.Get("role")); code != http.StatusOK {
			http.Error(w, "update failed", code)
		}
	})

	return mux
}

func run() error {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           newMux(newStore()),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("cloudamqp test-server listening on http://%s/api", *addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("test-server error: %v", err)
		os.Exit(1)
	}
}
