package models

type Quotation struct {
	USDBRL coin
}

type coin struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	CodeIN string `json:"codein"`
	Bid    string `json:"bid"`
}
