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

	// TODO: run application

	// TODO: graceful shutdown
}
