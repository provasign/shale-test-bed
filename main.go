// shale-test-bed: a deliberately small HTTP service used to demonstrate
// Shale (https://github.com/provasign/shale) on real pull requests.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/provasign/shale-test-bed/internal/auth"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", handleLogin)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !auth.Check(req.User, req.Password) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "user": req.User})
}
