package model

type Provider struct {
	Id           int    `json:"id" db:"id"`
	ProviderName string `json:"provider_name" db:"provider_name"`
	PhoneNumber  string `json:"phone_number" db:"phone_number"`
	Address      string `json:"address" db:"address"`
}
