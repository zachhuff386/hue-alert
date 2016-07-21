package config

import (
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/errortypes"
	"io/ioutil"
	"os"
)

const (
	filename          = "hue.json"
	logPathDefault    = "hue.log"
	serverPortDefault = 9300
	serverHostDefault = "localhost"
)

var Config = &ConfigData{}

type ConfigData struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	LogPath    string `json:"log_path"`
	ServerPort int    `json:"server_port"`
	ServerHost string `json:"server_host"`
	Google     struct {
		Rate         int    `json:"rate"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google"`
	Accounts map[string]*account.Account `json:"accounts"`
}

func (c *ConfigData) CommitAccount(acct *account.Account) (err error) {
	for _, act := range c.Accounts {
		if act.Type == acct.Type && act.Identity == acct.Identity {
			acct.Id = act.Id
		}
	}

	c.Accounts[acct.Id] = acct

	err = Save()
	if err != nil {
		return
	}

	return
}

func (c *ConfigData) RemoveAccount(acctId string) (ok bool, err error) {
	_, ok = c.Accounts[acctId]
	if ok {
		delete(c.Accounts, acctId)
	}

	err = Save()
	if err != nil {
		return
	}

	return
}

func Load() (err error) {
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = &errortypes.ReadError{
				errors.Wrap(err, "config: File stat error"),
			}
		}
		return
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File read error"),
		}
		return
	}

	err = json.Unmarshal(file, Config)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File unmarshal error"),
		}
		return
	}

	if Config.LogPath == "" {
		Config.LogPath = logPathDefault
	}

	if Config.ServerPort == 0 {
		Config.ServerPort = serverPortDefault
	}

	if Config.ServerHost == "" {
		Config.ServerHost = serverHostDefault
	}

	if Config.Accounts == nil {
		Config.Accounts = map[string]*account.Account{}
	}

	return
}

func Save() (err error) {
	data, err := json.Marshal(Config)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File marshal error"),
		}
		return
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File write error"),
		}
		return
	}

	return
}
