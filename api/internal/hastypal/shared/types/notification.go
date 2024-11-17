package types

type TelegramNotification struct {
	Id          string `json:"id"`
	SessionId   string `json:"sessionId"`
	BusinessId  string `json:"businessId"`
	BookingId   string `json:"bookingId"`
	ScheduledAt string `json:"scheduledAt"`
	ChatId      int    `json:"chatId"`
	From        string `json:"from"`
	CreatedAt   string `json:"createdAt"`
}
