package accounts

import (
	"fmt"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/oauth"
	"github.com/zachhuff386/hue-alert/utils"
)

const (
	google = "google"
)

var (
	googleConf *oauth.Oauth2
)

type GooglelClient struct {
	acct *account.Account
}

func (g *GooglelClient) SetAccount(acct *account.Account) {
	g.acct = acct
}

func (g *GooglelClient) Update() (err error) {
	client := googleConf.NewClient(g.acct)

	err = client.Refresh()
	if err != nil {
		return
	}

	data := struct {
		EmailAddress string `json:"emailAddress"`
	}{}

	err = client.GetJson(
		"https://www.googleapis.com/gmail/v1/users/me/profile", &data)
	if err != nil {
		return
	}

	g.acct.Identity = data.EmailAddress

	err = config.Config.CommitAccount(g.acct)
	if err != nil {
		return
	}

	return
}

func (g *GooglelClient) Sync() (err error) {
	client := googleConf.NewClient(g.acct)

	err = client.Refresh()
	if err != nil {
		return
	}

	data := struct {
		Messages           []interface{} `json:"messages"`
		resultSizeEstimate int           `json:"resultSizeEstimate"`
	}{}

	err = client.GetJson(
		fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages"+
			"?maxResults=3&q=\"%s\"",
			utils.Escape("in:inbox is:unread"),
		), &data)
	if err != nil {
		return
	}

	g.acct.Alert = len(data.Messages) != 0

	return
}

type GoogleAuth struct{}

func (g *GoogleAuth) Request() (url string, err error) {
	url, err = googleConf.Request()
	if err != nil {
		return
	}

	return
}

func (g *GoogleAuth) Authorize(state string, code string) (err error) {
	auth, err := googleConf.Authorize(state, code)
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

func googleInit() {
	googleConf = &oauth.Oauth2{
		Type:         google,
		ClientId:     config.Config.Google.ClientId,
		ClientSecret: config.Config.Google.ClientSecret,
		CallbackUrl: fmt.Sprintf("http://%s:%d/callback/google",
			config.Config.ServerHost, config.Config.ServerPort),
		AuthUrl:  "https://accounts.google.com/o/oauth2/auth",
		TokenUrl: "https://www.googleapis.com/oauth2/v3/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
	}
	googleConf.Config()
}

func init() {
	account.Register(google, Oauth2, GoogleAuth{}, GooglelClient{}, googleInit)
}
