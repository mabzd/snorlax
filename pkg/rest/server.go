package rest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mabzd/snorlax/api"
	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/internal/service"
)

func NewServerHandler(cfg config.Config) http.Handler {
	svc := service.NewSleepDiaryService(cfg)
	mux := http.NewServeMux()
	add(mux, "GET /sleep_diary/entries/{id}", getSleepDiaryEntry(svc))
	add(mux, "GET /sleep_diary/entries", getSleepDiaryEntries(svc))
	add(mux, "POST /sleep_diary/entries", createSleepDiaryEntry(svc))
	add(mux, "PUT /sleep_diary/entries/{id}", updateSleepDiaryEntry(svc))
	add(mux, "/", notFound())
	return mux
}

func NewServer(cfg config.Config) *http.Server {
	return &http.Server{
		Handler:      NewServerHandler(cfg),
		Addr:         fmt.Sprintf(":%s", cfg.ApiPort),
		WriteTimeout: time.Duration(cfg.ServerTimeoutInSec) * time.Second,
		ReadTimeout:  time.Duration(cfg.ServerTimeoutInSec) * time.Second,
	}
}

func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, api.ERR_NOT_FOUND, "path not found", nil)
	}
}

func add(mux *http.ServeMux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, withTrace(withLog(handler)))
}

func withTrace(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const TRACE_HEADER = "X-Trace-Id"
		traceId := r.Header.Get(TRACE_HEADER)
		if traceId == "" {
			traceId = uuid.NewString()
			r.Header.Set(TRACE_HEADER, traceId)
		}
		log.SetPrefix(fmt.Sprintf("[%s] ", traceId))
		next(w, r)
	}
}

func withLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		log.Printf("Completed '%s %s' in %v", r.Method, r.URL.Path, duration)
	}
}
