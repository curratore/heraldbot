package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)

type Handler func(m telebot.Message) error

type Herald struct {
	Users    []telebot.User
	Handlers map[string]Handler

	bot   *telebot.Bot
	start time.Time
}

func NewHerald(token string) (*Herald, error) {
	bot, err := telebot.NewBot(token)
	if err != nil {
		return nil, err
	}

	h := &Herald{bot: bot}
	h.Handlers = map[string]Handler{
		"/hi":   h.RegisterUser,
		"/help": h.GetUsage,
		"/kill": h.KillCommand,
		"/log":  h.GetOutput,
		//		"/stats":  h.GetOutput,
		"/status": h.CommandStatus,
		"/who":    h.GetUsers,
	}

	return h, nil
}

func (h *Herald) Run() error {
	fmt.Println(h)
	h.start = time.Now()

	go h.listen()
	h.print()

	return nil
}

func (h *Herald) listen() error {
	messages := make(chan telebot.Message)
	h.bot.Listen(messages, 1*time.Second)

	for m := range messages {
		if h, ok := h.Handlers[m.Text]; ok {
			h(m)
		}
	}
	return nil
}

func (h *Herald) print() error {
	i := 0
	for {
		time.Sleep(time.Second * 20)
		for _, user := range h.Users {
			h.bot.SendMessage(user, fmt.Sprintf("Message %d", i), nil)
		}

		i++
	}
	return nil
}

func (h *Herald) RegisterUser(m telebot.Message) error {
	h.Users = append(h.Users)

	reply := fmt.Sprintf(
		"Hello, %s! You logged on this herald session.",
		m.Sender.FirstName,
	)

	return h.bot.SendMessage(m.Chat, reply, nil)
}

func (h *Herald) GetUsers(m telebot.Message) error {
	var names []string
	for _, user := range h.Users {
		names = append(names, user.FirstName)
	}

	return h.bot.SendMessage(m.Chat, strings.Join(names, ","), nil)
}

func (h *Herald) KillCommand(m telebot.Message) error {
	msg := fmt.Sprintf("Command killed (running for %s)", time.Since(h.start))

	for _, user := range h.Users {
		if user.ID == m.Chat.ID {
			h.bot.SendMessage(m.Chat, msg, nil)
		} else {
			h.bot.SendMessage(m.Chat, fmt.Sprintf("%s by %s", msg, m.Chat.FirstName), nil)
		}
	}
	os.Exit(0)

	return nil
}

func (h *Herald) CommandStatus(m telebot.Message) error {
	return h.bot.SendMessage(m.Chat,
		fmt.Sprintf("*Start at:* %s\n*Uptime: %s", h.start, time.Since(h.start)),
		&telebot.SendOptions{
			ParseMode:             telebot.ModeMarkdown,
			DisableWebPagePreview: true,
		},
	)
}

//func (h *Herald) GetStats(m telebot.Message)

func (h *Herald) GetOutput(m telebot.Message) error {
	file, err := telebot.NewFile("/tmp/sloth.jpg")
	if err != nil {
		return nil
	}

	document := &telebot.Document{
		File:     file,
		FileName: "sloth.jpg",
		Mime:     "image/jpeg",
	}

	return h.bot.SendDocument(m.Chat, document, nil)
}

func (h *Herald) GetUsage(m telebot.Message) error {

	reply := fmt.Sprintf(
		`herald-bot %s herald-bot is a Telegram bot that can do 
	Usage:

	Where: on telegram chat     
	`, VERSION)
	os.Exit(1)

	return h.bot.SendMessage(m.Chat, reply, nil)
}
