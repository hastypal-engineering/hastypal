package types

type BookingSession struct {
	Id         string `json:"id"`
	BusinessId string `json:"business_id"`
	ChatId     int    `json:"chat_id"`
	ServiceId  string `json:"service"`
	Date       string `json:"date"`
	Hour       string `json:"hour"`
	CreatedAt  string `json:"created_at"`
	Ttl        int64  `json:"ttl"`
}
