package models

type Request struct {
	Instruction string  `json:"instruction"`
	Params      []int64 `json:"params"`
}

type Response struct {
	Message string `json:"message"`
	Value   int    `json:"value"`
	Success bool   `json:"success"`
}
