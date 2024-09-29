package types

type Business struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	CommunicationPhone string `json:"communicationPhone"`
	Email              string `json:email`
	Password           string `json:password`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type BusinessConfig struct {
	Step    int8   `json:"step"`
	Content string `json:"content"`
}
