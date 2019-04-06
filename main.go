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
	token   string
	verbose bool

	Version = "dev"
)

func init() {
	flag.StringVar(&token, "token", "", "telegram bot token")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()
}

func main() {
	logrus.Infof("init bot application: %s", Version)
	b, err := bot.New(token, verbose)
	if err != nil {
		logrus.Fatalf("error connecting to bot: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		logrus.Info("shutting down bot")
		b.Stop()
	}()

	logrus.Infof("starting bot %s", b.Info())
	b.Start()
}
