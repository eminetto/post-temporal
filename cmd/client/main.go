package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/eminetto/post-temporal/money_transfer"
	"github.com/eminetto/post-temporal/payment"

	"go.temporal.io/sdk/client"
)

// @@@SNIPSTART money-transfer-project-template-go-start-workflow
func main() {
	http.HandleFunc("/transfer", transferHandler)
	log.Println("Servidor iniciado na porta 8083")
	log.Fatal(http.ListenAndServe(":8083", nil))

}

type TransferRequest struct {
	Amount        float64 `json:"amount"`
	SourceAccount string  `json:"source_account"`
	TargetAccount string  `json:"target_account"`
	ReferenceID   string  `json:"reference_id"`
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
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

	var req TransferRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := payment.Details{
		SourceAccount: req.SourceAccount,
		TargetAccount: req.TargetAccount,
		Amount:        int(req.Amount),
		ReferenceID:   req.ReferenceID,
	}

	options := client.StartWorkflowOptions{
		ID:        req.ReferenceID,
		TaskQueue: payment.MoneyTransferTaskQueueName,
	}

	log.Printf("Starting transfer from account %s to account %s for %d", input.SourceAccount, input.TargetAccount, input.Amount)

	we, err := c.ExecuteWorkflow(r.Context(), options, money_transfer.MoneyTransfer, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sucesso"})
}
