package main

import (
	"fmt"
	"main/internal/config"
)

func main() {
	// TODO: setup config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: setup logger
	//log := logsetup.SetupLogger(cfg.Env)

	// TODO: run application
	//app := app.New(log, cfg)

	// TODO: graceful shutdown
}
