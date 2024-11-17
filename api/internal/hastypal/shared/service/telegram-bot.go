package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"io"
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

func (tb *TelegramBot) SendMsg(msg types2.SendTelegramMessage) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(msg)

	if jsonEncodeError != nil {
		return types2.ApiError{
			Msg:      jsonEncodeError.Error(),
			Function: "SendMsg -> json.Marshal()",
			File:     "service/telegram-bot.go",
		}
	}

	request, requestCreationError := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/sendMessage",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if requestCreationError != nil {
		return types2.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "SendMsg -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types2.ApiError{
			Msg:      requestError.Error(),
			Function: "SendMsg -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, ioReaderErr := io.ReadAll(response.Body)

		if ioReaderErr != nil {
			return types2.ApiError{
				Msg:      response.Status,
				Function: "SendMsg -> io.ReadAll()",
				File:     "service/telegram-bot.go",
			}
		}

		var data types2.TelegramHttpResponse
		if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
			return types2.ApiError{
				Msg:      unmarshalErr.Error(),
				Function: "SendMsg -> json.Unmarshal()",
				File:     "service/telegram-bot.go",
			}
		}

		return types2.ApiError{
			Msg:      fmt.Sprintf("Error code: %d, Description: %s", data.ErrorCode, data.Description),
			Function: "SendMsg",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}

func (tb *TelegramBot) AnswerCallbackQuery(msg types2.AnswerCallbackQuery) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(msg)

	if jsonEncodeError != nil {
		return types2.ApiError{
			Msg:      jsonEncodeError.Error(),
			Function: "AnswerCallbackQuery -> json.Marshal()",
			File:     "service/telegram-bot.go",
		}
	}

	request, requestCreationError := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/answerCallbackQuery",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if requestCreationError != nil {
		return types2.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "AnswerCallbackQuery -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types2.ApiError{
			Msg:      requestError.Error(),
			Function: "AnswerCallbackQuery -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, ioReaderErr := io.ReadAll(response.Body)

		if ioReaderErr != nil {
			return types2.ApiError{
				Msg:      response.Status,
				Function: "AnswerCallbackQuery -> io.ReadAll()",
				File:     "service/telegram-bot.go",
			}
		}

		var data types2.TelegramHttpResponse
		if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
			return types2.ApiError{
				Msg:      unmarshalErr.Error(),
				Function: "AnswerCallbackQuery -> json.Unmarshal()",
				File:     "service/telegram-bot.go",
			}
		}

		return types2.ApiError{
			Msg:      fmt.Sprintf("Error code: %d, Description: %s", data.ErrorCode, data.Description),
			Function: "AnswerCallbackQuery",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}

func (tb *TelegramBot) SetCommands(commands []types2.TelegramBotCommand) error {
	type SetCommand struct {
		Commands []types2.TelegramBotCommand `json:"commands"`
	}

	byteEncodedBody, jsonEncodeError := json.Marshal(SetCommand{Commands: commands})

	if jsonEncodeError != nil {
		return types2.ApiError{
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
		return types2.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "SetCommands -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types2.ApiError{
			Msg:      requestError.Error(),
			Function: "SetCommands -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types2.ApiError{
			Msg:      response.Status,
			Function: "SetCommands -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}

func (tb *TelegramBot) SetWebhook(webhook types2.TelegramWebhook) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(webhook)

	if jsonEncodeError != nil {
		return types2.ApiError{
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
		return types2.ApiError{
			Msg:      requestCreationError.Error(),
			Function: "SetWebhook -> http.NewRequest()",
			File:     "service/telegram-bot.go",
		}
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return types2.ApiError{
			Msg:      requestError.Error(),
			Function: "SetWebhook -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types2.ApiError{
			Msg:      response.Status,
			Function: "SetWebhook -> client.Do()",
			File:     "service/telegram-bot.go",
		}
	}

	return nil
}
