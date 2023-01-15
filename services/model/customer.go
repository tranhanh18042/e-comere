package model

type Customer struct {
	Id          int    `json:"id" db:"id"`
	Status      int    `json:"status" db:"status"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Address     string `json:"address" db:"address"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
}
