package models

import "time"

type RegisterRequest struct {
	Data struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}
