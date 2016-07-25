package config

import (
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/constants"
	"github.com/zachhuff386/hue-alert/errortypes"
	"io/ioutil"
	"os"
)

const (
	filename          = "hue.json"
	logPathDefault    = "hue.log"
	serverPortDefault = 9300
	serverHostDefault = "localhost"
	brightnessDefault = 254
	updateRateDefault = 60
)

var Config = &ConfigData{}

type ConfigData struct {
	path       string   `json:"path"`
	loaded     bool     `json:"-"`
	Host       string   `json:"host"`
	Username   string   `json:"username"`
	Lights     []string `json:"lights"`
	LogPath    string   `json:"log_path"`
	ServerPort int      `json:"server_port"`
	ServerHost string   `json:"server_host"`
	Mode       string   `json:"mode"`
	Brightness int      `json:"brightness"`
	UpdateRate int      `json:"update_rate"`
	Google     struct {
		Color        string `json:"color"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google"`
	Slack struct {
		Color        string `json:"color"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"slack"`
	Accounts map[string]*account.Account `json:"accounts"`
}

func (c *ConfigData) Load(path string) (err error) {
	c.path = path

	_, err = os.Stat(c.path)
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

	file, err := ioutil.ReadFile(c.path)
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

	if c.Lights == nil {
		c.Lights = []string{}
	}

	if c.LogPath == "" {
		c.LogPath = logPathDefault
	}

	if c.ServerPort == 0 {
		c.ServerPort = serverPortDefault
	}

	if c.ServerHost == "" {
		c.ServerHost = serverHostDefault
	}

	if c.Accounts == nil {
		c.Accounts = map[string]*account.Account{}
	}

	if !constants.Modes.Contains(c.Mode) {
		c.Mode = constants.Medium
	}

	if c.Brightness == 0 {
		c.Brightness = brightnessDefault
	} else if c.Brightness > 254 {
		c.Brightness = 254
	} else if c.Brightness < 1 {
		c.Brightness = 1
	}

	if c.UpdateRate == 0 {
		c.UpdateRate = updateRateDefault
	}

	c.loaded = true

	return
}

func (c *ConfigData) Save() (err error) {
	if !c.loaded {
		err = &errortypes.WriteError{
			errors.New("config: Config file has not been loaded"),
		}
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File marshal error"),
		}
		return
	}

	err = ioutil.WriteFile(c.path, data, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File write error"),
		}
		return
	}

	return
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
	err = Config.Load(filename)
	if err != nil {
		return
	}

	return
}

func Save() (err error) {
	err = Config.Save()
	if err != nil {
		return
	}

	return
}
