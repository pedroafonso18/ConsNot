package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Acesso1      string
	Acesso2      string
	Acesso3      string
	Acesso4      string
	Senha1       string
	Senha2       string
	Senha3       string
	Senha4       string
	Db_consultas string
	Db_search    string
	Db_storm     string
	ApiKey       string
}

func LoadEnv() Env {
	file, err := os.Open(".env")
	if err != nil {
		fmt.Printf("Could not open .env file: %v", err.Error())
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if pos := strings.Index(line, "="); pos != -1 {
			key := strings.TrimSpace(line[:pos])
			value := strings.TrimSpace(line[pos+1:])
			envVars[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading .env file: %v", err.Error())
	}

	getEnv := func(key string) string {
		val, ok := envVars[key]
		if !ok || val == "" {
			fmt.Printf("Missing required environment variable: %s", key)
		}
		return val
	}

	Acesso1 := getEnv("ACESSO_1")
	Acesso2 := getEnv("ACESSO_2")
	Acesso3 := getEnv("ACESSO_3")
	Acesso4 := getEnv("ACESSO_4")
	Senha1 := getEnv("SENHA_1")
	Senha2 := getEnv("SENHA_2")
	Senha3 := getEnv("SENHA_3")
	Senha4 := getEnv("SENHA_4")
	Db_consultas := getEnv("DB_CONSULTAS")
	Db_search := getEnv("DB_SEARCH")
	Db_storm := getEnv("DB_STORM")
	Apikey := getEnv("APIKEY")

	return Env{
		Acesso1:      Acesso1,
		Acesso2:      Acesso2,
		Acesso3:      Acesso3,
		Acesso4:      Acesso4,
		Senha1:       Senha1,
		Senha2:       Senha2,
		Senha3:       Senha3,
		Senha4:       Senha4,
		Db_consultas: Db_consultas,
		Db_search:    Db_search,
		Db_storm:     Db_storm,
		ApiKey:       Apikey,
	}
}
