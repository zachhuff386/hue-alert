package cmd

import (
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/alert"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/constants"
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

	alrt := &alert.Alert{
		Hue:        he,
		Lights:     config.Config.Lights,
		Rate:       time.Duration(config.Config.UpdateRate) * time.Second,
		Mode:       constants.Solid,
		Brightness: config.Config.Brightness,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		<-sig
		alrt.Stop()
		os.Exit(1)
	}()

	alrt.Run()

	return
}
