package deposit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eminetto/post-temporal/payment"
)

type DepositRequest struct {
	Amount  float64 `json:"amount"`
	Account string  `json:"account"`
}

func Deposit(ctx context.Context, data payment.Details) error {
	// Criar o payload da requisição
	depositReq := DepositRequest{
		Amount:  float64(data.Amount),
		Account: data.TargetAccount,
	}

	// Converter para JSON
	jsonData, err := json.Marshal(depositReq)
	if err != nil {
		return fmt.Errorf("erro ao converter para JSON: %w", err)
	}

	// Criar a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/deposit", bytes.NewBuffer(jsonData))
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
		return fmt.Errorf("erro na API de depósito: status %d", resp.StatusCode)
	}
	time.Sleep(1 * time.Second)
	return nil
}
