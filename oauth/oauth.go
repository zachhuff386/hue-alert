package oauth

var tokens = map[string]*Token{}

type Token struct {
	Id          string
	Type        string
	OauthSecret string
}
