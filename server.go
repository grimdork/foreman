package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/grimdork/foreman/api"
	"github.com/grimdork/foreman/clients"
	ll "github.com/grimdork/loglines"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Server structure.
type Server struct {
	sync.RWMutex
	sync.WaitGroup
	http.Server
	r  *chi.Mux
	db *pgxpool.Pool

	// All running scouts and canaries.
	scouts   map[string]*clients.Scout
	canaries map[string]*clients.Canary

	log     chan string
	err     chan string
	logquit chan any
}

const (
	// WEBHOST default
	WEBHOST = "127.0.0.1"
	// WEBPORT default
	WEBPORT = "80"
)

// NewServer init. Reads settings from environment:
// WEB_HOST - default 127.0.0.1
// WEB_PORT - default 80
// FOREMAN_* - no defaults
func NewServer() (*Server, error) {
	srv := &Server{
		scouts:   make(map[string]*clients.Scout),
		canaries: make(map[string]*clients.Canary),
	}

	err := srv.openDB()
	if err != nil {
		return nil, err
	}

	// Default timeouts
	srv.Server.IdleTimeout = time.Second * 20
	srv.Server.ReadTimeout = time.Second * 20
	srv.Server.ReadHeaderTimeout = time.Second * 5
	srv.Server.WriteTimeout = time.Second * 20
	srv.r = chi.NewRouter()
	srv.r.Use(
		middleware.NoCache,
		middleware.Heartbeat(api.EPHealth),
		middleware.RealIP,
		middleware.RequestID,
		srv.addLogger,
	)

	srv.r.Post(api.EPPulse, srv.pulse)
	srv.r.Route("/api", func(r chi.Router) {
		r.Use(
			AddCORS,
			AddJSONHeaders,
			srv.authAdmin,
		)

		r.Get(api.EPScouts, srv.scoutsGet)
		r.Get(api.EPScout, srv.scoutGet)
		r.Post(api.EPScout, srv.scoutPost)
		r.Delete(api.EPScout, srv.scoutDelete)

		r.Get(api.EPCanaries, srv.canariesGet)
		r.Get(api.EPCanary, srv.canaryGet)
		r.Post(api.EPCanary, srv.canaryPost)
		r.Delete(api.EPCanary, srv.canaryDelete)

		r.Get(api.EPKeys, srv.keysGet)
		r.Post(api.EPKey, srv.keyPost)
		r.Delete(api.EPKey, srv.keyDelete)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(api.Endpoints))
		})
	})

	srv.r.Route("/", func(r chi.Router) {
		r.Use(
			AddHTMLHeaders,
		)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("."))
		})
	})

	return srv, nil
}

// Start serving based on environment variables.
func (srv *Server) Start() error {
	srv.Lock()
	defer srv.Unlock()

	addr := net.JoinHostPort(
		os.Getenv("WEB_HOST"),
		os.Getenv("WEB_PORT"),
	)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv.startLogger()
	srv.log <- fmt.Sprintf("Starting web server on http://%s", addr)
	srv.Add(1)
	go func() {
		srv.Handler = srv.r
		err = srv.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			ll.Err("Error running server: %s", err.Error())
			os.Exit(2)
		}
		srv.log <- "Stopping web server."
		srv.stopLogger()
		srv.Done()
	}()

	return nil
}

// Stop serving.
func (srv *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		ll.Err("Shutdown error: %s", err.Error())
	}

	srv.Wait()
	srv.closeDB()
}
