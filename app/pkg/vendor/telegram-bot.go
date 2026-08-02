package vendor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"github.com/rotisserie/eris"
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

func (tb *TelegramBot) SendMsg(dto types.BookingTelegramMessage) error {
	telegramMessage := dto.Message

	textWithHeader := fmt.Sprintf(
		"*%s* \\#%s\n\n%s",
		dto.BusinessName,
		dto.BookingSessionId,
		telegramMessage.Text,
	)

	updatedTelegramMessage := types.TelegramMessage{
		ChatId:         telegramMessage.ChatId,
		Text:           textWithHeader,
		ParseMode:      telegramMessage.ParseMode,
		ProtectContent: telegramMessage.ProtectContent,
		ReplyMarkup:    telegramMessage.ReplyMarkup,
	}

	byteEncodedBody, err := json.Marshal(updatedTelegramMessage)

	if err != nil {
		return eris.Wrap(err, "Error marshaling struct")
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/sendMessage",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if err != nil {
		return eris.Wrap(err, "Error creating new http request")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return eris.Wrap(err, "Error performing http request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)

		if err != nil {
			return eris.Wrap(err, "Error reading http body buffer")
		}

		var data types.TelegramHttpResponse

		if err := json.Unmarshal(body, &data); err != nil {
			return eris.Wrap(err, "Error unmarshaling http response")
		}

		return eris.Wrapf(err, "Error code: %d, Description: %s", data.ErrorCode, data.Description)
	}

	return nil
}

func (tb *TelegramBot) AnswerCallbackQuery(msg types.AnswerCallbackQuery) error {
	byteEncodedBody, err := json.Marshal(msg)

	if err != nil {
		return eris.Wrap(err, "Error marshaling struct")
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/bot%s/answerCallbackQuery",
			tb.url,
			tb.token,
		),
		bytes.NewBuffer(byteEncodedBody),
	)

	if err != nil {
		return eris.Wrap(err, "Error creating http request")
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return eris.Wrap(err, "Error performing http request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)

		if err != nil {
			return eris.Wrap(err, "Error reading http body buffer")
		}

		var data types.TelegramHttpResponse

		if err := json.Unmarshal(body, &data); err != nil {
			return eris.Wrap(err, "Error unmarshaling http response")
		}

		return eris.Wrapf(err, "Error code: %d, Description: %s", data.ErrorCode, data.Description)
	}

	return nil
}
