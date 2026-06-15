package models

import "time"

type Product struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description,omitempty"`
	Quantity       float64    `json:"quantity"`
	Unit           string     `json:"unit"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CreateProductRequest struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Quantity       float64    `json:"quantity"`
	Unit           string     `json:"unit"`
	ExpirationDate *time.Time `json:"expiration_date"`
}
