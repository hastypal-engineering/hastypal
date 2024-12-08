package types

type Booking struct {
	Id         string `json:"id" db:"id"`
	SessionId  string `json:"sessionId" db:"session_id"`
	BusinessId string `json:"businessId" db:"business_id"`
	ServiceId  string `json:"serviceId" db:"service_id"`
	When       string `json:"when" db:"when"`
	CreatedAt  string `json:"createdAt" db:"created_at"`
}

func (s *Booking) DatabaseMappings() map[string]string {
	return map[string]string{
		"session_id": "session;id",
	}
}
