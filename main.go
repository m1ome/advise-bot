package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1ome/advise-bot/bot"
	"github.com/sirupsen/logrus"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "token", "", "telegram bot token")
	flag.Parse()
}

func main() {
	logrus.Info("intializing bot application")
	b, err := bot.New(token)
	if err != nil {
		logrus.Fatalf("error connecting to bot: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logrus.Info("shutting down bot")
		b.Stop()
	}()

	logrus.Infof("starting bot %s", b.Info())
	b.Start()
}
