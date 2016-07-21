package cmd

import (
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

func initHue() (he hue.Hue, err error) {
	he = hue.Hue{
		Host:     config.Config.Host,
		Username: config.Config.Username,
	}

	if config.Config.Username == "" {
		err = he.Register()
		if err != nil {
			return
		}

		config.Config.Username = he.Username
		err = config.Save()
		if err != nil {
			return
		}
	}

	return
}
