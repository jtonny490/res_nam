package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"res_nam/internal/app"
	"res_nam/internal/config"
	"res_nam/internal/repositories"
	"res_nam/internal/services"
	"syscall"
	"time"
)

func main() {
	c := config.Load()
	db, err := app.New(c)
	if err != nil {
		log.Fatal(err)
	}
	r := app.BuildRouter(c, db)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	statusJob := services.ReportStatusJob{Reports: repositories.ReportRepository{DB: db}}
	go statusJob.RunPeriodically(ctx, time.Hour, func(err error) {
		log.Printf("stale report job failed: %v", err)
	})

	server := &http.Server{Addr: ":" + c.Port, Handler: r}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}()
	log.Printf("listening on :%s", c.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
