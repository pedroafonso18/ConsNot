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
	username := cfg.Acesso1
	password := cfg.Senha1
	loginName := username

	tokenResp, err := api.GetAccessToken(username, password)
	if err != nil {
		fmt.Printf("Error getting token for CPF %s: %v\n", cpf, err)
		return
	}

	responseContent, err := api.GetApiReturn(tokenResp.AuthRes.AccessToken, cpf, cfg.ApiKey)
	if err != nil {
		fmt.Printf("Error calling API for CPF %s: %v\n", cpf, err)
		return
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

	err = database.InsertConsultaLog(dbConn, cpf, saldo, aviso, loginName, nome, numero, erro)
	if err != nil {
		fmt.Printf("Error inserting consulta log for CPF %s: %v\n", cpf, err)
	}

	err = database.UpdateConsultado(dbConn, cpf)
	if err != nil {
		fmt.Printf("Error marking as processed for CPF %s: %v\n", cpf, err)
	} else {
		fmt.Printf("CPF %s marked as processed\n", cpf)
	}
}

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
			fmt.Printf("No customers to process, sleeping...")
			time.Sleep(5 * time.Minute)
			continue
		}

		cpfChan := make(chan string)
		var wg sync.WaitGroup

		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for cpf := range cpfChan {
					processCPF(cpf, cfg, dbConn)
				}
			}()
		}

		for _, cpf := range forConsultar {
			cpfChan <- cpf
		}
		close(cpfChan)

		wg.Wait()

		fmt.Println("All records processed or processing paused! Checking for new records in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}
