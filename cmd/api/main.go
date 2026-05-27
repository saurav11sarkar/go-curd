package main

import (
	"curd/internal/config"
	"curd/internal/db"
	"curd/internal/server"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	client, database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Disconnect(client); err != nil {
			log.Fatal(err)
		}
	}()
	routes := server.NewRoutes(database)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	if err := routes.Run(addr); err != nil {
		log.Fatal(err)
	}
}
