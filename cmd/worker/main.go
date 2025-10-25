package main

import (
	"log"

	"github.com/eminetto/post-temporal/deposit"
	"github.com/eminetto/post-temporal/money_transfer"
	"github.com/eminetto/post-temporal/payment"
	"github.com/eminetto/post-temporal/refund"
	"github.com/eminetto/post-temporal/withdraw"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "temporal-frontend.temporal.svc.cluster.local:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, payment.MoneyTransferTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(money_transfer.MoneyTransfer)
	w.RegisterActivity(withdraw.Withdraw)
	w.RegisterActivity(deposit.Deposit)
	w.RegisterActivity(refund.Refund)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
