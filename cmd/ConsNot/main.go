package main

import (
	"ConsNot/internal/api"
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
		pause, err := database.IsPaused(dbConn)
		if err != nil {
			fmt.Printf("ERROR WHEN FETCHING PAUSED STATUS, RETRYING... %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if pause {
			fmt.Print("System is paused, waiting...")
			time.Sleep(5 * time.Minute)
			continue
		}

		fmt.Print("System is not paused, starting processing...")

		forConsultar, err := database.FetchCustomers(dbConn)
		if err != nil {
			fmt.Printf("Couldn't find customers: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if len(forConsultar) == 0 {
			fmt.Printf("Couldn't find customers: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		contador := 0
		batchSize := 50

		for services.IsAllowedTime() && contador < len(forConsultar) {
			if contador > 0 && contador%10 == 0 {
				pause, err := database.IsPaused(dbConn)
				if err != nil {
					fmt.Printf("ERROR WHEN FETCHING PAUSED STATUS, RETRYING... %v\n", err)
					break
				}
				if pause {
					fmt.Print("System is paused, breaking processing loop...\n")
					break
				}
			}

			if !services.IsAllowedTime() {
				fmt.Println("Out of allowed time, waiting...")
				break
			}

			count, err := database.CountConsultas(
				dbConn,
				cfg.Acesso1,
				cfg.Acesso2,
				cfg.Acesso3,
				cfg.Acesso4,
			)
			if err != nil {
				fmt.Printf("Error counting consultas: %v\n", err)
				time.Sleep(10 * time.Second)
				continue
			}

			fmt.Printf("\n=== Consultas para cada login ===\n")
			fmt.Printf("Login 1: %v\n", count.Login1)
			fmt.Printf("Login 2: %v\n", count.Login2)
			fmt.Printf("Login 3: %v\n", count.Login3)
			fmt.Printf("Login 4: %v\n", count.Login4)

			batchEnd := contador + batchSize
			if batchEnd > len(forConsultar) {
				batchEnd = len(forConsultar)
			}
			batch := forConsultar[contador:batchEnd]

			for _, cpf := range batch {
				if !services.IsAllowedTime() {
					fmt.Println("Out of allowed time, breaking...")
					break
				}

				loginIdx := 0
				loginCounts := []int{count.Login1, count.Login2, count.Login3, count.Login4}
				loginLimits := []int{5000, 2000, 2000, 2000}
				loginUsers := []string{cfg.Acesso1, cfg.Acesso2, cfg.Acesso3, cfg.Acesso4}
				loginPasswords := []string{cfg.Senha1, cfg.Senha2, cfg.Senha3, cfg.Senha4}

				minCount := loginCounts[0]
				for i := 1; i < 4; i++ {
					if loginCounts[i] < minCount && loginCounts[i] < loginLimits[i] {
						minCount = loginCounts[i]
						loginIdx = i
					}
				}
				if loginCounts[loginIdx] >= loginLimits[loginIdx] {
					fmt.Println("All logins reached their limits, waiting for next window...")
					break
				}

				username := loginUsers[loginIdx]
				password := loginPasswords[loginIdx]
				loginCounts[loginIdx]++

				fmt.Printf("\nProcessing CPF: %s (%d/%d)\n", cpf, contador+1, len(forConsultar))

				tokenResp, err := api.GetAccessToken(username, password)
				if err != nil {
					fmt.Printf("Error getting token: %v\n", err)
					contador++
					continue
				}

				responseContent, err := api.GetApiReturn(tokenResp.AuthRes.AccessToken, cpf, cfg.ApiKey)
				if err != nil {
					fmt.Printf("Error calling API: %v\n", err)
					contador++
					continue
				}

				searchDb, err := database.ConnectDb(cfg.Db_search)
				if err != nil {
					fmt.Printf("Error connecting to search DB: %v\n", err)
				}
				stormDb, err := database.ConnectDb(cfg.Db_storm)
				if err != nil {
					fmt.Printf("Error connecting to storm DB: %v\n", err)
				}
				var nome, numero string
				if searchDb != nil && stormDb != nil {
					pessoa, err := database.FetchConsultas(searchDb, stormDb, cpf)
					if err == nil {
						nome = pessoa.Nome
						numero = pessoa.Numero
					}
				}

				saldo := ""
				aviso := ""
				erro := false
				if responseContent.Simulacoes != nil {
					saldo = responseContent.Simulacoes.ValorLiberado
				}
				if responseContent.Avisos != nil && len(*responseContent.Avisos) > 0 {
					aviso = (*responseContent.Avisos)[0].Aviso
				}
				if responseContent.Error != nil {
					erro = *responseContent.Error
				}

				loginName := loginUsers[loginIdx]

				err = database.InsertConsultaLog(dbConn, cpf, saldo, aviso, loginName, nome, numero, erro)
				if err != nil {
					fmt.Printf("Error inserting consulta log: %v\n", err)
				}

				err = database.UpdateConsultado(dbConn, cpf)
				if err != nil {
					fmt.Printf("Error marking as processed: %v\n", err)
				} else {
					fmt.Printf("CPF %s marked as processed\n", cpf)
				}

				contador++
			}

			if contador >= len(forConsultar) {
				fmt.Println("All records processed!")
				break
			}
		}

		fmt.Println("All records processed or processing paused! Checking for new records in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}
