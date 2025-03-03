package models

type AddReq struct {
	Expression string `json:"expression"`
}

type AddResp struct {
	ID string `json:"id"`
}
