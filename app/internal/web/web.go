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
	Slots     *BookingDatesDTO
}

type BookingPatchDTO struct {
	SessionID string
	ServiceID int
}

type TimeTable struct {
	Data []bool
}

func (t *TimeTable) MarkTimeSlot(start int, end int, symbol bool) {
	// Normal intra-day range (e.g., 09:00 to 17:00)
	if start <= end {
		for i := start; i < end; i++ {
			t.Data[i] = symbol
		}

		return
	}

	// Overnight range (e.g., 22:00 to 04:00 wrap-around)
	for i := start; i < 1440; i++ {
		t.Data[i] = symbol
	}
	for i := 0; i < end; i++ {
		t.Data[i] = symbol
	}
}

func (t *TimeTable) IsChunkAllTrue(start int, interval int) bool {
	if start+interval > len(t.Data) {
		return false
	}
	
	for _, minute := range t.Data[start : start+interval] {
		if !minute {
			return false
		}
	}
	return true
}

type SlotDTO struct {
	Hour        string
	IsAvailable bool
}

type BookingDatesDTO struct {
	Day   string
	Slots []*SlotDTO
}
