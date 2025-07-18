# 🤖 Robô de Consultas

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

## 📋 Descrição

O Robô de Consultas é uma aplicação automatizada desenvolvida em Go para realizar consultas em massa de informações de clientes através da API do simulador FGTS. O sistema gerencia múltiplos logins para distribuir as consultas, respeitando o limite de 5000 consultas por dia para cada credencial.

## ✨ Funcionalidades

- **Leitura de dados**: Importa CPFs e informações de clientes a partir de arquivos CSV
- **Autenticação inteligente**: Utiliza múltiplas credenciais de acesso, distribuindo as consultas entre elas
- **Controle de limites**: Monitora o número de consultas realizadas por cada login nas últimas 24h
- **Persistência de dados**: Armazena logs de consulta em banco de dados PostgreSQL
- **Tratamento de erros**: Identifica e registra erros nas consultas
- **Resiliência**: Continua operando mesmo quando algumas credenciais não estão disponíveis

## 🚀 Como usar

### Pré-requisitos

- Go (versão recente)
- PostgreSQL
- Acesso às credenciais da API

### Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/robo_consultas.git
cd robo_consultas
```

2. Configure o arquivo `.env` na raiz do projeto com suas credenciais:
```env
ACESSO1=consultas.noturnas.meuconsig@vipfinanceira.com.br
SENHA1=sua_senha_1
ACESSO2=consultas.noturnas.meuconsig2@vipfinanceira.com.br
SENHA2=sua_senha_2
ACESSO3=consultas.noturnas.meuconsig3@vipfinanceira.com.br
SENHA3=sua_senha_3
ACESSO4=consultas.noturnas.meuconsig4@vipfinanceira.com.br
SENHA4=sua_senha_4
DB_CONSULTAS=postgres://usuario:senha@host/banco_de_dados
DB_SEARCH=postgres://usuario:senha@host/banco_de_dados_search
DB_STORM=postgres://usuario:senha@host/banco_de_dados_storm
APIKEY=sua_apikey
```

3. Compile o projeto:
```bash
go build -o robo_consultas ./cmd/ConsNot
```

### Execução

1. Prepare o arquivo CSV com os CPFs para consulta. O arquivo deve seguir o formato:
```csv
CPF;Data Nascimento;Telefone;Valor Liberado
123.456.789-00;01/01/1990;(11) 99999-9999;0,00
```

2. Execute o programa:
```bash
./robo_consultas
```

## 🧰 Estrutura do projeto

```
robo_consultas/
├── cmd/
│   └── ConsNot/         # Ponto de entrada da aplicação
├── internal/
│   ├── api/             # Funções de autenticação e consulta à API
│   ├── config/          # Carregamento de configurações do .env
│   ├── database/        # Funções de acesso ao banco de dados
│   └── services/        # Serviços auxiliares (ex: controle de horário)
├── .env                 # Arquivo de configuração (credenciais)
├── go.mod               # Dependências do projeto
├── go.sum               # Sums das dependências
└── README.md            # Este arquivo
```

## 📊 Banco de dados

### Tabela: `logs_consultas`

| Coluna          | Tipo     | Descrição                                |
|-----------------|----------|-----------------------------------------|
| id              | SERIAL   | Identificador único da consulta         |
| cpf             | TEXT     | CPF consultado                           |
| saldo_consultado| TEXT     | Valor do saldo consultado (se disponível)|
| erro            | TEXT     | Indicador se houve erro (true/false)     |
| aviso           | TEXT     | Mensagem de aviso ou erro (se houver)    |
| login           | TEXT     | Credencial utilizada para a consulta     |
| created_at      | TIMESTAMP| Data e hora da consulta                  |

## 📝 Notas importantes

- O sistema limita-se a 5000 consultas por login em um período de 24 horas
- Caracteres especiais nos campos de senha (como $ e #) são tratados adequadamente
- O sistema utiliza as credenciais com menor número de consultas para balanceamento de carga
- Os resultados e erros são registrados no banco de dados para posterior análise

## 🔒 Segurança

- Credenciais são armazenadas apenas no arquivo `.env` (não incluído no controle de versão)
- A leitura direta do arquivo `.env` garante que caracteres especiais nas senhas sejam interpretados corretamente
- O banco de dados deve ser protegido por firewall e credenciais seguras

## 📄 Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE) - veja o arquivo LICENSE para mais detalhes.

---

Desenvolvido com ❤️ por Pedro Afonso