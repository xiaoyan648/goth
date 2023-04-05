package goth

import "fmt"

type Token struct {
	Raw          map[string]interface{}
	UID          string
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Raw       map[string]interface{}
	UID       string
	Email     string
	Name      string
	Sex       int8 // 1为男性，2为女性
	AvatarURL string
	Location  string
}

type Provider interface {
	Name() string
	Token(code string) (*Token, error)
	GetUserInfo(*Token) (*User, error)
}

// Providers is list of known/available providers.
type Providers map[string]Provider

var providers = Providers{}

// UseProviders adds a list of available providers for use with Goth.
// Can be called multiple times. If you pass the same provider more
// than once, the last will be used.
func UseProviders(viders ...Provider) {
	for _, provider := range viders {
		providers[provider.Name()] = provider
	}
}

// GetProvider returns a previously created provider. If Goth has not
// been told to use the named provider it will return an error.
func GetProvider(name string) (Provider, error) {
	provider := providers[name]
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}
