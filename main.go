package main

import (
	"flag"
	"github.com/zachhuff386/hue-alert/cmd"
)

func main() {
	flag.Parse()

	var err error

	switch flag.Arg(0) {
	case "hue-setup":
		err = cmd.HueSetup()
	case "hue-lights":
		err = cmd.HueLights()
	case "google-add":
		err = cmd.GoogleAdd()
	case "google-setup":
		err = cmd.GoogleSetup()
	case "slack-add":
		err = cmd.SlackAdd()
	case "slack-setup":
		err = cmd.SlackSetup()
	case "accounts":
		err = cmd.Accounts()
	case "account-remove":
		err = cmd.AccountRemove(flag.Arg(1), flag.Arg(2))
	case "start":
		err = cmd.Start()
	}

	if err != nil {
		panic(err)
	}
}
