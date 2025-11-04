package models

type LoginRequest struct {
	Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"data"`
}
