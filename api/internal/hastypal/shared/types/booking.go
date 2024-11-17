package types

type Booking struct {
	Id         string `json:"id"`
	SessionId  string `json:"sessionId"`
	BusinessId string `json:"businessId"`
	ServiceId  string `json:"serviceId"`
	When       string `json:"when"`
	CreatedAt  string `json:"createdAt"`
}
