// Package business
package business

type Business struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ContactPhone string `json:"contactPhone"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Country      string `json:"country"`
	Password     string `json:"password"`
	Lang         string `json:"lang"`
	DateAdd      string `json:"dateAdd"`
	DateUpd      string `json:"dateUpd"`
}

