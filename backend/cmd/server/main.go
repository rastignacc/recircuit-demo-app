package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rmarko/electronics-marketplace/backend/internal/config"
	"github.com/rmarko/electronics-marketplace/backend/internal/handler"
	"github.com/rmarko/electronics-marketplace/backend/internal/middleware"
	"github.com/rmarko/electronics-marketplace/backend/internal/repository"
	"github.com/rmarko/electronics-marketplace/backend/internal/seed"
	"github.com/rmarko/electronics-marketplace/backend/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("connected to database")

	runMigrations(cfg.DatabaseURL, logger)

	seedCtx, seedCancel := context.WithTimeout(context.Background(), 30*time.Second)
	seed.Run(seedCtx, pool, logger)
	seedCancel()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	productSvc := service.NewProductService(productRepo)
	orderSvc := service.NewOrderService(orderRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	healthHandler := handler.NewHealthHandler(pool)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))

	r.Get("/healthz", healthHandler.Check)

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// Public product read
		r.Get("/categories", productHandler.ListCategories)
		r.Get("/products", productHandler.List)
		r.Get("/products/{id}", productHandler.GetByID)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(authSvc))

			r.Get("/orders", orderHandler.ListOrders)
			r.Get("/orders/{id}", orderHandler.GetByID)

			// Buyer-only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireBuyer)
				r.Post("/orders", orderHandler.PlaceOrder)
			})

			// Seller-only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireSeller)
				r.Post("/products", productHandler.Create)
				r.Put("/products/{id}", productHandler.Update)
				r.Delete("/products/{id}", productHandler.Delete)
				r.Get("/seller/products", productHandler.ListSellerProducts)
				r.Get("/seller/stats", orderHandler.SellerStats)
			})
		})
	})

	// Serve Vue SPA static files (production deployment)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "/static"
	}
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		logger.Info("serving static files", slog.String("dir", staticDir))
		fileServer := http.FileServer(http.Dir(staticDir))
		r.NotFound(spaHandler(staticDir, fileServer))
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("server shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", slog.String("error", err.Error()))
	}
	logger.Info("server stopped")
}

// spaHandler serves static files if they exist, otherwise falls back to index.html
// so that Vue Router's client-side routes work on page refresh.
func spaHandler(staticDir string, fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
		_, err := os.Stat(path)
		if os.IsNotExist(err) || err != nil {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, r)
	}
}

func runMigrations(databaseURL string, logger *slog.Logger) {
	m, err := migrate.New("file:///migrations", databaseURL)
	if err != nil {
		// Try local path for development
		m, err = migrate.New("file://migrations", databaseURL)
		if err != nil {
			logger.Warn("could not load migrations", slog.String("error", err.Error()))
			return
		}
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("migration failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	v, dirty, _ := m.Version()
	logger.Info("migrations applied", slog.Uint64("version", uint64(v)), slog.Bool("dirty", dirty))
}

