package bot

import (
	"fmt"
	"time"

	"github.com/m1ome/advise-bot/lib/fga"
	"github.com/sirupsen/logrus"

	"gopkg.in/tucnak/telebot.v2"
)

type (
	Bot struct {
		bot    *telebot.Bot
		client *fga.Client
	}
)

func New(token string) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: time.Second * 10,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:    bot,
		client: fga.New(),
	}, nil
}

func (b *Bot) Info() string {
	return fmt.Sprintf("[#%d:@%s] %s %s", b.bot.Me.ID, b.bot.Me.Username, b.bot.Me.FirstName, b.bot.Me.LastName)
}

func (b *Bot) Start() {
	// Handle command
	b.bot.Handle("/advise", func(message *telebot.Message) {
		// Getting recipient
		var rec telebot.Recipient = message.Sender
		if message.FromGroup() {
			rec = message.Chat
		}

		advise, err := b.client.Random()
		if err != nil {
			logrus.Errorf("error fetching advise: %v", err)
			if _, err := b.bot.Send(rec, "error getting advise: %s", err.Error()); err != nil {
				logrus.Errorf("error sending failure message: %v", err)
			}

			return
		}

		if _, err := b.bot.Send(rec, advise); err != nil {
			logrus.Errorf("error sending message: %v", err)
		}
	})

	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
