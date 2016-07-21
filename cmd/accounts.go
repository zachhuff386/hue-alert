package cmd

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/accounts"
	"github.com/zachhuff386/hue-alert/config"
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
	}

	return
}

func AccountRemove(typ, identity string) (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	account.InitAccounts()

	accts, err := accounts.GetAccounts()
	if err != nil {
		return
	}

	removed := false

	for _, acct := range accts {
		if acct.Type == typ && acct.Identity == identity {
			ok, e := config.Config.RemoveAccount(acct.Id)
			if e != nil {
				err = e
				return
			}

			if ok {
				removed = true
			}
		}
	}

	if removed {
		fmt.Println("Account successfully removed")
	} else {
		fmt.Println("Account not found")
	}

	return
}
