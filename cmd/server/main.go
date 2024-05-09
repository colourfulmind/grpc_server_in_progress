package main

import (
	"fmt"
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"main/pkg/logger/logsetup"
	"main/pkg/logger/sl"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: setup config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: setup logger
	log := logsetup.SetupLogger(cfg.Env)

	// TODO: run application
	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("cannot create server", sl.Err(err))
	}

	go application.Server.MustRun()

	// TODO: graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sgl := <-stop
	log.Info("stopping application", slog.String("signal", sgl.String()))
	application.Server.Stop()
	log.Info("application stopped")
}
