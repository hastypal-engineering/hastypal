package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type TelegramBot struct {
	url   string
	token string
}

func NewTelegramBot(url string, token string) *TelegramBot {
	return &TelegramBot{
		url:   url,
		token: token,
	}
}

func (tb *TelegramBot) SendMsg() error {
	return nil
}

func (tb *TelegramBot) SetCommands(commands []types.TelegramBotCommand) error {
	type SetCommand struct {
		Commands []types.TelegramBotCommand `json:"commands"`
	}

	byteEncodedBody, jsonEncodeError := json.Marshal(SetCommand{Commands: commands})

	if jsonEncodeError != nil {
		return types.ApiError{
			Msg:      jsonEncodeError.Error(),
			Function: "SetCommands -> json.Marshal()",
			File:     "service/telegram-bot.go",
		}
	}

	request, requestCreationError := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/setMyCommands",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if requestCreationError != nil {
		return types.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "SetCommands -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types.ApiError{
			Msg:      requestError.Error(),
			Function: "SetCommands -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types.ApiError{
			Msg:      response.Status,
			Function: "SetCommands -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}

func (tb *TelegramBot) SetWebhook(webhook types.TelegramWebhook) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(webhook)

	if jsonEncodeError != nil {
		return types.ApiError{
			Msg:      jsonEncodeError.Error(),
			Function: "SetWebhook -> json.Marshal()",
			File:     "service/telegram-bot.go",
		}
	}

	request, requestCreationError := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/setWebhook",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if requestCreationError != nil {
		return types.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "SetWebhook -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types.ApiError{
			Msg:      requestError.Error(),
			Function: "SetWebhook -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types.ApiError{
			Msg:      response.Status,
			Function: "SetWebhook -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}
