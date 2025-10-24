package withdraw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eminetto/post-temporal/payment"
)

const LIMIT = 10000

func Withdraw(ctx context.Context, data payment.Details) error {
	if data.Amount > LIMIT {
		fmt.Println("Limite de saque excedido")
		return &payment.OverLimitError{}
	}

	// Criar o payload da requisição
	withdrawReq := WithdrawRequest{
		Amount:  float64(data.Amount),
		Account: data.SourceAccount,
	}

	// Converter para JSON
	jsonData, err := json.Marshal(withdrawReq)
	if err != nil {
		return fmt.Errorf("erro ao converter para JSON: %w", err)
	}

	// Criar a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8081/withdraw", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir o header Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Fazer a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Verificar o status da resposta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro na API de saque: status %d", resp.StatusCode)
	}
	time.Sleep(1 * time.Second)
	return nil
}

type WithdrawRequest struct {
	Amount  float64 `json:"amount"`
	Account string  `json:"account"`
}

type WithdrawResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Balance float64 `json:"balance,omitempty"`
}
