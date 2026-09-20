// Package web
package web

type ErrorDTO struct {
	SessionID string
	Business  *BusinessDTO
}

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
	PublicID    string
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

type BookingPatchDTO struct {
	SessionID string
	ServiceID int
}
