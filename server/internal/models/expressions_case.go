package models

type Expression struct {
	ID     string  `json:"id"`
	Value  string  `json:"-"`
	Status string  `json:"status"`
	Result string `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}
