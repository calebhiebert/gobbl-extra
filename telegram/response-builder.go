package tel

import (
	"fmt"
	"time"

	"github.com/calebhiebert/gobbl"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var ParseModeHTML = "HTML"
var ParseModeMarkdown = "Markdown"

// Response is a collection of telegram messages to be sent to the chat
type Response struct {
	Messages []*Message

	KeyboardMarkup      *tgbotapi.ReplyKeyboardMarkup
	CallbackQueryAnswer *tgbotapi.CallbackConfig
	editMessage         *Message
	removeKeyboard      *tgbotapi.ReplyKeyboardRemove
	forceReply          *tgbotapi.ForceReply
}

// Message is a single message that will be sent back to the chat
type Message struct {
	MinTypingTime  time.Duration
	Text           string
	InlineButtons  [][]tgbotapi.InlineKeyboardButton
	ImageParameter *string
}

// CreateResponse will add a response object to the context
// if it does not exist. It will then return this response object
func CreateResponse(c *gbl.Context) *Response {
	var response *Response

	if c.R == nil {
		response = &Response{
			Messages: []*Message{},
		}

		c.R = response
	}

	return response
}

// Text will add a single text message to the response
func (r *Response) Text(text string) *Message {
	message := &Message{
		Text:          text,
		MinTypingTime: 1 * time.Second,
	}

	r.Messages = append(r.Messages, message)

	return message
}

// RemoveKeyboard will remove any custom keyboards from a user's options
func (r *Response) RemoveKeyboard(selective bool) {
	r.removeKeyboard = &tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      selective,
	}
}

// ForceReply will force the user to reply to the bot's message
// https://core.telegram.org/bots/api#forcereply
func (r *Response) ForceReply(selective bool) {
	r.forceReply = &tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  selective,
	}
}

// EditText will create a new message that edits an existing message
// this is mostly useful in response to callback queries
func (r *Response) EditText(text string) *Message {
	message := &Message{
		Text:          text,
		MinTypingTime: 1 * time.Second,
	}

	r.editMessage = message

	return message
}

// Image will create an image message, the imageParameter should be either
// a URL or a telegram file id
func (r *Response) Image(caption, imageParameter string) *Message {
	message := &Message{
		Text:           caption,
		ImageParameter: &imageParameter,
		MinTypingTime:  1 * time.Second,
	}

	r.Messages = append(r.Messages, message)
	return message
}

// KeyboardConfig will set keyboard configuration options
func (r *Response) KeyboardConfig(resize, oneTime, selective bool) {
	if r.KeyboardMarkup == nil {
		r.KeyboardMarkup = &tgbotapi.ReplyKeyboardMarkup{}
	}

	r.KeyboardMarkup.ResizeKeyboard = resize
	r.KeyboardMarkup.OneTimeKeyboard = oneTime
	r.KeyboardMarkup.Selective = selective
}

// Keyboard will add or add to a row of keyboard buttons
func (r *Response) Keyboard(row int, buttons ...*tgbotapi.KeyboardButton) {
	if r.KeyboardMarkup == nil {
		r.KeyboardMarkup = &tgbotapi.ReplyKeyboardMarkup{}
	}

	if r.KeyboardMarkup.Keyboard == nil {
		r.KeyboardMarkup.Keyboard = [][]tgbotapi.KeyboardButton{}
	}

	if row > len(r.KeyboardMarkup.Keyboard)+1 {
		panic(fmt.Sprintf("row of %d invalid: current row count %d", row, len(r.KeyboardMarkup.Keyboard)))
	}

	if len(r.KeyboardMarkup.Keyboard) == 0 || len(r.KeyboardMarkup.Keyboard) < row+1 {
		r.KeyboardMarkup.Keyboard = append(r.KeyboardMarkup.Keyboard, []tgbotapi.KeyboardButton{})
	}

	for _, btn := range buttons {
		r.KeyboardMarkup.Keyboard[row] = append(r.KeyboardMarkup.Keyboard[row], *btn)
	}
}

// AnswerQuery will use telegrams answerCallbackQuery api to send the
// user a notification or alert in response to a callback query.
// This method will have no effect if the request is not a callback query
func (r *Response) AnswerQuery(text, url string, alert bool) {
	r.CallbackQueryAnswer = &tgbotapi.CallbackConfig{
		Text:      text,
		ShowAlert: alert,
		URL:       url,
		CacheTime: 0,
	}
}

// InlineButton will add an inline button to the response message
// Row is the 0-based row where the button will be placed
// this method will panic if the row parameter is not possible
// ie. there is currently 1 row, and a row parameter of 3 was specified
func (m *Message) InlineButton(row int, buttons ...*tgbotapi.InlineKeyboardButton) *Message {

	if m.InlineButtons == nil {
		m.InlineButtons = [][]tgbotapi.InlineKeyboardButton{}
	}

	if row > len(m.InlineButtons)+1 {
		panic(fmt.Sprintf("row of %d invalid: current row count %d", row, len(m.InlineButtons)))
	}

	if len(m.InlineButtons) == 0 || len(m.InlineButtons) < row+1 {
		m.InlineButtons = append(m.InlineButtons, []tgbotapi.InlineKeyboardButton{})
	}

	for _, btn := range buttons {
		m.InlineButtons[row] = append(m.InlineButtons[row], *btn)
	}

	return m
}

// IBURL is a helper for constructing an inline url button
func IBURL(text, url string) *tgbotapi.InlineKeyboardButton {
	return &tgbotapi.InlineKeyboardButton{
		Text: text,
		URL:  &url,
	}
}

// IBCallback is a helper for constructing an inline callback button
func IBCallback(text, callbackData string) *tgbotapi.InlineKeyboardButton {
	return &tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}
}

// KeyboardButton is a helper for constructing a keyboard button
func KeyboardButton(text string) *tgbotapi.KeyboardButton {
	return &tgbotapi.KeyboardButton{
		Text: text,
	}
}

// SpecialKeyboardButton is a helper for constructing a keyboard button
func SpecialKeyboardButton(text string, contact, location bool) *tgbotapi.KeyboardButton {
	return &tgbotapi.KeyboardButton{
		Text:            text,
		RequestContact:  contact,
		RequestLocation: location,
	}
}
