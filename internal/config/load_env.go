package config

import (
	"os"

	"github.com/joho/godotenv"
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
}

func LoadEnv() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, err
	}
	Acesso1 := os.Getenv("ACESSO_1")
	Acesso2 := os.Getenv("ACESSO_2")
	Acesso3 := os.Getenv("ACESSO_3")
	Acesso4 := os.Getenv("ACESSO_4")
	Senha1 := os.Getenv("SENHA_1")
	Senha2 := os.Getenv("SENHA_2")
	Senha3 := os.Getenv("SENHA_3")
	Senha4 := os.Getenv("SENHA_4")
	Db_consultas := os.Getenv("DB_CONSULTAS")
	Db_search := os.Getenv("DB_SEARCH")
	Db_storm := os.Getenv("DB_STORM")

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
	}, nil

}
