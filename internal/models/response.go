package models

type SuccesResponse struct {
	SuccesData interface{}
}

type FailResponse struct {
	Error   string
	Details string
}
