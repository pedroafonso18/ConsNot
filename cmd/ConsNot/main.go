package main

import (
	"ConsNot/internal/config"
	"ConsNot/internal/database"
	"ConsNot/internal/services"
	"fmt"
	"time"
)

func main() {
	cfg := config.LoadEnv()
	for {
		dbConn, err := database.ConnectDb(cfg.Db_consultas)
		if err != nil {
			fmt.Printf("ERROR WHEN TRYING TO CONNECT TO DB : %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if !services.IsAllowedTime() {
			fmt.Printf("System is not in allowed time, waiting...")
			time.Sleep(5 * time.Minute)
			continue
		}
	}
}
