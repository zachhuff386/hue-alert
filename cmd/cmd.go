package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/hue"
	"github.com/zachhuff386/hue-alert/logger"
)

func initConfig() (err error) {
	logger.Init()

	err = config.Load()
	if err != nil {
		return
	}

	err = config.Save()
	if err != nil {
		return
	}

	return
}

func initHue() (he *hue.Hue) {
	if config.Config.Host == "" || config.Config.Username == "" {
		fmt.Println("Hue Bridge has not been setup. Run `hue-alert hue-setup'")
		return
	}

	he = &hue.Hue{
		Host:     config.Config.Host,
		Username: config.Config.Username,
	}

	return
}
