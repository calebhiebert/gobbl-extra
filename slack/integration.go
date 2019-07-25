package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	gbl "github.com/calebhiebert/gobbl"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

const (
	FlagIsMessage       = "slack:ismessage"
	FlagEventMessage    = "slack:e:message"
	FlagChannel         = "slack:channel"
	FlagIsBot           = "slack:bot"
	FlagIsAppMention    = "slack:isappmention"
	FlagEventAppMention = "slack:e:appmention"
	FlagUser            = "slack:user"
)

// SlackIntegration is a slack gobbl integration
type SlackIntegration struct {
	api               *slack.Client
	verificationToken string
	bot               *gbl.Bot
}

// New creates a new slack integration
func New(verificationToken, authToken string, bot *gbl.Bot) *SlackIntegration {
	i := SlackIntegration{
		verificationToken: verificationToken,
		api:               slack.New(authToken),
		bot:               bot,
	}

	return &i
}

// User extracts the current user from the slack message
func (m *SlackIntegration) User(c *gbl.Context) (gbl.User, error) {
	var userId string
	var user gbl.User

	if c.HasFlag(FlagIsMessage) {
		msg := c.GetFlag(FlagEventMessage).(*slackevents.MessageEvent)

		if msg.BotID != "" {
			userId = msg.BotID
		}

		userId = msg.User
	} else if c.HasFlag(FlagIsAppMention) {
		am := c.GetFlag(FlagEventAppMention).(*slackevents.AppMentionEvent)

		userId = am.User
	}

	if userId != "" {
		slackUser, err := m.api.GetUserInfo(userId)
		if err != nil {
			return user, err
		}

		user.ID = slackUser.ID
		user.FirstName = slackUser.Name
		c.Flag(FlagUser, slackUser)
	}

	return user, nil
}

// Respond responds to an incoming request using the c.R property
func (m *SlackIntegration) Respond(c *gbl.Context) (*interface{}, error) {
	if c.R == nil {
		c.Debug("skipping slack response due to empty c.R")
		return nil, nil
	}

	if !c.HasFlag(FlagChannel) {
		c.Debugf("skipping slack response due to missing %s flag", FlagChannel)
		return nil, nil
	}

	mso := c.R.(slack.MsgOption)

	_, _, err := m.api.PostMessage(c.GetStringFlag(FlagChannel), mso)
	return nil, err
}

// GenericRequest extracts a generic request from slack
func (m *SlackIntegration) GenericRequest(c *gbl.Context) (gbl.GenericRequest, error) {
	genericRequest := gbl.GenericRequest{}

	req := c.RawRequest.(slackevents.EventsAPIEvent)

	switch evt := req.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		c.Flag(FlagIsMessage, true)
		c.Flag(FlagEventMessage, evt)
		c.Flag(FlagChannel, evt.Channel)
		genericRequest.Text = evt.Text

		if evt.SubType == "bot_message" && evt.BotID != "" {
			c.Flag(FlagIsBot, true)
		}

		break
	case *slackevents.AppMentionEvent:
		c.Flag(FlagIsAppMention, true)
		c.Flag(FlagEventAppMention, evt)
		c.Flag(FlagChannel, evt.Channel)
	}

	return genericRequest, nil
}

// ServeHTTP is a http request handler that is specifically built for accepting facebook webhook requests
func (m *SlackIntegration) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SLACK PANIC", r)
			fmt.Println(string(debug.Stack()))
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		panic(err)
	}

	body := buf.String()

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: m.verificationToken}))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		rw.Header().Set("Content-Type", "text")
		_, err = rw.Write([]byte(r.Challenge))
		if err != nil {
			panic(err)
		}
	case slackevents.CallbackEvent:
		inputContext := gbl.InputContext{
			RawRequest:  eventsAPIEvent,
			Integration: m,
			Response:    nil,
		}

		m.bot.Execute(&inputContext)
	default:

	}
}

func (m *SlackIntegration) GetAPI() *slack.Client {
	return m.api
}
