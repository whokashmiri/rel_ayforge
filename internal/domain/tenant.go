package domain

import "time"

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
