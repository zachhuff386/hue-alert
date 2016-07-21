package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/accounts"
)

func Accounts() (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	account.InitAccounts()

	accts, err := accounts.GetAccounts()
	if err != nil {
		return
	}

	for _, acct := range accts {
		fmt.Printf("%s: %s\n", acct.Type, acct.Identity)

		client, e := acct.GetClient()
		if e != nil {
			err = e
			return
		}

		err = client.Sync()
		if err != nil {
			return
		}
	}

	return
}
