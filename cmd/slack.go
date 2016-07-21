package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/server"
	"strings"
	"time"
)

func SlackAdd() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	if config.Config.Slack.ClientId == "" ||
		config.Config.Slack.ClientSecret == "" {

		fmt.Println(
			"Slack Oauth has not been setup. Run 'hue-alert slack-setup'")
		return
	}

	account.InitAccounts()

	auth, _, err := account.GetAuth("slack")
	if err != nil {
		return
	}

	url, err := auth.Request()
	if err != nil {
		return
	}

	fmt.Println("Open URL below to authenticate Slack account:")
	fmt.Println(url)

	go func() {
		err = server.Server()
		if err != nil {
			panic(err)
		}
	}()

	<-account.Authenticated

	time.Sleep(1 * time.Second)

	fmt.Println("Slack account successfully authenticated")

	return
}

func SlackSetup() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	clientId := ""
	clientSecret := ""

	fmt.Print("Enter Slack API ClientID: ")
	fmt.Scanln(&clientId)

	fmt.Print("Enter Slack API ClientSecret: ")
	fmt.Scanln(&clientSecret)

	config.Config.Slack.ClientId = strings.TrimSpace(clientId)
	config.Config.Slack.ClientSecret = strings.TrimSpace(clientSecret)

	err = config.Save()
	if err != nil {
		return
	}

	return
}
