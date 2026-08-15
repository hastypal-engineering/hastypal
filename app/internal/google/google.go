package google

import "time"

type GoogleToken struct {
	BusinessID   string
	AccessToken  string
	TokenType    string
	RefreshToken string
	DateAdd      time.Time
	DateUpd      time.Time
}
