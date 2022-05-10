package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Update is a Telegram object that the handler receives every time a user interacts with the bot.
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object containing the user's message data and metadata.
type Message struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From indicates the user to which the message belongs.
type From struct {
	ID int `json:"id"`
}

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	ID int `json:"id"`
}

// ParseUpdate handles incoming update from the Telegram webhook
func ParseUpdate(jsonData string) (Update, error) {
	var update Update

	if err := json.Unmarshal([]byte(jsonData), &update); err != nil {
		return Update{}, fmt.Errorf("error on parsing telegram update: %s", err)
	}

	return update, nil
}

// SendMessage sends a text message to the Telegram chat identified by the given chat ID
func SendMessage(chatID int, text, parseMode string) (string, error) {
	endpoint := "https://api.telegram.org/bot" + os.Getenv("BOT_TOKEN") + "/sendMessage"

	res, err := http.PostForm(
		endpoint,
		url.Values{
			"chat_id":    {strconv.Itoa(chatID)},
			"text":       {text},
			"parse_mode": {parseMode},
		})
	if err != nil {
		return "", fmt.Errorf("error on sending message to telegram: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error on parsing answer from telegram: %s", err)
	}

	return string(body), nil
}
