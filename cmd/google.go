package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/server"
	"time"
)

func Google() (err error) {
	err = initConfig()
	if err != nil {
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
