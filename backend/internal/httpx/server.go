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
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/admin"
	"phsar/internal/auth"
	"phsar/internal/catalogmeta"
	"phsar/internal/config"
	"phsar/internal/devtools"
	"phsar/internal/health"
	"phsar/internal/i18n"
	"phsar/internal/seller"
	"phsar/internal/storage"
	"phsar/internal/webui"
)

func NewRouter(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool, authSvc *auth.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(i18n.Middleware(i18n.Lang(cfg.DefaultLang)))

	if cfg.Env == "development" {
		r.Get("/dev/login", devtools.Login)
	}

	h := health.New(pool)
	r.Get("/healthz", h.Live)
	r.Get("/readyz", h.Ready)

	meta := catalogmeta.New(pool)
	r.Route("/v1/app", func(r chi.Router) {
		r.Get("/meta", meta.Meta)
		r.Get("/service-areas/check", meta.CheckPoint)
		r.Get("/categories", meta.Categories)
		r.Mount("/auth", authSvc.Routes("app"))
		r.Group(func(r chi.Router) {
			r.Use(authSvc.RequireUser("app"))
			r.Get("/me", authSvc.HandleMe)
		})
	})
	files := storage.Local{Dir: cfg.StorageDir, BaseURL: cfg.PublicBaseURL + "/media"}
	uploader := storage.Uploader{Store: files, Pool: pool}
	sellers := seller.NewService(pool, files)

	// File publik (logo toko, foto titik penjemputan). File private/… tidak pernah disajikan di sini.
	r.Handle("/media/public/*", http.StripPrefix("/media/public/",
		http.FileServer(noDirListing{http.Dir(cfg.StorageDir + "/public")})))

	r.Route("/v1/seller", func(r chi.Router) {
		r.Use(authSvc.RequireUser("app", "seller_web"))
		r.Mount("/", seller.NewHandler(sellers, uploader, log).Routes())
	})
	r.Route("/admin", func(r chi.Router) {
		// Login admin web lewat Telegram/SMS (MFA menyusul). Token app juga diterima sementara.
		r.Mount("/auth", authSvc.CookieRoutes("admin", "phsar_admin", "/admin/auth"))
		r.With(authSvc.RequireUser("admin", "app")).Mount("/api", admin.NewHandler(pool, sellers, files, log).Routes())
		r.Get("/", webui.AdminPage)
	})

	// Seller Center web
	r.Handle("/ui/*", webui.Assets())
	r.Route("/seller", func(r chi.Router) {
		r.Mount("/auth", authSvc.CookieRoutes("seller_web", "phsar_seller", "/seller/auth"))
		r.Get("/", webui.SellerPage)
	})
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/telegram/{key}", authSvc.TelegramWebhook)
		r.Post("/*", notYet)
	})

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		Error(w, http.StatusNotFound, "not_found", "route not found")
	})
	return r
}

// noDirListing mencegah daftar isi folder ditampilkan.
type noDirListing struct{ fs http.FileSystem }

func (n noDirListing) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if st, err := f.Stat(); err == nil && st.IsDir() {
		f.Close()
		return nil, os.ErrNotExist
	}
	return f, nil
}

func notYet(w http.ResponseWriter, _ *http.Request) {
	Error(w, http.StatusNotImplemented, "not_implemented", "coming in a later phase")
}

func Serve(ctx context.Context, cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) error {
	authSvc := auth.NewService(pool, cfg, log)
	if cfg.TelegramMode == "polling" {
		go authSvc.RunTelegramPolling(ctx)
	}
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           NewRouter(cfg, log, pool, authSvc),
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
