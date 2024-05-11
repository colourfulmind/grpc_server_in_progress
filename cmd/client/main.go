package main

import (
	"context"
	"log/slog"
	grpcclient "main/internal/clients/blog/grpc"
	"main/internal/config"
	"main/pkg/logger/logsetup"
	"main/pkg/logger/sl"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const op = "cmd/client/New"

	cfg := config.MustLoad()

	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting client", slog.Any("config", "cfg.Clients.GRPCClient"))

	cc, err := grpcclient.NewConnection(
		context.Background(),
		log,
		"localhost:8080",
		5,
		3*time.Second,
	)
	if err != nil {
		log.Error("failed to connect to server", op, sl.Err(err))
		os.Exit(1)
	}
	defer cc.Close()

	client := grpcclient.New(cc, log, cfg)
	go client.ConfigureRouter()
	client.Router.Run(":8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping client", slog.String("signal", sgl.String()))
	log.Info("client stopped")
}
