package database

import (
	"database/sql"
)

func InsertConsultaLog(client *sql.DB, cpf, saldo, aviso, login, nome, numero string, erro bool) error {
	erroStr := "true"
	if !erro {
		erroStr = "false"
	}
	if nome != "" && numero != "" {
		query := "INSERT INTO logs_consultas (cpf, saldo_consultado, erro, aviso, login, nome, numero) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err := client.Exec(query, cpf, saldo, erroStr, aviso, login, nome, numero)
		if err != nil {
			return err
		}
	} else {
		query := "INSERT INTO logs_consultas (cpf, saldo_consultado, erro, aviso, login) VALUES ($1, $2, $3, $4, $5)"
		_, err := client.Exec(query, cpf, saldo, erroStr, aviso, login)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateConsultado(client *sql.DB, cpf string) error {
	campanha, err := FetchCurrentCampaign(client)
	if err != nil {
		return err
	}
	_, err = client.Exec("UPDATE consultar SET consultado = true WHERE cpf = $1 AND campanha = $2", cpf, campanha)
	if err != nil {
		return err
	}
	return nil
}
