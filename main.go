package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Estruturas da API do Telegram
type UpdateResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int `json:"update_id"`

	Message struct {
		Text string `json:"text"`

		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`

	} `json:"message"`
}

// ---------------- FUNÇÃO CPF ---------------- //

func validarCPF(cpf string) bool {

	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	for _, c := range cpf {

		if c < '0' || c > '9' {
			return false
		}
	}

	todosIguais := true

	for i := 1; i < 11; i++ {

		if cpf[i] != cpf[0] {

			todosIguais = false
			break
		}
	}

	if todosIguais {
		return false
	}

	// Primeiro dígito
	soma := 0

	for i := 0; i < 9; i++ {

		soma += int(cpf[i]-'0') * (10 - i)
	}

	digito1 := (soma * 10) % 11

	if digito1 == 10 {
		digito1 = 0
	}

	// Segundo dígito
	soma = 0

	for i := 0; i < 10; i++ {

		soma += int(cpf[i]-'0') * (11 - i)
	}

	digito2 := (soma * 10) % 11

	if digito2 == 10 {
		digito2 = 0
	}

	if digito1 == int(cpf[9]-'0') &&
		digito2 == int(cpf[10]-'0') {

		return true
	}

	return false
}
