package payu

type PaymentRequest struct {
	Key         string `json:"key"`
	TxnID       string `json:"txnid"`
	Amount      string `json:"amount"`
	ProductInfo string `json:"productinfo"`
	Firstname   string `json:"firstname"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Hash        string `json:"hash"`
	Furl        string `json:"furl"`
	Surl        string `json:"surl"`
}

func (p *PaymentRequest) SetHash(hash string) {
	p.Hash = hash
}