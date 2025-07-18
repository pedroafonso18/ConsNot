package api

type ApiResponse struct {
	HtmlString string `json:"htmlString"`
	AppVersion string `json:"appVersion"`
}

type ResponseContent struct {
	Error               *bool      `json:"error,omitempty"`
	Avisos              *[]Aviso   `json:"avisos,omitempty"`
	IsInstability       *bool      `json:"isInstability,omitempty"`
	Nome                *string    `json:"nome,omitempty"`
	Cpf                 *string    `json:"cpf,omitempty"`
	DataNascimento      *string    `json:"dataNascimento,omitempty"`
	IdadeAnos           *uint32    `json:"idadeAnos,omitempty"`
	IdadeMeses          *uint32    `json:"idadeMeses,omitempty"`
	Telefone            *string    `json:"telefone,omitempty"`
	AbordarComo         *string    `json:"abordarComo,omitempty"`
	Simulacoes          *Simulacao `json:"simulacoes,omitempty"`
	ValorLiberado       *string    `json:"valorLiberado,omitempty"`
	PercentualComissao  *float64   `json:"percentualComissao,omitempty"`
	DisplayCardCreation *string    `json:"displayCardCreation,omitempty"`
}

type Aviso struct {
	Aviso string `json:"aviso"`
}

type Simulacao struct {
	ParcelaNova                  string           `json:"parcelaNova"`
	ValorOperacao                string           `json:"valorOperacao"`
	Prazo                        uint32           `json:"prazo"`
	BancoEmprestimo              string           `json:"bancoEmprestimo"`
	Taxa                         string           `json:"taxa"`
	NomeTabela                   string           `json:"nomeTabela"`
	Atualizar                    bool             `json:"atualizar"`
	IdOrgao                      uint32           `json:"idOrgao"`
	Coeficiente                  float64          `json:"coeficiente"`
	CodigoInternoTabela          string           `json:"codigoInternoTabela"`
	PercentualComissao           float64          `json:"percentualComissao"`
	PercentualProducao           float64          `json:"percentualProducao"`
	Ativa                        bool             `json:"ativa"`
	Produto                      string           `json:"produto"`
	ValorLiberado                string           `json:"valorLiberado"`
	SaldoTotal                   *string          `json:"saldoTotal,omitempty"`
	SalarioBruto                 string           `json:"salarioBruto"`
	ReservaDeSaldo               string           `json:"reservaDeSaldo"`
	Periodos                     []Periodo        `json:"periodos"`
	ApelidoBancoEmprestimo       string           `json:"apelidoBancoEmprestimo"`
	JurosTotais                  string           `json:"jurosTotais"`
	JurosPorAno                  string           `json:"jurosPorAno"`
	JurosPorMes                  string           `json:"jurosPorMes"`
	JurosEmprestimoPessoalPorAno string           `json:"jurosEmprestimoPessoalPorAno"`
	SimulacaoConfiavel           bool             `json:"simulacaoConfiavel"`
	SimulacaoFactivel            bool             `json:"simulacaoFactivel"`
	DataExpiracao                string           `json:"dataExpiracao"`
	MensagemProposta             MensagemProposta `json:"mensagemProposta"`
	Avisos                       []interface{}    `json:"avisos"`
	IsInstability                bool             `json:"isInstability"`
	Informacoes                  []interface{}    `json:"informacoes"`
	Objecoes                     []Objecao        `json:"objecoes"`
}

type Periodo struct {
	DataRepasse         string  `json:"DataRepasse"`
	RepasseMaximo       string  `json:"RepasseMaximo"`
	ValorFinanciado     string  `json:"ValorFinanciado"`
	PercentualRepassado *string `json:"PercentualRepassado,omitempty"`
	SaqueRestante       *string `json:"SaqueRestante,omitempty"`
}

type MensagemProposta struct {
	Nome                          string    `json:"nome"`
	PrimeiroNome                  string    `json:"primeiroNome"`
	Cpf                           string    `json:"cpf"`
	ApelidoBancoEmprestimoPort    string    `json:"apelidoBancoEmprestimoPort"`
	ValorLiberado                 string    `json:"valorLiberado"`
	TaxaContrato                  string    `json:"taxaContrato"`
	ReservaDeSaldo                string    `json:"reservaDeSaldo"`
	MensagemPeriodos              []Periodo `json:"mensagemPeriodos"`
	NumeroPeriodos                uint32    `json:"numeroPeriodos"`
	MensagemDocumentosNecessarios string    `json:"mensagemDocumentosNecessarios"`
}

type Objecao struct {
	Objecao          string  `json:"objecao"`
	Resposta         string  `json:"resposta"`
	RespostaCopiavel *string `json:"respostaCopiavel,omitempty"`
}

type AuthenticationResult struct {
	AccessToken string `json:"AccessToken"`
}

type AccessTokenResp struct {
	AuthRes AuthenticationResult `json:"AuthenticationResult"`
}
