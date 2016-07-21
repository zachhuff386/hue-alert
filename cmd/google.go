package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/server"
	"strings"
	"time"
)

func GoogleAdd() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	if config.Config.Google.ClientId == "" ||
		config.Config.Google.ClientSecret == "" {

		fmt.Println(
			"Google Oauth has not been setup. Run 'hue-alert google-setup'")
		return
	}

	account.InitAccounts()

	auth, _, err := account.GetAuth("google")
	if err != nil {
		return
	}

	url, err := auth.Request()
	if err != nil {
		return
	}

	fmt.Println("Open URL below to authenticate Google account:")
	fmt.Println(url)

	go func() {
		err = server.Server()
		if err != nil {
			panic(err)
		}
	}()

	<-account.Authenticated

	time.Sleep(1 * time.Second)

	fmt.Println("Google account successfully authenticated")

	return
}

func GoogleSetup() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	clientId := ""
	clientSecret := ""

	fmt.Print("Enter Google API ClientID: ")
	fmt.Scanln(&clientId)

	fmt.Print("Enter Google API ClientSecret: ")
	fmt.Scanln(&clientSecret)

	config.Config.Google.ClientId = strings.TrimSpace(clientId)
	config.Config.Google.ClientSecret = strings.TrimSpace(clientSecret)

	err = config.Save()
	if err != nil {
		return
	}

	return
}
