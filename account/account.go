package account

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/errortypes"
	"reflect"
	"time"
)

var Authenticated chan bool

type Client interface {
	Update() error
	Sync() error
	SetAccount(acct *Account)
}

type Auth interface {
	Request() (string, error)
	Authorize(string, string) error
}

type Account struct {
	Id            string    `json:"id"`
	Type          string    `json:"type"`
	Identity      string    `json:"identity"`
	IdentityId    string    `json:"identity_id,omitempty"`
	OauthTokn     string    `json:"oauth_tokn,omitempty"`
	OauthSec      string    `json:"oauth_sec,omitempty"`
	Oauth2AccTokn string    `json:"oauth2_acc_tokn,omitempty"`
	Oauth2RefTokn string    `json:"oauth2_ref_tokn,omitempty"`
	Oauth2Exp     time.Time `json:"oauth2_exp,omitempty"`
	Alert         bool      `json:"-"`
}

func (a *Account) GetClient() (client Client, err error) {
	typ, ok := clientRegistry[a.Type]
	if !ok {
		err = &errortypes.UnknownError{
			errors.New("account: Invalid account type"),
		}
		return
	}

	val := reflect.New(typ).Elem()

	client = val.Addr().Interface().(Client)
	client.SetAccount(a)

	return
}

func (a *Account) GetColor() string {
	return colorRegistry[a.Type]
}

func GetAuth(acctType string) (auth Auth, authTyp int, err error) {
	typ, ok := authRegistry[acctType]
	if !ok {
		err = &errortypes.UnknownError{
			errors.New("account: Invalid account type"),
		}
		return
	}

	authTyp = authTypes[acctType]
	val := reflect.New(typ).Elem()

	auth = val.Addr().Interface().(Auth)

	return
}
