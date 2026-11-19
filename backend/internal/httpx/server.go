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
	"phsar/internal/ads"
	"phsar/internal/chat"
	"phsar/internal/content"
	"phsar/internal/flash"
	"phsar/internal/invoices"
	"phsar/internal/legal"
	"phsar/internal/reports"
	"phsar/internal/returns"
	"phsar/internal/tickets"
	"phsar/internal/tour"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/admin"
	"phsar/internal/auth"
	"phsar/internal/catalogmeta"
	"phsar/internal/config"
	"phsar/internal/courierreturn"
	"phsar/internal/delivery"
	"phsar/internal/demo"
	"phsar/internal/devtools"
	"phsar/internal/health"
	"phsar/internal/i18n"
	"phsar/internal/integ"
	"phsar/internal/jobs"
	"phsar/internal/media"
	"phsar/internal/notify"
	"phsar/internal/paygw"
	"phsar/internal/payment"
	"phsar/internal/payout"
	"phsar/internal/photosearch"
	"phsar/internal/secretbox"
	"phsar/internal/seller"
	"phsar/internal/shop"
	"phsar/internal/storage"
	"phsar/internal/storefront"
	"phsar/internal/translate"
	"phsar/internal/webui"
)

func NewRouter(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool, authSvc *auth.Service) http.Handler {
	r := chi.NewRouter()
	// Notifikasi web: satu hub realtime (LISTEN) + satu pengirim Web Push untuk semua audience.
	notifyHub, notifyPush := notify.NewHub(pool, log), NewPusher(cfg, pool)
	notifyFor := func(aud string) *notify.Handler {
		return &notify.Handler{Pool: pool, Log: log, Audience: aud, Hub: notifyHub, Push: notifyPush}
	}
	// Halaman error HTML (storefront); diisi setelah storefront dibuat di bawah.
	var htmlError func(http.ResponseWriter, *http.Request, int)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.GetHead) // HEAD (monitor uptime, crawler) dijawab seperti GET
	r.Use(requestLogger(log))
	r.Use(recoverer(log, func(w http.ResponseWriter, r *http.Request, code int) {
		if htmlError == nil {
			http.Error(w, http.StatusText(code), code)
			return
		}
		htmlError(w, r, code)
	}))
	r.Use(timeoutExceptUploads(30 * time.Second))
	r.Use(i18n.Middleware(i18n.Lang(cfg.DefaultLang)))

	if cfg.Env == "development" {
		r.Get("/dev/login", devtools.Login)
	}

	h := health.New(pool)
	r.Get("/healthz", h.Live)
	r.Get("/readyz", h.Ready)

	// Provider SMS & Login Facebook — diatur super admin (admin → Login & SMS), kredensial terenkripsi.
	contentStore := &content.Store{Pool: pool, Files: newFiles(cfg)} // content: dipakai bersama admin & website (cache 30 detik)
	legalStore := &legal.Store{Pool: pool}
	if err := legalStore.EnsureDefaults(context.Background()); err != nil {
		log.Warn("legal: dokumen bawaan gagal dipasang", "err", err)
	}
	logins := NewIntegrations(cfg, log, pool)
	authSvc.UseSMSGateway(logins)
	authSvc.UseFacebook(logins)

	meta := catalogmeta.New(pool)
	tourH := &tour.Handler{Pool: pool, Log: log} // tour: status "tur cara pakai sudah dilihat"
	r.Route("/v1/app", func(r chi.Router) {
		r.Get("/meta", meta.Meta)
		r.Get("/service-areas/check", meta.CheckPoint)
		r.Get("/categories", meta.Categories)
		r.Mount("/auth", authSvc.Routes("app"))
		r.Group(func(r chi.Router) {
			r.Use(authSvc.RequireUser("app"))
			r.Get("/me", authSvc.HandleMe)
			r.With(auth.NoGuest).Mount("/tours", tourH.BuyerRoutes()) // tour:
		})
	})
	files := newFiles(cfg)
	uploader := storage.Uploader{Store: files, Pool: pool}
	if media.Available() {
		uploader.CheckVideo = media.CheckVideo(files) // tolak video > 3 menit / rusak langsung saat upload
	}
	sellers := seller.NewService(pool, files)
	if tr := newTranslator(cfg, AIKeys(cfg, NewIntegrations(cfg, log, pool)), nil); tr == nil {
		sellers.TranslationOn = false
	} else if d, ok := tr.(translate.Dynamic); ok {
		sellers.TranslationReady = d.Ready // kunci di Admin → Integrasi
	}
	sellers.MediaOn = media.Available()

	// File publik (logo toko, foto titik penjemputan). File private/… tidak pernah disajikan di sini.
	r.Handle("/media/public/*", http.StripPrefix("/media/public/",
		http.FileServer(noDirListing{http.Dir(cfg.StorageDir + "/public")})))

	shopSvc := NewShop(cfg, log, pool)
	pay := shopSvc.Pay.(*paygw.Store)
	shopH := &shop.Handler{S: shopSvc, Log: log}
	reportSvc := &reports.Service{Pool: pool, Log: log} // reports:
	rh := &returns.Handler{S: &returns.Service{Pool: pool, Files: files, Log: log}, Uploader: uploader}
	th := &tickets.Handler{S: &tickets.Service{Pool: pool, Log: log}}                             // support: tiket bantuan
	ch := &chat.Handler{S: &chat.Service{Pool: pool, Files: files, Log: log}, Uploader: uploader} // chat: pembeli ↔ toko
	shopH.Extra = func(r chi.Router) { rh.BuyerRoutes(r); th.BuyerRoutes(r); ch.BuyerRoutes(r) }
	adsH := &ads.Handler{S: &ads.Service{Pool: pool, Log: log}} // ads: iklan bersponsor
	r.Get("/ad/{id}", adsH.Click)
	askbotSvc := NewAskbot(cfg, log, pool) // askbot: draf jawaban untuk seller + FAQ toko
	// Saldo & payout seller
	ph := &payout.Handler{S: NewPayouts(cfg, log, pool), Log: log}
	registerJobs(cfg, log, pool, shopSvc, ph.S, rh.S) // jobs: tombol "Run now" di admin
	// Kurir Phsar: web HP di /courier, API di /v1/courier
	dh := &delivery.Handler{S: &delivery.Service{Pool: pool, Files: files, Uploader: uploader, Log: log}, Log: log}
	// courierreturn: gagal antar → kembalikan ke toko (kurir, seller, admin) + job terima otomatis
	crh := &courierreturn.Handler{S: &courierreturn.Service{Pool: pool, Log: log}, Files: files, Log: log}
	jobs.SetRunner(courierreturn.JobName, crh.S.AutoReceive)
	r.With(authSvc.RequireUser("app"), auth.NoGuest).Mount("/v1/courier/returns", crh.CourierRoutes())
	r.With(authSvc.RequireUser("app"), auth.NoGuest).Mount("/v1/courier", dh.CourierRoutes())
	r.Route("/courier", func(r chi.Router) {
		r.Mount("/auth", authSvc.CookieRoutes("app", "phsar_courier", "/courier/auth"))
		r.Get("/", webui.CourierPage)
	})
	r.Route("/v1/seller", func(r chi.Router) {
		r.Use(authSvc.RequireUser("app", "seller_web"), auth.NoGuest)
		r.Mount("/notifications", notifyFor(notify.AudSeller).Routes())
		r.Mount("/tours", tourH.SellerRoutes()) // tour:
		sh := seller.NewHandler(sellers, uploader, log)
		sh.Orders = shopH.SellerRoutes()
		// invoices: faktur komisi untuk pemilik toko di /payouts/invoices (ikut pemeriksaan pemilik).
		payoutsR := chi.NewRouter()
		payoutsR.Mount("/invoices", (&invoices.Handler{S: &invoices.Service{Pool: pool}, Log: log}).SellerRoutes())
		payoutsR.Mount("/", ph.SellerRoutes())
		sh.Payouts = payoutsR
		sh.Returns = rh.SellerRoutes()
		sh.Strikes = reportSvc.SellerStrikes // reports:
		sh.Tickets = th.SellerRoutes()       // support
		sh.Chat = ch.SellerRoutes()          // chat:
		sh.Ads = adsH.SellerRoutes()
		sh.Flash = (&flash.Service{Pool: pool, Files: files, Log: log}).SellerRoutes()
		sh.Askbot = askbotSvc.SellerRoutes()   // askbot:
		sh.CourierReturns = crh.SellerRoutes() // courierreturn
		r.Mount("/", sh.Routes())
	})
	r.Route("/admin", func(r chi.Router) {
		// Login admin: HP (Telegram/SMS) → kode aplikasi authenticator (MFA). API hanya untuk sesi yang lolos MFA.
		r.Mount("/auth", authSvc.CookieRoutes("admin", "phsar_admin", "/admin/auth"))
		ah := admin.NewHandler(pool, sellers, files, log)
		ah.Refunds = &admin.RefundRoutes{List: shopH.AdminListRefunds, Done: shopH.AdminMarkRefunded}
		ah.Delivery = &admin.DeliveryRoutes{Couriers: dh.AdminCouriers, SaveCourier: dh.AdminSaveCourier, List: dh.AdminDeliveries,
			Assign: dh.AdminAssign, Proof: dh.AdminProof}
		ah.Payouts = &admin.PayoutRoutes{GetSettings: ph.AdminGetSettings, SaveSettings: ph.AdminSaveSettings, List: ph.AdminList,
			Generate: ph.AdminGenerate, Paid: ph.AdminPaid, Cancel: ph.AdminCancel}
		confirm := func(r *http.Request, userID, code string) error {
			return authSvc.ConfirmTOTP(r.Context(), userID, code, auth.ClientFrom(r))
		}
		ah.Payments = &admin.PaymentRoutes{Store: pay, ConfirmMFA: confirm}
		ah.Integrations = &admin.IntegrationRoutes{Store: logins, ConfirmMFA: confirm, TestAI: AITester(cfg)}
		ah.Ads = &admin.AdsRoutes{Handler: adsH.AdminRoutes()}
		ah.CourierReturns = crh.AdminRoutes() // courierreturn
		ah.Returns = &admin.ReturnRoutes{List: rh.AdminList, Get: rh.AdminGet, Decide: rh.AdminDecide, Refunded: rh.AdminRefunded, Photo: rh.AdminPhoto}
		ah.Content = &admin.ContentRoutes{Store: contentStore}
		ah.Legal = &admin.LegalRoutes{Store: legalStore, BaseURL: cfg.PublicBaseURL}
		ah.Askbot = askbotSvc.AdminRoutes(cfg.AnthropicModel) // askbot:
		// Notifikasi admin (dibuat per izin peran oleh trigger database) — sebelum /api supaya rutenya lebih spesifik.
		r.With(authSvc.RequireUser("admin"), authSvc.RequireMFA).Mount("/api/notifications", notifyFor(notify.AudAdmin).Routes())
		r.With(authSvc.RequireUser("admin"), authSvc.RequireMFA).Mount("/api", ah.Routes())
		r.Get("/", webui.AdminPage)
		r.Get("/sw.js", webui.ServiceWorker) // Web Push untuk /admin/
	})

	// Seller Center web
	r.Get("/ui/config.js", webui.ConfigJS(logins.Maps, func(ctx context.Context) map[string]bool {
		return map[string]bool{"telegram": cfg.TelegramMode != "off", "sms": authSvc.SMSEnabled(ctx), "facebook": authSvc.FacebookEnabled(ctx)}
	}))
	r.Handle("/ui/*", webui.Assets())
	r.Route("/seller", func(r chi.Router) {
		r.Mount("/auth", authSvc.CookieRoutes("seller_web", "phsar_seller", "/seller/auth"))
		r.Get("/", webui.SellerPage)
		r.Get("/sw.js", webui.ServiceWorker) // Web Push untuk /seller/
	})
	// Keranjang, checkout, pembayaran KHQR, pesanan pembeli
	r.With(authSvc.RequireUser("app")).Mount("/v1/shop/notifications", notifyFor(notify.AudBuyer).Routes()) // pembeli (tamu juga)
	r.Mount("/v1/shop", shopH.Routes(authSvc.RequireUser("app")))
	r.Get("/sw.js", webui.ServiceWorker) // Web Push untuk website pembeli
	// reports: pembeli melaporkan produk / toko
	r.With(authSvc.RequireUser("app"), auth.NoGuest).Post("/v1/reports", reportSvc.BuyerReport)
	r.Mount("/shop/auth", authSvc.CookieRoutes("app", "phsar_web", "/shop/auth"))
	if cfg.Env == "development" { // hanya bekerja saat mode pembayaran = fake
		r.With(authSvc.RequireUser("app")).Post("/dev/pay/{id}", shopH.DevPay)
	}

	// Website pembeli (SSR + SEO): /, /km/…, /en/…, /zh/…, /sitemap.xml, /robots.txt
	sf := storefront.New(&storefront.Catalog{Pool: pool, Files: files}, cfg.PublicBaseURL, cfg.DefaultLang, log)
	htmlError = sf.ErrorPage
	sf.TelegramLogin = cfg.TelegramBotToken != ""
	sf.Payments = pay   // pita "segera dibuka" & tombol simulasi bayar mengikuti pengaturan admin
	sf.Logins = authSvc // tombol SMS / Facebook mengikuti pengaturan admin
	sf.Photo = NewPhotoSearch(cfg, log, pool)
	sf.Content = contentStore
	sf.Legal = legalStore
	sf.Ads = adsH.S // ads: slot bersponsor di hasil cari & kategori
	// askbot: blok "Tanya jawab" di halaman produk
	sf.FAQ = askbotFAQ{askbotSvc}
	sf.Routes(r)

	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/telegram/{key}", authSvc.TelegramWebhook)
		r.Post("/payway", shopH.PayWayCallback)
		r.Post("/*", notYet)
	})

	// Browser (navigasi halaman) mendapat halaman error Phsar; API & fetch tetap JSON.
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		if wantsHTML(req) {
			sf.ErrorPage(w, req, http.StatusNotFound)
			return
		}
		Error(w, http.StatusNotFound, "not_found", "route not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		if wantsHTML(req) {
			sf.ErrorPage(w, req, http.StatusNotFound) // bagi pembeli: halaman ini tidak ada
			return
		}
		Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
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

// timeoutExceptUploads: batas waktu 30 detik untuk semua request kecuali upload video/foto.
func timeoutExceptUploads(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		withTimeout := middleware.Timeout(d)(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// upload video/foto & aliran notifikasi realtime (SSE) boleh berjalan lama
			if strings.HasSuffix(r.URL.Path, "/uploads") || strings.HasSuffix(r.URL.Path, "/notifications/stream") {
				next.ServeHTTP(w, r)
				return
			}
			withTimeout.ServeHTTP(w, r)
		})
	}
}

func Serve(ctx context.Context, cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) error {
	authSvc := auth.NewService(pool, cfg, log)
	if cfg.AdminMFADevSkip {
		log.Warn("MODE DEV: login admin melewati MFA (ADMIN_MFA_SKIP=false untuk mematikan; tidak berlaku di staging/production)")
	}
	if cfg.TelegramMode == "polling" {
		go authSvc.RunTelegramPolling(ctx)
	} else if cfg.TelegramMode == "webhook" {
		go authSvc.RegisterTelegramWebhook(ctx)
	}
	if cfg.WorkerInServe {
		StartWorkers(ctx, cfg, log, pool)
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

// StartWorkers menjalankan worker latar: terjemahan AI + notifikasi Telegram.
func StartWorkers(ctx context.Context, cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) {
	go NewTranslationWorker(cfg, log, pool).Run(ctx)
	go NewNotifyWorker(cfg, log, pool).Run(ctx)
	go (&media.Worker{Pool: pool, Files: newFiles(cfg), Log: log}).Run(ctx)
	go NewShop(cfg, log, pool).Worker(ctx)                                        // checkout tidak dibayar → batal, stok kembali
	go NewPayouts(cfg, log, pool).Worker(ctx)                                     // pesanan selesai → saldo; jadwal payout
	go (&returns.Service{Pool: pool, Files: newFiles(cfg), Log: log}).Worker(ctx) // retur tak dijawab seller 48 jam → admin
	go NewPhotoSearch(cfg, log, pool).Worker(ctx)                                 // foto produk → vektor untuk cari-dengan-foto
	go (&courierreturn.Service{Pool: pool, Log: log}).Worker(ctx)                 // courierreturn: paket kembali tak dikonfirmasi seller → diterima
}

// Pembayaran palsu (pengembangan) disimpan di memori → satu instance per proses.
var sharedFake = &payment.Fake{}

// FakePayments: penyedia pembayaran palsu bersama (dipakai /dev/pay dan tes).
func FakePayments() *payment.Fake { return sharedFake }

// NewPayments: payment gateway yang diatur super admin (database), dengan .env sebagai cadangan.
func NewPayments(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *paygw.Store {
	box, err := secretbox.New(cfg.DataKey)
	if err != nil {
		panic(err) // config sudah memvalidasi
	}
	st := &paygw.Store{Pool: pool, Box: box, Log: log, Def: paygw.Defaults{
		Env: cfg.Env, Mode: cfg.PaymentProvider, PayWayBaseURL: cfg.PayWayBaseURL, PayWayMerchantID: cfg.PayWayMerchantID,
		PayWayAPIKey: cfg.PayWayAPIKey, LifetimeMin: cfg.PaymentLifetimeMin, PublicCallbackURL: cfg.PublicBaseURL + "/webhooks/payway",
	}}
	if cfg.Env != "production" {
		st.Fake = sharedFake
	}
	return st
}

// NewIntegrations: provider SMS & Login Facebook yang diatur super admin (.env SMS_PROVIDER sebagai cadangan).
func NewIntegrations(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *integ.Store {
	box, err := secretbox.New(cfg.DataKey)
	if err != nil {
		panic(err)
	}
	return &integ.Store{Pool: pool, Box: box, Log: log, Def: integ.Defaults{Env: cfg.Env, SMSProvider: cfg.SMSProvider, BaseURL: cfg.PublicBaseURL,
		AnthropicKey: cfg.AnthropicKey, AnthropicModel: cfg.AnthropicModel, GeminiKey: cfg.GeminiKey, GeminiModel: cfg.GeminiModel,
		MapsKey: cfg.GoogleMapsKey, MapsMapID: cfg.GoogleMapsMapID,
		TelegramSet: cfg.TelegramBotToken != "", TelegramUsername: cfg.TelegramBotUsername}}
}

// AIKeys: kunci AI yang berlaku (Admin → Integrasi, cadangan .env), dibaca ulang tiap pemakaian.
// UsageRecorder: catat setiap percobaan AI ke laporan pemakaian (Admin → Integrasi → Pemakaian & saldo).
func UsageRecorder(st *integ.Store, purpose string) translate.Recorder {
	return func(ctx context.Context, provider, model string, in, out int, err error) {
		st.RecordUsage(ctx, integ.Usage{Provider: provider, Purpose: purpose, Model: model, In: in, Out: out, Err: err})
	}
}

func AIKeys(cfg config.Config, st *integ.Store) func(context.Context) translate.Keys {
	return func(ctx context.Context) translate.Keys {
		k := st.AI(ctx)
		return translate.Keys{AnthropicKey: k.AnthropicKey, AnthropicModel: k.AnthropicModel, AnthropicBaseURL: cfg.AnthropicBaseURL,
			GeminiKey: k.GeminiKey, GeminiModel: k.GeminiModel}
	}
}

// AITester: uji satu penyedia dengan menerjemahkan satu kata (dipakai tombol "Tes" di Admin → Integrasi).
func AITester(cfg config.Config) integ.Tester {
	return func(ctx context.Context, provider, key, model string) error {
		var t translate.Translator = &translate.Gemini{Key: key, Model: model}
		if provider == "anthropic" {
			t = &translate.Claude{Key: key, Model: model, BaseURL: cfg.AnthropicBaseURL}
		}
		tctx, cancel := context.WithTimeout(ctx, 25*time.Second)
		defer cancel()
		_, err := t.Translate(tctx, translate.Request{Kind: "product", Source: "en", Targets: []string{"km"}, Fields: map[string]string{"name": "Red cotton T-shirt"}})
		return err
	}
}

// NewPayouts: layanan payout. Kunci enkripsi rekening dari DATA_ENCRYPTION_KEY (dicek saat start).
func NewPayouts(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *payout.Service {
	box, err := secretbox.New(cfg.DataKey)
	if err != nil {
		panic(err) // config sudah memvalidasi; sampai sini berarti kunci rusak
	}
	return &payout.Service{Pool: pool, Box: box, Log: log}
}

func NewShop(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *shop.Service {
	return &shop.Service{Pool: pool, Pay: NewPayments(cfg, log, pool), Files: newFiles(cfg), BaseURL: cfg.PublicBaseURL,
		LifetimeMin: cfg.PaymentLifetimeMin, Log: log}
}

var (
	photoOnce sync.Once
	photoSvc  *photosearch.Service
)

// NewPhotoSearch: satu instance per proses (indeks vektor disimpan di memori).
func NewPhotoSearch(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *photosearch.Service {
	photoOnce.Do(func() {
		photoSvc = &photosearch.Service{Pool: pool, Files: newFiles(cfg), Log: log, Emb: &photosearch.HTTPEmbedder{URL: cfg.EmbedderURL}}
	})
	return photoSvc
}

func newFiles(cfg config.Config) storage.Local {
	return storage.Local{Dir: cfg.StorageDir, BaseURL: strings.TrimRight(cfg.PublicBaseURL, "/") + "/media"}
}

// NewNotifyWorker: tanpa TELEGRAM_BOT_TOKEN notifikasi tetap tersimpan (lonceng), hanya tidak dikirim ke Telegram.
func NewNotifyWorker(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *notify.Worker {
	w := &notify.Worker{Pool: pool, BaseURL: cfg.PublicBaseURL, SellerURL: cfg.SellerBaseURL, AdminURL: cfg.AdminBaseURL, Log: log, Push: NewPusher(cfg, pool)}
	if cfg.TelegramBotToken != "" {
		w.Bot = &notify.TelegramBot{Token: cfg.TelegramBotToken}
	}
	return w
}

// NewTranslationWorker dipakai `phsar serve` (WORKER_IN_SERVE=true) dan `phsar worker`.
func NewTranslationWorker(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *translate.Worker {
	ist := NewIntegrations(cfg, log, pool)
	t := newTranslator(cfg, AIKeys(cfg, ist), UsageRecorder(ist, "translate"))
	if t != nil {
		log.Info("terjemahan AI aktif: Anthropic utama, Gemini cadangan (kunci di Admin → Integrasi)", "provider", cfg.Translator,
			"cadangan_google", cfg.GoogleTranslateKey != "" && cfg.TranslationFallback == "google")
	}
	return &translate.Worker{Pool: pool, T: t, Log: log}
}

// newTranslator: "claude" (bawaan) = kunci dinamis dari Admin → Integrasi (Anthropic → Gemini → Google .env);
// fake/google/off tetap seperti sebelumnya. keys boleh nil (hanya .env).
func newTranslator(cfg config.Config, keys func(context.Context) translate.Keys, rec translate.Recorder) translate.Translator {
	if cfg.Translator == "claude" {
		if keys == nil {
			keys = func(context.Context) translate.Keys {
				return translate.Keys{AnthropicKey: cfg.AnthropicKey, AnthropicModel: cfg.AnthropicModel, AnthropicBaseURL: cfg.AnthropicBaseURL,
					GeminiKey: cfg.GeminiKey, GeminiModel: cfg.GeminiModel}
			}
		}
		var extra translate.Translator
		if cfg.TranslationFallback == "google" && cfg.GoogleTranslateKey != "" {
			extra = &translate.Google{Key: cfg.GoogleTranslateKey}
		}
		return translate.Dynamic{Keys: keys, Extra: extra, Record: rec}
	}
	return translate.New(translate.Config{Provider: cfg.Translator, AnthropicKey: cfg.AnthropicKey, AnthropicModel: cfg.AnthropicModel,
		GoogleKey: cfg.GoogleTranslateKey, Fallback: cfg.TranslationFallback})
}

// NewTranslator dipakai `phsar translate-test` (kunci dari .env saja).
func NewTranslator(cfg config.Config) translate.Translator { return newTranslator(cfg, nil, nil) }

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

// NewDemoSeeder: untuk `phsar seed-demo` (teks & media sudah disiapkan, jadi tanpa antrean terjemahan/kompres).
func NewDemoSeeder(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *demo.Seeder {
	files := newFiles(cfg)
	sellers := seller.NewService(pool, files)
	return &demo.Seeder{Pool: pool, Sellers: sellers, Uploader: storage.Uploader{Store: files, Pool: pool}, Files: files, Log: log}
}

// NewPusher: Web Push (VAPID). Kunci privat VAPID dienkripsi DATA_ENCRYPTION_KEY kalau diisi.
// Kontak untuk layanan push: PUSH_CONTACT (mailto:… / https://…), cadangan PUBLIC_BASE_URL kalau https.
func NewPusher(cfg config.Config, pool *pgxpool.Pool) *notify.Pusher {
	p := &notify.Pusher{Pool: pool, Subject: cfg.PushContact}
	if p.Subject == "" && strings.HasPrefix(cfg.PublicBaseURL, "https://") {
		p.Subject = strings.TrimRight(cfg.PublicBaseURL, "/")
	}
	if cfg.DataKey != "" {
		if box, err := secretbox.New(cfg.DataKey); err == nil {
			p.Box = box
		}
	}
	return p
}
