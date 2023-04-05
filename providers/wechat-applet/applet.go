package wechatapplet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xiaoyan648/goth"
)

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type Provider struct {
	AppID     string
	AppSecret string
	cli       *http.Client
	// sync.Pool cache body buffer
}

type AccessToken struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`

	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 必填 appid, secret.
// 选填 cli， 不填默认为 DefaultClient.
func NewProvider(appid, secret string, cli *http.Client) *Provider {
	return &Provider{
		AppID:     appid,
		AppSecret: secret,
		cli:       cli,
	}
}

func (p *Provider) GetToken(code string) (*goth.Token, error) {
	url := fmt.Sprintf(code2sessionURL, p.AppID, p.AppSecret, code)
	client := http.DefaultClient
	if p.cli != nil {
		client = p.cli
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	at := AccessToken{}
	if err := json.Unmarshal(data, &at); err != nil {
		return nil, err
	}
	if at.Errcode != 0 {
		return nil, fmt.Errorf("errcode: %v, errmsg: %v", at.Errcode, at.Errmsg)
	}

	return &goth.Token{
		Raw: raw,
		UID: at.Openid,
	}, nil
}

func (e *Provider) GetUserInfo(access_token string, openid string) (*goth.User, error) {
	return nil, nil
}
