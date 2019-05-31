package tel

import (
	"fmt"
	"time"

	"github.com/calebhiebert/gobbl"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var EventTypeMessage = "message"
var EventTypeCallbackQuery = "callback_query"
var EventTypeUserLeave = "user_leave"
var EventTypeNewMembers = "user_join"
var EventTypeNewChatTitle = "new_chat_title"
var EventTypeNewChatPhoto = "new_chat_photo"
var EventTypeCommand = "command"

// GenericRequest will extract a generic request from a telegram event
func (i *Integration) GenericRequest(c *gbl.Context) (gbl.GenericRequest, error) {
	request := gbl.GenericRequest{}
	raw := c.RawRequest.(*tgbotapi.Update)

	fmt.Println(raw)
	fmt.Printf("%+v\n", raw)

	if raw.Message != nil {
		if raw.Message.LeftChatMember != nil {
			c.Flag("tel:eventtype", EventTypeUserLeave)
			c.Flag("tel:left_chat_member", raw.Message.LeftChatMember)
		} else if raw.Message.NewChatMembers != nil {
			c.Flag("tel:eventtype", EventTypeNewMembers)
			c.Flag("tel:new_chat_members", raw.Message.NewChatMembers)
		} else if raw.Message.NewChatTitle != "" {
			c.Flag("tel:eventtype", EventTypeNewChatTitle)
			c.Flag("tel:new_chat_title", raw.Message.NewChatTitle)
		} else if raw.Message.NewChatPhoto != nil {
			c.Flag("tel:eventtype", EventTypeNewChatPhoto)
			c.Flag("tel:new_chat_photo", raw.Message.NewChatPhoto)
		} else if raw.Message.IsCommand() {
			c.Flag("tel:eventtype", EventTypeCommand)
			c.Flag("tel:command", raw.Message.Command())
			c.Flag("tel:command_args", raw.Message.CommandArguments())
		} else {
			request.Text = raw.Message.Text
			c.Flag("tel:eventtype", EventTypeMessage)
		}

	} else if raw.CallbackQuery != nil {
		fmt.Printf("%+v\n", raw.CallbackQuery)

		request.Text = raw.CallbackQuery.Data
		c.Flag("tel:eventtype", EventTypeCallbackQuery)
		c.Flag("tel:callback_query_id", raw.CallbackQuery.ID)
		if raw.CallbackQuery.Message != nil {
			c.Flag("tel:callback_query_message_id", raw.CallbackQuery.Message.MessageID)
		}
	}

	return request, nil
}

// User will extract user data from telegram requests
func (i *Integration) User(c *gbl.Context) (gbl.User, error) {
	user := gbl.User{}
	raw := c.RawRequest.(*tgbotapi.Update)

	if raw.Message != nil {
		user.ID = fmt.Sprintf("%d", raw.Message.Chat.ID)
		user.FirstName = raw.Message.Chat.FirstName
		user.LastName = raw.Message.Chat.LastName
	} else if raw.CallbackQuery != nil {
		user.ID = fmt.Sprintf("%d", raw.CallbackQuery.From.ID)
		user.FirstName = raw.CallbackQuery.From.FirstName
		user.LastName = raw.CallbackQuery.From.LastName
	}

	return user, nil
}

// Respond will respond to telegram incoming messages
func (i *Integration) Respond(c *gbl.Context) (*interface{}, error) {
	if c.R == nil {
		return nil, nil
	}

	if c.User.ID == "" {
		return nil, fmt.Errorf("cannot respond because chatid missing")
	}

	var chatID int64

	raw := c.RawRequest.(*tgbotapi.Update)

	if raw.Message != nil {
		chatID = raw.Message.Chat.ID
	} else if raw.CallbackQuery != nil {
		chatID = raw.CallbackQuery.Message.Chat.ID
	}

	response := c.R.(*Response)

	if response.editMessage != nil {
		editMessage := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    chatID,
				MessageID: raw.CallbackQuery.Message.MessageID,
			},
			Text: response.editMessage.Text,
		}

		if response.editMessage.InlineButtons != nil {
			editMessage.BaseEdit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: response.editMessage.InlineButtons,
			}
		}

		i.API.Send(editMessage)
	}

	if len(response.Messages) > 0 {
		for idx, msg := range response.Messages {
			var chattable tgbotapi.Chattable
			var baseChat *tgbotapi.BaseChat

			baseChat = &tgbotapi.BaseChat{
				ChatID: chatID,
			}

			// Construct the outgoing telegram message
			if msg.InlineButtons != nil {
				baseChat.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: msg.InlineButtons,
				}
			}

			if idx == len(response.Messages)-1 && response.KeyboardMarkup != nil {
				if msg.InlineButtons != nil {
					c.Warn("Last message in group had inline keyboard, but was overwritten by reply keyboard")
				}

				baseChat.ReplyMarkup = *response.KeyboardMarkup
			}

			if idx == len(response.Messages)-1 && response.removeKeyboard != nil {
				if baseChat.ReplyMarkup != nil {
					c.Warn("remove keyboard reply markup overriding existing markup")
				}

				baseChat.ReplyMarkup = *response.removeKeyboard
			}

			if idx == len(response.Messages)-1 && response.forceReply != nil {
				if baseChat.ReplyMarkup != nil {
					c.Warn("force reply markup overriding existing markup")
				}

				baseChat.ReplyMarkup = *response.forceReply
			}

			if msg.ImageParameter != nil {
				chattable = tgbotapi.PhotoConfig{
					BaseFile: tgbotapi.BaseFile{
						BaseChat:    *baseChat,
						FileID:      *msg.ImageParameter,
						UseExisting: true,
					},
					Caption: msg.Text,
				}
			} else {
				chattable = tgbotapi.MessageConfig{
					BaseChat: *baseChat,
					Text:     msg.Text,
				}
			}

			// Set typing if required
			if msg.MinTypingTime > 0 {
				_, err := i.API.Send(tgbotapi.ChatActionConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: chatID,
					},
					Action: "typing",
				})
				if err != nil {
					c.Errorf("Error setting typing %v", err)
				} else {
					c.Trace("Set typing success")
				}

				time.Sleep(msg.MinTypingTime)
			}

			resp, err := i.API.Send(chattable)
			if err != nil {
				c.Errorf("Error sending telegram message %v", err)
			} else {
				c.Tracef("Sent telegram message %v", resp.MessageID)
			}
		}
	}

	if c.HasFlag("tel:callback_query_id") && response.CallbackQueryAnswer != nil {
		answer := response.CallbackQueryAnswer
		answer.CallbackQueryID = c.GetStringFlag("tel:callback_query_id")

		_, err := i.API.AnswerCallbackQuery(*answer)
		if err != nil {
			c.Errorf("Error answering callback query %v", err)
		}
	}

	return nil, nil
}
