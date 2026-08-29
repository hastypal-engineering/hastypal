package whatsapp

/*
================================================================================
WHATSAPP API DTO's
================================================================================
*/

type WhatsappUpdate struct {
	Object string          `json:"object"`
	Entry  []WhatsappEntry `json:"entry"`
}

type WhatsappEntry struct {
	ID      string           `json:"id"`
	Changes []WhatsappChange `json:"changes"`
}

type WhatsappChange struct {
	Value WhatsappValue `json:"value"`
	Field string        `json:"field"`
}

type WhatsappValue struct {
	MessagingProduct string            `json:"messaging_product"`
	Metadata         WhatsappMetadata  `json:"metadata"`
	Contacts         []WhatsappContact `json:"contacts"`
	Messages         []WhatsappMessage `json:"messages"`
}

type WhatsappMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WhatsappContact struct {
	Profile         WhatsappProfile `json:"profile"`
	WaID            string          `json:"wa_id"`
	IdentityKeyHash string          `json:"identity_key_hash"`
}

type WhatsappProfile struct {
	Name string `json:"name"`
}

type WhatsappMessage struct {
	From      string            `json:"from"`
	ID        string            `json:"id"`
	Timestamp string            `json:"timestamp"`
	Type      string            `json:"type"`
	Text      *WhatsappText     `json:"text,omitempty"`
	Context   *WhatsappContext  `json:"context,omitempty"`
	Referral  *WhatsappReferral `json:"referral,omitempty"`
}

type WhatsappText struct {
	Body string `json:"body"`
}

type WhatsappContext struct {
	From                string                   `json:"from,omitempty"`
	ID                  string                   `json:"id,omitempty"`
	ReferredProduct     *WhatsappReferredProduct `json:"referred_product,omitempty"`
	Forwarded           *bool                    `json:"forwarded,omitempty"`
	FrequentlyForwarded *bool                    `json:"frequently_forwarded,omitempty"`
}

type WhatsappReferredProduct struct {
	CatalogID         string `json:"catalog_id"`
	ProductRetailerID string `json:"product_retailer_id"`
}

type WhatsappReferral struct {
	SourceURL      string                  `json:"source_url"`
	SourceID       string                  `json:"source_id"`
	SourceType     string                  `json:"source_type"`
	Body           string                  `json:"body"`
	Headline       string                  `json:"headline"`
	MediaType      string                  `json:"media_type"`
	ImageURL       string                  `json:"image_url"`
	VideoURL       string                  `json:"video_url"`
	ThumbnailURL   string                  `json:"thumbnail_url"`
	CTWAClid       string                  `json:"ctwa_clid"`
	WelcomeMessage *WhatsappWelcomeMessage `json:"welcome_message"`
}

type WhatsappWelcomeMessage struct {
	Text string `json:"text"`
}
