package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetAccessToken(access, password string) (AccessTokenResp, error) {
	url := "https://cognito-idp.us-east-2.amazonaws.com"
	body := map[string]interface{}{
		"AuthParameters": map[string]string{
			"USERNAME": access,
			"PASSWORD": password,
		},
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "63ccaojkma1th1pucikhn1n19k",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error when trying to marshal the JSON: %v\n", err)
		return AccessTokenResp{}, err
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		return AccessTokenResp{}, err
	}
	httpReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
	httpReq.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("HTTP request failed: %v", err)
		return AccessTokenResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 226 || resp.StatusCode < 200 {
		fmt.Printf("Request failed with status: %s", resp.Status)
		return AccessTokenResp{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var accessTokenResp AccessTokenResp
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return AccessTokenResp{}, err
	}

	err = json.Unmarshal(bodyBytes, &accessTokenResp)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\nBody: %s\n", err, string(bodyBytes))
		return AccessTokenResp{}, err
	}

	return accessTokenResp, nil

}

func GetApiReturn(token, cpf string) (ResponseContent, error) {
	cleaned_cpf := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
	url := "https://consig.private.app.br/api/secao/simulador-safra-fgts"
	body := map[string]interface{}{
		"numCpf":           cleaned_cpf,
		"bancoDestinoNovo": "9993-MB",
		"autorizacao":      true,
		"salarioBruto":     "",
		"mesesTrabalhados": "",
		"dtNascimentoAux":  false,
		"numTelefone":      "",
		"numeroDeParcelas": "10",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error when Marshalling body: %v", err)
		return ResponseContent{}, nil
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error when creating request: %v", err)
		return ResponseContent{}, nil
	}
	httpReq.Header.Set("accesstoken", token)
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("HTTP request failed: %v", err)
		return ResponseContent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 226 || resp.StatusCode < 200 {
		fmt.Printf("Request failed with status: %s", resp.Status)
		return ResponseContent{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var responseContent ResponseContent
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return ResponseContent{}, err
	}

	err = json.Unmarshal(bodyBytes, &responseContent)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\nBody: %s\n", err, string(bodyBytes))
		return ResponseContent{}, err
	}

	return responseContent, nil

}
