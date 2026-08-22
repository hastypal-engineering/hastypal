package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adriein/hastypal/pkg/constants"
	"github.com/rotisserie/eris"
)

func (stm *TelegramMessage) SessionExpired() TelegramMessage {
	var markdownText strings.Builder

	expiredSession := "![🙂‍↕️](tg://emoji?id=5368324170671202286) Lo sentimos, la sesión ha caducado\\!\n\n"

	processInstructionsIcon := "![‍ℹ️️](tg://emoji?id=5368324170671202286)"
	processInstructions := " *Pulsa Volver a empezar y te redirigiremos al canal de donde vienes*"

	markdownText.WriteString(expiredSession)
	markdownText.WriteString(processInstructionsIcon)
	markdownText.WriteString(processInstructions)

	startAgainButton := KeyboardButton{
		Text: "Volver a empezar",
		Url:  "t.me/+0djgKpMfYY5lY2I8",
	}

	chunked := [][]KeyboardButton{{startAgainButton}}

	return TelegramMessage{
		ChatId:         stm.ChatId,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: chunked},
	}
}

type AnswerCallbackQuery struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
}

type TelegramBot interface {
	SendMsg(dto BookingTelegramMessage) error
	AnswerCallbackQuery(msg AnswerCallbackQuery) error
}

type Bot struct {
	url   string
	token string
}

func NewTelegramBot(url string, token string) *Bot {
	return &Bot{
		url:   url,
		token: token,
	}
}

func (tb *Bot) SendMsg(dto BookingTelegramMessage) error {
	telegramMessage := dto.Message

	textWithHeader := fmt.Sprintf(
		"*%s* \\#%s\n\n%s",
		dto.BusinessName,
		dto.BookingSessionId,
		telegramMessage.Text,
	)

	updatedTelegramMessage := TelegramMessage{
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

		var data TelegramHttpResponse

		if err := json.Unmarshal(body, &data); err != nil {
			return eris.Wrap(err, "Error unmarshaling http response")
		}

		return eris.Wrapf(err, "Error code: %d, Description: %s", data.ErrorCode, data.Description)
	}

	return nil
}

func (tb *Bot) AnswerCallbackQuery(msg AnswerCallbackQuery) error {
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

		var data TelegramHttpResponse

		if err := json.Unmarshal(body, &data); err != nil {
			return eris.Wrap(err, "Error unmarshaling http response")
		}

		return eris.Wrapf(err, "Error code: %d, Description: %s", data.ErrorCode, data.Description)
	}

	return nil
}
