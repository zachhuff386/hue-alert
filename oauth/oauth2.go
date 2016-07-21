package oauth

import (
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/errortypes"
	"github.com/zachhuff386/hue-alert/utils"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type Oauth2 struct {
	Type         string
	ClientId     string
	ClientSecret string
	CallbackUrl  string
	AuthUrl      string
	TokenUrl     string
	Scopes       []string
	conf         *oauth2.Config
}

func (o *Oauth2) Config() {
	o.conf = &oauth2.Config{
		ClientID:     o.ClientId,
		ClientSecret: o.ClientSecret,
		RedirectURL:  o.CallbackUrl,
		Scopes:       o.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  o.AuthUrl,
			TokenURL: o.TokenUrl,
		},
	}
}

func (o *Oauth2) Request() (url string, err error) {
	state := utils.RandStr(32)

	url = o.conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}

	tokn := &Token{
		Id:   state,
		Type: o.Type,
	}

	tokens[state] = tokn

	return
}

func (o *Oauth2) Authorize(state string, code string) (
	client *Oauth2Client, err error) {

	tokn, ok := tokens[state]
	if !ok || tokn.Type != o.Type {
		err = &errortypes.NotFoundError{
			errors.New("oauth.oauth2: State not found"),
		}
		return
	}

	accessTokn, err := o.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}

	acct := &account.Account{
		Id:            utils.RandStr(16),
		Type:          o.Type,
		Oauth2AccTokn: accessTokn.AccessToken,
		Oauth2RefTokn: accessTokn.RefreshToken,
		Oauth2Exp:     accessTokn.Expiry,
	}

	client = &Oauth2Client{
		Account: acct,
		Token:   *accessTokn,
		client:  o.conf.Client(oauth2.NoContext, accessTokn),
		conf:    o,
	}

	return
}

func (o *Oauth2) NewClient(acct *account.Account) (client *Oauth2Client) {
	tokn := &oauth2.Token{
		AccessToken:  acct.Oauth2AccTokn,
		TokenType:    "Bearer",
		RefreshToken: acct.Oauth2RefTokn,
		Expiry:       acct.Oauth2Exp,
	}

	client = &Oauth2Client{
		Account: acct,
		client:  o.conf.Client(oauth2.NoContext, tokn),
		conf:    o,
	}

	return
}

type Oauth2Client struct {
	oauth2.Token
	client  *http.Client
	conf    *Oauth2
	Account *account.Account
}

func (c *Oauth2Client) Refresh() (err error) {
	refreshed, err := c.Check()
	if err != nil {
		return
	}

	if !refreshed {
		return
	}

	c.Account.Oauth2AccTokn = c.AccessToken
	c.Account.Oauth2RefTokn = c.RefreshToken
	c.Account.Oauth2Exp = c.Expiry

	return
}

func (c *Oauth2Client) Check() (refreshed bool, err error) {
	tokn, err := c.client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown token error"),
		}
		return
	}

	refreshed = tokn.AccessToken != c.AccessToken
	if refreshed {
		c.AccessToken = tokn.AccessToken
		c.RefreshToken = tokn.RefreshToken
		c.Expiry = tokn.Expiry
	}

	return
}

func (c *Oauth2Client) GetJson(url string, resp interface{}) (err error) {
	httpResp, err := c.client.Get(url)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown parse error"),
		}
		return
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown parse error"),
		}
		return
	}

	return
}

func (c *Oauth2Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
