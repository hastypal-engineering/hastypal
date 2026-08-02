package types

import "time"

type TelegramNotification struct {
	Id           string `json:"id"`
	SessionId    string `json:"sessionId"`
	BusinessId   string `json:"businessId"`
	BookingId    string `json:"bookingId"`
	ScheduledAt  string `json:"scheduledAt"`
	ChatId       int    `json:"chatId"`
	BusinessName string `json:"businessName"`
	ServiceName  string `json:"serviceName"`
	BookingDate  string `json:"bookingDate"`
	Sent         bool   `json:"sent"`
	SentAt       string `json:"sentAt,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

func (tn *TelegramNotification) MarkAsSent() {
	tn.Sent = true
	tn.SentAt = time.Now().Format(time.DateTime)
}
