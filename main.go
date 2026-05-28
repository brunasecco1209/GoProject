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

// ---------------- ENVIAR MENSAGEM ---------------- //

func enviarMensagem(token string, chatID int64, mensagem string) {

	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		token,
	)

	data := url.Values{}

	data.Set(
		"chat_id",
		fmt.Sprintf("%d", chatID),
	)

	data.Set(
		"text",
		mensagem,
	)

	http.PostForm(apiURL, data)
}

// ---------------- MAIN ---------------- //

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Erro ao carregar .env")
		return
	}

	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		fmt.Println("Token não encontrado!")
		return
	}

	fmt.Println("Bot conectado com sucesso!")

	ultimoUpdateID := 0

	for {

		link := fmt.Sprintf(
			"https://api.telegram.org/bot%s/getUpdates?offset=%d",
			token,
			ultimoUpdateID+1,
		)

		resp, err := http.Get(link)

		if err != nil {
			fmt.Println("Erro:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)

		resp.Body.Close()

		if err != nil {
			fmt.Println("Erro ao ler resposta:", err)
			continue
		}

		var updates UpdateResponse

		json.Unmarshal(body, &updates)

		for _, update := range updates.Result {

			ultimoUpdateID = update.UpdateID

			texto := strings.TrimSpace(update.Message.Text)

			chatID := update.Message.Chat.ID

			fmt.Println("Mensagem recebida:", texto)

			var resposta string

			if texto == "/start" {

				resposta = `🤖 Olá! Eu sou o Jarvis.

Eu consigo:
✅ Validar CPF
✅ Verificar se um CPF é válido ou inválido

📌 Basta enviar um CPF no chat.

Exemplo:
52998224725`

			} else {

				if validarCPF(texto) {
					resposta = "✅ CPF válido"
				} else {
					resposta = "❌ CPF inválido"
				}
			}

			enviarMensagem(
				token,
				chatID,
				resposta,
			)

			fmt.Println("Resposta enviada!")
		}

		time.Sleep(2 * time.Second)
	}
}