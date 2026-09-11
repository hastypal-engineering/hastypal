// Package business
package business

import "time"

type Business struct {
	ID             int
	Name           string
	ContactPhone   string
	Email          string
	Address        string
	Country        string
	Password       string
	Lang           string
	ServiceCatalog []*ServiceCatalog
	DateAdd        time.Time
	DateUpd        time.Time
}

type ServiceCatalog struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Currency    string
	Duration    string
	BusinessID  int
	DateAdd     time.Time
	DateUpd     time.Time
}
