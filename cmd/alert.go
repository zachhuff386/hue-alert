package cmd

import (
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/alert"
	"github.com/zachhuff386/hue-alert/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	he := initHue()
	if he == nil {
		return
	}

	account.InitAccounts()

	alrt := alert.Alert{
		Hue:    he,
		Lights: config.Config.Lights,
		Rate:   5 * time.Second,
	}

	alrt.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		<-sig
		alrt.Stop()
		os.Exit(1)
	}()

	for {
		time.Sleep(1 * time.Minute)
	}

	return
}
