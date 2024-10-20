package types

type Business struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	ContactPhone   string              `json:"contactPhone"`
	Email          string              `json:"email"`
	Password       string              `json:"password"`
	ServiceCatalog []ServiceCatalog    `json:"serviceCatalog"`
	OpeningHours   map[string][]string `json:"openingHours"`
	ChannelName    string              `json:"channelName"`
	Location       string              `json:"location"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

type ServiceCatalog struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Currency string `json:"currency"`
}

type BusinessConfig struct {
	Step    int8   `json:"step"`
	Content string `json:"content"`
}
