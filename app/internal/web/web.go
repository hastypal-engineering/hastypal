// Package web
package web

type GetServicesReq struct {
	BusinessPublicID string `uri:"publicID" binding:"required"`
}

type ServiceDTO struct {
	ID          int
	Name        string
	Price       float64
	Currency    string
	Duration    string
	Description string
}

type BusinessDTO struct {
	Email       string
	Phone       string
	Address     string
	Name        string
	Description string
}

type BookingDTO struct {
	SessionID string
	Step      int
	Services  []*ServiceDTO
	Business  *BusinessDTO
}
