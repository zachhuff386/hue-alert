package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/hue"
	"strings"
)

func HueSetup() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	host := ""

	fmt.Print("Enter Hue Bridge hostname: ")
	fmt.Scanln(&host)

	config.Config.Host = strings.TrimSpace(host)

	he := hue.Hue{
		Host: config.Config.Host,
	}

	fmt.Print(
		"Press the link button on top of the Hue Bridge then press enter...")
	fmt.Scanln()

	err = he.Register()
	if err != nil {
		return
	}

	config.Config.Username = he.Username

	err = config.Save()
	if err != nil {
		return
	}

	return
}
