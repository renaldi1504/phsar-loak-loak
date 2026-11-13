// Package httpx menyusun router HTTP: satu backend dengan beberapa "pintu".
//
//	/v1/app/*     app Flutter (pembeli + mode seller)
//	/v1/seller/*  API Seller Center (selalu dibatasi store_id)
//	/admin/*      admin Phsar (host terpisah + Cloudflare Access + MFA)
//	/webhooks/*   callback PayWay dll.
//	/             website customer (SSR, SEO) — Tahap 2
package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/catalogmeta"
	"phsar/internal/config"
	"phsar/internal/health"
	"phsar/internal/i18n"
)

func NewRouter(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(i18n.Middleware(i18n.Lang(cfg.DefaultLang)))

	h := health.New(pool)
	r.Get("/healthz", h.Live)
	r.Get("/readyz", h.Ready)

	meta := catalogmeta.New(pool)
	r.Route("/v1/app", func(r chi.Router) {
		r.Get("/meta", meta.Meta)
		r.Get("/service-areas/check", meta.CheckPoint)
		r.Get("/categories", meta.Categories)
	})
	r.Route("/v1/seller", func(r chi.Router) {
		r.Get("/*", notYet)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Get("/*", notYet)
	})
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/*", notYet)
	})

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		Error(w, http.StatusNotFound, "not_found", "route not found")
	})
	return r
}

func notYet(w http.ResponseWriter, _ *http.Request) {
	Error(w, http.StatusNotImplemented, "not_implemented", "coming in a later phase")
}

func Serve(ctx context.Context, cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) error {
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           NewRouter(cfg, log, pool),
		ReadHeaderTimeout: 10 * time.Second,
	}
	errc := make(chan error, 1)
	go func() {
		log.Info("http server listening", "addr", cfg.HTTPAddr, "env", cfg.Env)
		errc <- srv.ListenAndServe()
	}()
	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		log.Info("shutting down")
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	}
}

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("request",
				"method", r.Method, "path", r.URL.Path, "status", ww.Status(),
				"ms", time.Since(start).Milliseconds(), "req_id", middleware.GetReqID(r.Context()))
		})
	}
}
