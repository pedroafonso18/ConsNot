package database

import (
	"database/sql"
	"fmt"
	"strings"
)

func FetchConsultas(client, storm_client *sql.DB, cpf string) (Pessoa, error) {
	cleaned_cpf := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
	var pessoa Pessoa

	query := "SELECT nome, numero FROM contatos WHERE cpf is not null AND (cpf = $1 OR cpf = $2)"
	fmt.Printf("Executing query: %s\n", query)
	fmt.Printf("With parameters: %s and %s\n", cpf, cleaned_cpf)

	rows, err := client.Query(query, cpf, cleaned_cpf)
	if err != nil {
		fmt.Printf("Erro ao buscar contato na base principal: %v\n", err)
		fmt.Println("Tentando query na base Storm...")
	} else {
		defer rows.Close()
		if rows.Next() {
			err := rows.Scan(&pessoa.Nome, &pessoa.Numero)
			if err != nil {
				fmt.Printf("Erro ao fazer scan do contato: %v\n", err)
				return Pessoa{}, err
			}
			fmt.Printf("Contato encontrado na base principal: %s - %s\n", pessoa.Nome, pessoa.Numero)
			return pessoa, nil
		}
		fmt.Printf("Nenhum contato encontrado na base principal para o CPF: %s\n", cpf)
		fmt.Println("Tentando buscar na base Storm...")
	}

	stormQuery := "SELECT cliente, telefone_celular FROM digitados_sistema WHERE cpf_cliente = $1 LIMIT 1"
	stormRows, err := storm_client.Query(stormQuery, cpf)
	if err != nil {
		fmt.Printf("Erro ao buscar na Storm: %v\n", err)
		return Pessoa{}, nil
	}
	defer stormRows.Close()
	if stormRows.Next() {
		err := stormRows.Scan(&pessoa.Nome, &pessoa.Numero)
		if err != nil {
			fmt.Printf("Erro ao fazer scan do contato na Storm: %v\n", err)
			return Pessoa{}, err
		}
		fmt.Printf("Contato encontrado na Storm: %s - %s\n", pessoa.Nome, pessoa.Numero)
		return pessoa, nil
	}
	fmt.Printf("Nenhum contato encontrado para o CPF: %s dentro da Storm.\n", cpf)
	return Pessoa{}, nil
}

func CountConsultas(client *sql.DB, user1, user2, user3, user4 string) (Logins, error) {
	query := `
		SELECT
			COUNT(CASE WHEN login = $1 THEN 1 END) as login1,
			COUNT(CASE WHEN login = $2 THEN 1 END) as login2,
			COUNT(CASE WHEN login = $3 THEN 1 END) as login3,
			COUNT(CASE WHEN login = $4 THEN 1 END) as login4
		FROM logs_consultas
		WHERE created_at >= NOW() - INTERVAL '24 hours'
	`
	var logins Logins
	err := client.QueryRow(query, user1, user2, user3, user4).Scan(
		&logins.Login1,
		&logins.Login2,
		&logins.Login3,
		&logins.Login4,
	)
	if err != nil {
		return Logins{}, err
	}
	return logins, nil
}

func FetchCurrentCampaign(client *sql.DB) (string, error) {
	var Campaign string
	err := client.QueryRow("SELECT campanha_ativa FROM config LIMIT 1").Scan(&Campaign)
	if err != nil {
		return "", err
	}
	return Campaign, nil
}

func IsPaused(client *sql.DB) (bool, error) {
	var isActive bool
	err := client.QueryRow("SELECT pausado FROM config").Scan(&isActive)
	if err != nil {
		return true, err
	}
	return isActive, nil
}

func FetchCustomers(client *sql.DB) ([]string, error) {
	var cpfs []string
	campanha_ativa, err := FetchCurrentCampaign(client)
	if err != nil {
		return nil, err
	}
	query := "SELECT cpf FROM consultar WHERE consultado = false AND campanha = $1 LIMIT 10"
	rows, err := client.Query(query, campanha_ativa)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cpf string
		if err := rows.Scan(&cpf); err != nil {
			return nil, err
		}
		cpfs = append(cpfs, cpf)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cpfs, nil
}
