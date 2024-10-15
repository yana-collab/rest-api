package main

import (
	"context"
	"log/slog"
	"os"
	"test/rest_api/internal/config"
	"test/rest_api/internal/storage/mongo"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()
	// TODO: init logger - slog
	log := setUpLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are inabled")

	// TODO: init storage - mongo
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		log.Warn("MONGODB_URI is not set, using default URI")
	}

	client, err := mongo.InitializeMongoDB(mongoURI)
	if err != nil {
		log.Error("Failed to initialize MongoDB", slog.String("error", err.Error()))
		return
	}

	mongo.EnsureCollectionAndDocumentExists(client, log)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Error("Error disconnecting MongoDB", slog.String("error", err.Error()))
		}
	}()

	log.Info("MongoDB initialized and ready for use!")

	// TODO: init router - chi, chi render

	// TODO: run server
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
