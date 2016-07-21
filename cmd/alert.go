package cmd

import (
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/hue"
	"github.com/zachhuff386/hue-alert/notification"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Alert() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	he := hue.Hue{
		Host:     config.Config.Host,
		Username: config.Config.Username,
	}

	lights, err := he.GetLights()
	if err != nil {
		return
	}

	notf := notification.Notification{
		Transition: 500 * time.Millisecond,
		Rate:       500 * time.Millisecond,
	}

	for _, light := range lights {
		notf.AddLight(light)
	}

	alrt := notification.Alert{
		Type:     "google",
		Color:    "#dd4c40",
		Duration: 500 * time.Millisecond,
	}

	notf.AddAlert(alrt)

	notf.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		<-sig
		notf.Stop()
		os.Exit(1)
	}()

	for {
		time.Sleep(1 * time.Minute)
	}

	return
}
