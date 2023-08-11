package optmiz

type Request struct {
	TransactionID string `json:"transaction_id"`
	Payload       []int  `json:"payload"`
}

type Response struct {
	TransactionID string `json:"transaction_id"`
	Experssion    string `json:"experssion"`
}
