package main

import (
	"ConsNot/internal/api"
	"ConsNot/internal/config"
	"ConsNot/internal/database"
	"ConsNot/internal/services"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

const workerCount = 12

func processCPF(cpf string, cfg config.Env, dbConn *sql.DB) {
	start := time.Now()
	fmt.Printf("[%s] [WORKER] Starting processing for CPF: %s\n", time.Now().Format(time.RFC3339), cpf)

	username := cfg.Acesso1
	password := cfg.Senha1
	loginName := username

	tokenResp, err := api.GetAccessToken(username, password)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error getting token for CPF %s: %v\n", time.Now().Format(time.RFC3339), cpf, err)
		return
	}
	fmt.Printf("[%s] [WORKER] Got access token for CPF: %s\n", time.Now().Format(time.RFC3339), cpf)

	responseContent, err := api.GetApiReturn(tokenResp.AuthRes.AccessToken, cpf, cfg.ApiKey)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error calling API for CPF %s: %v\n", time.Now().Format(time.RFC3339), cpf, err)
		return
	}
	fmt.Printf("[%s] [WORKER] API call successful for CPF: %s\n", time.Now().Format(time.RFC3339), cpf)

	searchDb, err := database.ConnectDb(cfg.Db_search)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error connecting to search DB: %v\n", time.Now().Format(time.RFC3339), err)
	}
	stormDb, err := database.ConnectDb(cfg.Db_storm)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error connecting to storm DB: %v\n", time.Now().Format(time.RFC3339), err)
	}
	var nome, numero string
	if searchDb != nil && stormDb != nil {
		pessoa, err := database.FetchConsultas(searchDb, stormDb, cpf)
		if err == nil {
			nome = pessoa.Nome
			numero = pessoa.Numero
			fmt.Printf("[%s] [WORKER] Found nome/numero for CPF %s: %s / %s\n", time.Now().Format(time.RFC3339), cpf, nome, numero)
		} else {
			fmt.Printf("[%s] [WORKER] No nome/numero found for CPF %s\n", time.Now().Format(time.RFC3339), cpf)
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

	err = database.InsertConsultaLog(dbConn, cpf, saldo, aviso, loginName, nome, numero, erro)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error inserting consulta log for CPF %s: %v\n", time.Now().Format(time.RFC3339), cpf, err)
	} else {
		fmt.Printf("[%s] [WORKER] Inserted consulta log for CPF %s\n", time.Now().Format(time.RFC3339), cpf)
	}

	err = database.UpdateConsultado(dbConn, cpf)
	if err != nil {
		fmt.Printf("[%s] [WORKER] Error marking as processed for CPF %s: %v\n", time.Now().Format(time.RFC3339), cpf, err)
	} else {
		fmt.Printf("[%s] [WORKER] CPF %s marked as processed\n", time.Now().Format(time.RFC3339), cpf)
	}

	fmt.Printf("[%s] [WORKER] Finished processing CPF: %s in %v\n", time.Now().Format(time.RFC3339), cpf, time.Since(start))
}

func main() {
	cfg := config.LoadEnv()
	for {
		fmt.Printf("[%s] [MAIN] Starting new processing cycle\n", time.Now().Format(time.RFC3339))
		dbConn, err := database.ConnectDb(cfg.Db_consultas)
		if err != nil {
			fmt.Printf("[%s] [MAIN] ERROR WHEN TRYING TO CONNECT TO DB : %v\n", time.Now().Format(time.RFC3339), err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if !services.IsAllowedTime() {
			fmt.Printf("[%s] [MAIN] System is not in allowed time, waiting...\n", time.Now().Format(time.RFC3339))
			time.Sleep(5 * time.Minute)
			continue
		}
		pause, err := database.IsPaused(dbConn)
		if err != nil {
			fmt.Printf("[%s] [MAIN] ERROR WHEN FETCHING PAUSED STATUS, RETRYING... %v\n", time.Now().Format(time.RFC3339), err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if pause {
			fmt.Printf("[%s] [MAIN] System is paused, waiting...\n", time.Now().Format(time.RFC3339))
			time.Sleep(5 * time.Minute)
			continue
		}

		fmt.Printf("[%s] [MAIN] System is not paused, starting processing...\n", time.Now().Format(time.RFC3339))

		forConsultar, err := database.FetchCustomers(dbConn)
		if err != nil {
			fmt.Printf("[%s] [MAIN] Couldn't find customers: %v\n", time.Now().Format(time.RFC3339), err)
			time.Sleep(5 * time.Minute)
			continue
		}
		if len(forConsultar) == 0 {
			fmt.Printf("[%s] [MAIN] No customers to process, sleeping...\n", time.Now().Format(time.RFC3339))
			time.Sleep(5 * time.Minute)
			continue
		}

		fmt.Printf("[%s] [MAIN] Queueing %d CPFs for processing with %d workers\n", time.Now().Format(time.RFC3339), len(forConsultar), workerCount)

		cpfChan := make(chan string)
		var wg sync.WaitGroup

		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func(workerID int) {
				fmt.Printf("[%s] [WORKER-%d] Started\n", time.Now().Format(time.RFC3339), workerID)
				defer func() {
					fmt.Printf("[%s] [WORKER-%d] Finished\n", time.Now().Format(time.RFC3339), workerID)
					wg.Done()
				}()
				for cpf := range cpfChan {
					fmt.Printf("[%s] [WORKER-%d] Picked up CPF: %s\n", time.Now().Format(time.RFC3339), workerID, cpf)
					processCPF(cpf, cfg, dbConn)
				}
			}(i + 1)
		}

		for _, cpf := range forConsultar {
			fmt.Printf("[%s] [MAIN] Queueing CPF: %s\n", time.Now().Format(time.RFC3339), cpf)
			cpfChan <- cpf
		}
		close(cpfChan)

		wg.Wait()

		fmt.Printf("[%s] [MAIN] All records processed or processing paused! Checking for new records in 10 seconds...\n", time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Second)
	}
}
