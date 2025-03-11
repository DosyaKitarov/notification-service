package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/DosyaKitarov/notification-service/internal/config"
	"github.com/DosyaKitarov/notification-service/internal/database"
)

func main() {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// DB connection setup
	db, err := database.ConnectToDB(cfg)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	defer db.Close()
	log.Printf("\n DB connected, configs:\n host: %v \n port: %v \n user: %v \n password: %v \n dbname: %v \n",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	// Router setup
	r := SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server launch error: %v", err)
	}
}
