package main

import (
	"net/http"

	"github-stats/internal/config"
	"github-stats/internal/githubapi"
	"github-stats/internal/server"
	"github-stats/internal/stats"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	client := githubapi.NewClient(cfg.GitHubToken)
	statsSvc := stats.NewService(client, log)

	c := cron.New()
	_, err = c.AddFunc("*/5 * * * *", func() {
		if _, err := statsSvc.GetUserStats("elyor04", true); err != nil {
			log.Warn("cron: failed to update user stats", zap.Error(err))
		}
		if _, err := statsSvc.GetLanguageStats("elyor04", true); err != nil {
			log.Warn("cron: failed to update language stats", zap.Error(err))
		}
	})
	if err != nil {
		log.Fatal("failed to schedule cron job", zap.Error(err))
	}
	c.Start()
	defer c.Stop()

	srv := server.New(statsSvc, log)

	log.Info("server starting", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, srv.Handler()); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
