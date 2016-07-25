package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/constants"
	"github.com/zachhuff386/hue-alert/hue"
	"strings"
	"time"
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
		"Press the link button on top of the Hue bridge then press enter...")
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

	fmt.Println("Hue successfully linked")

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

	for {
		brightness := 0

		fmt.Print("Enter alert light brightness: [1-254] ")
		fmt.Scanln(&brightness)

		if brightness >= 1 && brightness <= 254 {
			config.Config.Brightness = brightness
			break
		}

		fmt.Println("Brightness is invalid...")
	}

	for {
		mode := ""

		fmt.Print("Enter alert light mode: [solid,slow,medium,fast] ")
		fmt.Scanln(&mode)

		if constants.Modes.Contains(mode) {
			config.Config.Mode = mode
			break
		}

		fmt.Println("Mode is invalid...")
	}

	err = config.Save()
	if err != nil {
		return
	}

	return
}

func HueTest(color string) (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	he := initHue()
	if he == nil {
		return
	}

	lights, err := he.GetLightsById(config.Config.Lights)
	if err != nil {
		return
	}

	if !strings.HasPrefix(color, "#") {
		color = "#" + color
	}

	for _, light := range lights {
		light.SetState(true)
		light.SetBrightness(254)
		light.SetColorHex(color)
		light.SetTransition(500 * time.Millisecond)

		err = light.Commit()
		if err != nil {
			return
		}
	}

	return
}
