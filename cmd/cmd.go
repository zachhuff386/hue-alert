package cmd

import (
	"github.com/zachhuff386/hue-alert/config"
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
