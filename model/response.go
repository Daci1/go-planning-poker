package model

type Response struct {
	HasErrors bool    `json:"has_errors"`
	Errors    []Error `json:"errors"`
}
