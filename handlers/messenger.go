package handlers

import (
	"fmt"
	"github.com/abhinavdahiya/go-messenger-bot"
	"github.com/yeyus/witgo/v1/witgo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MessengerHandler struct {
	Actions      map[string]func(*witgo.Session, witgo.EntityMap) (*witgo.Session, error)
	Bot          *mbotapi.BotAPI
	Mux          *http.ServeMux
	BotCallbacks <-chan mbotapi.Callback
	outgoing     chan OutgoingMessage
}

func NewMessengerHandler(bot *mbotapi.BotAPI) *MessengerHandler {
	callbacks, mux := bot.SetWebhook("/webhook")
	return &MessengerHandler{
		Actions:      make(map[string]func(*witgo.Session, witgo.EntityMap) (*witgo.Session, error)),
		Bot:          bot,
		Mux:          mux,
		BotCallbacks: callbacks,
	}
}

func (h *MessengerHandler) Action(session *witgo.Session, entities witgo.EntityMap, action string) (response *witgo.Session, err error) {
	return h.Actions[action](session, entities)
}

func (h *MessengerHandler) Say(session *witgo.Session, msg string) (response *witgo.Session, err error) {
	userID, err := h.parseSessionID(session.ID())
	if err != nil {
		return nil, err
	}

	nmsg := mbotapi.NewMessage(msg)
	h.outgoing <- OutgoingMessage{
		User:        mbotapi.User{ID: userID},
		Message:     nmsg,
		MessageType: mbotapi.RegularNotif,
	}

	return session, err
}

func (h *MessengerHandler) Merge(session *witgo.Session, entities witgo.EntityMap) (response *witgo.Session, err error) {
	// TODO: implement merging routines
	return nil, nil
}

func (h *MessengerHandler) Error(session *witgo.Session, msg string) {
	// TODO: implement error handling routines
}

func (h *MessengerHandler) Run() (chan<- witgo.SessionID, <-chan witgo.InputRecord) {
	var (
		requests = make(chan witgo.SessionID)
		records  = make(chan witgo.InputRecord)
	)
	// https://github.com/kurrik/witgo/blob/master/examples/02-twitter/main.go

	h.outgoing = make(chan OutgoingMessage, 10)
	h.runSend(h.outgoing)

	log.Printf("Callbacks setup")
	go func() {
		for callback := range h.BotCallbacks {
			log.Printf("[%#v] %s", callback.Sender, callback.Message.Text)

			records <- witgo.InputRecord{
				SessionID: h.makeSessionID(callback),
				Query:     callback.Message.Text,
			}
		}
	}()

	return requests, records
}

type OutgoingMessage struct {
	User        mbotapi.User
	Message     interface{}
	MessageType string
}

func (h *MessengerHandler) makeSessionID(callback mbotapi.Callback) witgo.SessionID {
	return witgo.SessionID(fmt.Sprintf("%v-%v", callback.Sender.ID, callback.Timestamp))
}

func (h *MessengerHandler) parseSessionID(sessionID witgo.SessionID) (userID int64, err error) {
	var lines []string
	lines = strings.Split(string(sessionID), "-")
	userID, err = strconv.ParseInt(lines[0], 10, 64)
	return
}

func (h *MessengerHandler) runSend(messages <-chan OutgoingMessage) {
	var (
		request OutgoingMessage
	)
	go func() {
		for request = range messages {
			log.Printf("Sending [User: %v] %v", request.User, request.Message)
			h.Bot.Send(request.User, request.Message, request.MessageType)
		}
	}()
}
