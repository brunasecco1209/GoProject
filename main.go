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