package payment

const MoneyTransferTaskQueueName = "TRANSFER_MONEY_TASK_QUEUE"

type Details struct {
	SourceAccount string
	TargetAccount string
	Amount        int
	ReferenceID   string
}

// InsufficientFundsError is raised when the account doesn't have enough money.
type InsufficientFundsError struct{}

func (m *InsufficientFundsError) Error() string {
	return "Insufficient Funds"
}

// InvalidAccountError is raised when the account number is invalid
type InvalidAccountError struct{}

func (m *InvalidAccountError) Error() string {
	return "Account number supplied is invalid"
}

// OverLimitError is raised when the account number is invalid
type OverLimitError struct{}

func (m *OverLimitError) Error() string {
	return "Withdrawal amount exceeded limit"
}
