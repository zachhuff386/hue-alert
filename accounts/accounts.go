package accounts

import (
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/config"
)

const (
	Oauth1 = 1
	Oauth2 = 2
)

func GetAccounts() (accts []*account.Account, err error) {
	accts = []*account.Account{}

	for _, acct := range config.Config.Accounts {
		accts = append(accts, acct)
	}

	return
}
