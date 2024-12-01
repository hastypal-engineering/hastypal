package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
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

func (tb *TelegramBot) SendMsg(msg types.SendTelegramMessage) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(msg)

	if jsonEncodeError != nil {
		return exception.
			New(jsonEncodeError.Error()).
			Trace("json.Marshal", "telegram-bot.go")
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
		return exception.
			New(requestCreationError.Error()).
			Trace("http.NewRequest", "telegram-bot.go")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return exception.
			New(requestError.Error()).
			Trace("client.Do", "telegram-bot.go")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, ioReaderErr := io.ReadAll(response.Body)

		if ioReaderErr != nil {
			return exception.
				New(ioReaderErr.Error()).
				Trace("io.ReadAll", "telegram-bot.go")
		}

		var data types.TelegramHttpResponse

		if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
			return exception.
				New(unmarshalErr.Error()).
				Trace("json.Unmarshal", "telegram-bot.go")
		}

		return exception.
			New(fmt.Sprintf("Error code: %d, Description: %s", data.ErrorCode, data.Description)).
			Trace("SendMsg", "telegram-bot.go")
	}

	return nil
}

func (tb *TelegramBot) AnswerCallbackQuery(msg types.AnswerCallbackQuery) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(msg)

	if jsonEncodeError != nil {
		return exception.
			New(jsonEncodeError.Error()).
			Trace("json.Marshal", "telegram-bot.go")
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
		return exception.
			New(requestCreationError.Error()).
			Trace("http.NewRequest", "telegram-bot.go")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return exception.
			New(requestError.Error()).
			Trace("client.Do", "telegram-bot.go")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, ioReaderErr := io.ReadAll(response.Body)

		if ioReaderErr != nil {
			return exception.
				New(ioReaderErr.Error()).
				Trace("io.ReadAll", "telegram-bot.go")
		}

		var data types.TelegramHttpResponse
		if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
			return exception.
				New(unmarshalErr.Error()).
				Trace("json.Unmarshal", "telegram-bot.go")
		}

		return exception.
			New(fmt.Sprintf("Error code: %d, Description: %s", data.ErrorCode, data.Description)).
			Trace("AnswerCallbackQuery", "telegram-bot.go")

	}

	return nil
}

func (tb *TelegramBot) SetCommands(commands []types.TelegramBotCommand) error {
	type SetCommand struct {
		Commands []types.TelegramBotCommand `json:"commands"`
	}

	byteEncodedBody, jsonEncodeError := json.Marshal(SetCommand{Commands: commands})

	if jsonEncodeError != nil {
		return exception.
			New(jsonEncodeError.Error()).
			Trace("json.Marshal", "telegram-bot.go")
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
		return exception.
			New(requestCreationError.Error()).
			Trace("http.NewRequest", "telegram-bot.go")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return exception.
			New(requestError.Error()).
			Trace("client.Do", "telegram-bot.go")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return exception.
			New("Failed request").
			Trace("client.Do", "telegram-bot.go")
	}

	return nil
}

func (tb *TelegramBot) SetWebhook(webhook types.TelegramWebhook) error {
	byteEncodedBody, jsonEncodeError := json.Marshal(webhook)

	if jsonEncodeError != nil {
		return exception.
			New(jsonEncodeError.Error()).
			Trace("json.Marshal", "telegram-bot.go")
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
		return exception.
			New(requestCreationError.Error()).
			Trace("http.NewRequest", "telegram-bot.go")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return exception.
			New(requestError.Error()).
			Trace("client.Do", "telegram-bot.go")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return exception.
			New("Response failed").
			Trace("client.Do", "telegram-bot.go")
	}

	return nil
}
