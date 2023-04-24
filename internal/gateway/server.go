package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/polyxia-org/gateway/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	HEALTH_ENDPOINT        = "/healthz"
	SKILLS_ENDPOINT        = "/v1/skills"
	DEVICE_DEMAND_ENDPOINT = "/v1/nlu"
)

type (
	Server struct {
		cfg *config.Config
	}

	APIError struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}

	APISucess struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}
)

func NewServer() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	server := &Server{
		cfg: cfg,
	}

	return server, nil
}

func (s *Server) Serve() {
	p := s.cfg.Port
	h := s.cfg.Addr

	ctx, stop := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h, p),
		Handler: s.router(),
	}

	// Listen for syscall signals for process to interrupt/quit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				cancel()
				log.Fatal("graceful shutdown timed out... forcing exit")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		stop()
	}()

	log.Printf("polyxia gateway is listening on %s:%d\n", h, p)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-ctx.Done()
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)

	r.Post(SKILLS_ENDPOINT, s.SkillsHandler)
	r.Post(DEVICE_DEMAND_ENDPOINT, s.DeviceDemandHandler)
	r.Get(HEALTH_ENDPOINT, s.HealthcheckHandler)

	return r
}

func (s *Server) JSONResponse(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) APIErrorResponse(w http.ResponseWriter, err *APIError) {
	log.Error(err.Message)
	s.JSONResponse(w, err.StatusCode, err)
}
