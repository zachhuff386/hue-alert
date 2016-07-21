package accounts

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/errortypes"
	"github.com/zachhuff386/hue-alert/oauth"
)

const (
	slack = "slack"
)

var (
	slackConf *oauth.Oauth2
)

type SlackClient struct {
	acct *account.Account
}

func (g *SlackClient) SetAccount(acct *account.Account) {
	g.acct = acct
}

func (g *SlackClient) Update() (err error) {
	client := slackConf.NewClient(g.acct)

	err = client.Refresh()
	if err != nil {
		return
	}

	data := struct {
		Ok     bool   `json:"ok"`
		Error  string `json:"error"`
		User   string `json:"user"`
		UserId string `json:"user_id"`
	}{}

	err = client.GetJson(
		fmt.Sprintf("https://slack.com/api/auth.test?token=%s",
			g.acct.Oauth2AccTokn), &data)
	if err != nil {
		return
	}

	if !data.Ok {
		err = &errortypes.ApiError{
			errors.Newf("accounts.slack: Slack api error '%s'", data.Error),
		}
		return
	}

	g.acct.Identity = data.User
	g.acct.IdentityId = data.UserId

	err = config.Config.CommitAccount(g.acct)
	if err != nil {
		return
	}

	return
}

func (g *SlackClient) checkChannel(channelId string) (unread bool, err error) {
	client := slackConf.NewClient(g.acct)

	data := struct {
		Ok      bool   `json:"ok"`
		Error   string `json:"error"`
		Channel struct {
			Id                 string `json:"id"`
			Name               string `json:"name"`
			UnreadCount        int    `json:"unread_count"`
			UnreadCountDisplay int    `json:"unread_count_display"`
		} `json:"channel"`
		UserId string `json:"user_id"`
	}{}

	err = client.GetJson(
		fmt.Sprintf("https://slack.com/api/channels.info?token=%s&channel=%s",
			g.acct.Oauth2AccTokn, channelId), &data)
	if err != nil {
		return
	}

	if !data.Ok {
		err = &errortypes.ApiError{
			errors.Newf("accounts.slack: Slack api error '%s'", data.Error),
		}
		return
	}

	if data.Channel.UnreadCountDisplay != 0 {
		unread = true
	}

	return
}

func (g *SlackClient) Sync() (err error) {
	client := slackConf.NewClient(g.acct)

	err = client.Refresh()
	if err != nil {
		return
	}

	data := struct {
		Ok       bool   `json:"ok"`
		Error    string `json:"error"`
		Channels []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"channels"`
		UserId string `json:"user_id"`
	}{}

	err = client.GetJson(
		fmt.Sprintf("https://slack.com/api/channels.list"+
			"?token=%s&exclude_archived=1",
			g.acct.Oauth2AccTokn), &data)
	if err != nil {
		return
	}

	if !data.Ok {
		err = &errortypes.ApiError{
			errors.Newf("accounts.slack: Slack api error '%s'", data.Error),
		}
		return
	}

	alert := false

	for _, channel := range data.Channels {
		unread, e := g.checkChannel(channel.Id)
		if e != nil {
			err = e
			return
		}

		if unread {
			alert = true
		}
	}

	g.acct.Alert = alert

	return
}

type SlackAuth struct{}

func (g *SlackAuth) Request() (url string, err error) {
	url, err = slackConf.Request()
	if err != nil {
		return
	}

	return
}

func (g *SlackAuth) Authorize(state string, code string) (err error) {
	auth, err := slackConf.Authorize(state, code)
	if err != nil {
		return
	}

	client, err := auth.Account.GetClient()
	if err != nil {
		return
	}

	err = client.Update()
	if err != nil {
		return
	}

	return
}

func slackInit() {
	slackConf = &oauth.Oauth2{
		Type:         slack,
		ClientId:     config.Config.Slack.ClientId,
		ClientSecret: config.Config.Slack.ClientSecret,
		CallbackUrl: fmt.Sprintf("http://%s:%d/callback/slack",
			config.Config.ServerHost, config.Config.ServerPort),
		AuthUrl:  "https://slack.com/oauth/authorize",
		TokenUrl: "https://slack.com/api/oauth.access",
		Scopes: []string{
			"channels:read",
		},
	}
	slackConf.Config()
}

func init() {
	account.Register(slack, Oauth2, SlackAuth{}, SlackClient{}, slackInit)
}
