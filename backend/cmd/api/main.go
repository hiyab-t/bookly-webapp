package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"cafe_store.hiyabnako/internal/data"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)


const version = "1.0.0"

type config struct {
	port int
	env string
	db struct { 
		dns string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime time.Duration
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func init() {
	godotenv.Load()
}

func main() {
	var cfg config

	dns := os.Getenv("DATABASE_URL")

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dns, "dns", dns,"database connection string")

	flag.IntVar(&cfg.db.maxOpenConns, "mxOpenConn", 25, "PostgreSQL maximum in-use and idle connecitons")
	flag.IntVar(&cfg.db.maxIdleConns, "mxIdleConns", 25, "PostgreSQL maximum idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "mzIdleTime", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbpool, err := OpenDatabase(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer dbpool.Close()

	logger.Info("database connection pool established")

	app := &application{config: cfg,
		logger: logger,
		models: data.NewModels(dbpool),}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.router(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server...", "addr:", srv.Addr, "env:", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func OpenDatabase(cfg config) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(context.Background(), cfg.db.dns)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbpool.Ping(ctx)
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}