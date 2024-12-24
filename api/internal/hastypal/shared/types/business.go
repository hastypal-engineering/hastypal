package types

type Business struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	ContactPhone   string              `json:"contactPhone"`
	Email          string              `json:"email"`
	Password       string              `json:"password"`
	ServiceCatalog []ServiceCatalog    `json:"serviceCatalog"`
	OpeningHours   map[string][]string `json:"openingHours"`
	Holidays       map[string][]string `json:"holidays"`
	ChannelName    string              `json:"channelName"`
	Street         string              `json:"street"`
	PostCode       string              `json:"postCode"`
	City           string              `json:"city"`
	Country        string              `json:"country"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

type ServiceCatalog struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Currency   string `json:"currency"`
	Duration   string `json:"duration"`
	BusinessId string `json:"businessId"`
}

type BusinessConfig struct {
	Step    int8   `json:"step"`
	Content string `json:"content"`
}
