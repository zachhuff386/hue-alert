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

func HueLights() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	he := initHue()
	if he == nil {
		return
	}

	lights, err := he.GetLights()
	if err != nil {
		return
	}

	lightIds := []string{}

	for _, light := range lights {
		confirm := ""

		fmt.Printf("Add %s (%s)? [y/N] ", light.Name, light.Type)
		fmt.Scanln(&confirm)

		if strings.HasPrefix(strings.ToLower(confirm), "y") {
			fmt.Printf("Light '%s' has been added...\n", light.Name)
			lightIds = append(lightIds, light.UniqueId)
		}
	}

	config.Config.Lights = lightIds

	err = config.Save()
	if err != nil {
		return
	}

	return
}
