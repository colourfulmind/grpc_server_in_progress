package main

import (
	"fmt"
	"main/internal/app"
	"main/internal/config"
	"main/pkg/logger/logsetup"
)

func main() {
	// TODO: setup config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: setup logger
	log := logsetup.SetupLogger(cfg.Env)

	// TODO: run application
	app := app.New(log, cfg)

	// TODO: graceful shutdown
}
