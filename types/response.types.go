package types

type Response struct {
	Error   []ErrorDetail
	Status  uint
	Message string
}

type ErrorDetail struct {
	ErrorType    string
	ErrorMessage string
}

type SimpleError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
