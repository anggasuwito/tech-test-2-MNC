package entity

type (
	Transaction struct {
		ID                  string `json:"id"`
		TransactionDetailID string `json:"transaction_detail_id"`
		Category            string `json:"category"`
		Status              string `json:"status"`
		Type                string `json:"type"`
		Amount              int64  `json:"amount"`
		BalanceBefore       int64  `json:"balance_before"`
		BalanceAfter        int64  `json:"balance_after"`
		Description         string `json:"description"`
		CreatedAt           string `json:"created_at"`
		UpdatedAt           string `json:"updated_at"`
	}
	TransactionTopupRequest struct {
		AccountID string `json:"-"`
		Amount    int64  `json:"amount"`
	}

	TransactionTopupResponse struct {
		*Transaction
	}

	TransactionTransferRequest struct {
		AccountID       string `json:"-"`
		TargetAccountID string `json:"target_account_id"`
		Amount          int64  `json:"amount"`
		Description     string `json:"description"`
	}

	TransactionTransferResponse struct {
		*Transaction
	}

	TransactionPaymentRequest struct {
		AccountID   string `json:"-"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
	}

	TransactionPaymentResponse struct {
		*Transaction
	}

	TransactionReportRequest struct {
		*ListPaginationRequest
	}

	TransactionReportResponse struct {
		*ListPaginationResponse
		Data []*Transaction `json:"data"`
	}
)
