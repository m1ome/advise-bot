package bot

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/m1ome/advise-bot/lib/fga"
	"github.com/sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type (
	Bot struct {
		bot    *telebot.Bot
		client *fga.Client
	}
)

func New(token string, verbose bool) (*Bot, error) {
	longPoller := &telebot.LongPoller{
		Timeout: time.Second * 10,
	}

	var poller telebot.Poller
	if verbose {
		poller = telebot.NewMiddlewarePoller(longPoller, func(update *telebot.Update) bool {
			text := ""
			if update.Message != nil {
				text = update.Message.Text
			}

			logrus.WithFields(logrus.Fields{
				"update_id":   update.ID,
				"update_text": text,
			}).Info("incoming telegram update")
			return true
		})
	} else {
		poller = longPoller
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: poller,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Reporter: func(err error) {
			if strings.Contains(err.Error(), telebot.ErrCouldNotUpdate.Error()) && verbose {
				logrus.Infof("cannot fetch updates from telegram: %v", err)
			} else {
				logrus.Errorf("error from telegram: %v", err)
			}
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
			if _, err := b.bot.Send(rec, fmt.Sprintf("error getting advise: %s", err.Error())); err != nil {
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
