// Package web
package web

type GetServicesReq struct {
	BusinessID int
}

type ServiceDTO struct {
	ID          int
	SessionID   string
	Name        string
	Price       float64
	Currency    string
	Duration    string
	Description string
}
