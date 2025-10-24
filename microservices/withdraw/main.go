package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WithdrawRequest struct {
	Amount  float64 `json:"amount"`
	Account string  `json:"account"`
}

type WithdrawResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Balance float64 `json:"balance,omitempty"`
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req WithdrawRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(WithdrawResponse{
			Status:  "erro",
			Message: "Valor deve ser maior que zero",
		})
		return
	}

	if req.Account == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(WithdrawResponse{
			Status:  "erro",
			Message: "Conta é obrigatória",
		})
		return
	}

	// Aqui você pode processar o saque
	// Por exemplo, verificar saldo, debitar da conta, etc.

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(WithdrawResponse{
		Status:  "sucesso",
		Message: "Saque realizado com sucesso",
		Balance: 1000.0 - req.Amount, // Exemplo de saldo fictício
	})
}

func main() {
	http.HandleFunc("/withdraw", withdrawHandler)
	log.Println("Servidor de saque iniciado na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
