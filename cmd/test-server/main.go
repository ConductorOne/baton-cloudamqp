// Command test-server is an in-process mock of the CloudAMQP customer API,
// covering the surface baton-cloudamqp uses: GET /team, POST /team/invite,
// POST /team/remove, and PUT /team/{id}. It exists so CI (and local runs) can
// exercise sync + account provisioning end-to-end without a real CloudAMQP
// tenant or credentials.
//
// Point the connector at it with: --base-url http://localhost:8080/api
// (or BATON_BASE_URL). Any non-empty Basic auth is accepted.
//
// Invite semantics mirror the real API (and the baton-workato mock it was modeled
// on): POST /team/invite is ASYNC — the invitee only becomes a team member after
// accepting the emailed invitation, so an invited address does NOT appear in
// /team. This exercises the connector's pending-invite (ActionRequiredResult)
// path. A duplicate invite (already a member, or already invited) returns 409
// Conflict, which the connector classifies as already-exists. Delete/role-update
// operate on the seeded members.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type user struct {
	Id    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type store struct {
	mu      sync.Mutex
	users   []user
	invited map[string]bool
}

func newStore() *store {
	return &store{
		users: []user{
			{Id: "1", Email: "owner@example.com", Roles: []string{"owner"}},
			{Id: "2", Email: "dev@example.com", Roles: []string{"developer"}},
		},
		invited: map[string]bool{},
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

// invite records a pending invitation. The invitee does NOT become a member
// until they accept (async), so this never adds to /team. Returns 409 if the
// email is already a member or already invited — the connector's
// IsAlreadyExistsError path.
func (s *store) invite(email string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.findByEmail(email); ok {
		return http.StatusConflict
	}
	key := strings.ToLower(email)
	if s.invited[key] {
		return http.StatusConflict
	}
	s.invited[key] = true
	return http.StatusOK
}

// remove deletes a member by email. Returns 404 if no such member (the connector
// treats not-found-on-delete as success). Also clears any pending invitation.
func (s *store) remove(email string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.invited, strings.ToLower(email))
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
		_ = r.ParseForm()
		if code := s.invite(r.PostForm.Get("email")); code != http.StatusOK {
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
		_ = r.ParseForm()
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
		_ = r.ParseForm()
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
