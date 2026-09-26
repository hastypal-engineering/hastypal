// Package web
package web

import "time"

type ErrorDTO struct {
	SessionID string
	Business  *BusinessDTO
}

type GetServicesReq struct {
	BusinessPublicID string `uri:"publicID" binding:"required"`
}

type GetDatesReq struct {
	BusinessPublicID string `uri:"publicID" form:"publicID" binding:"required"`
	SessionID        string `form:"sessionID"`
	Day              time.Time
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

type SelectableHourDTO struct {
	Hour        string
	IsAvailable bool
}

type BookingDatesDTO struct {
	Calendar    string
	HoursPerDay map[string]SelectableHourDTO
}
