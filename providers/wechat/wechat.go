package wechat

import (
	"encoding/json" //	"fmt"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xiaoyan648/goth"
)

type Provider struct {
	Appid  string
	Appkey string
	cli    *http.Client
}

// 必填 appid, appkey.
// 选填 cli， 不填默认为 DefaultClient.
func NewProvider(appid, appkey string, cli *http.Client) *Provider {
	return &Provider{
		Appid:  appid,
		Appkey: appkey,
		cli:    cli,
	}
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`

	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type UserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int64    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`

	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//获取token
func (e *Provider) GetToken(code string) (*goth.Token, error) {
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + e.Appid + "&secret=" + e.Appkey + "&code=" + code + "&grant_type=authorization_code"
	reqest, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	if e.cli != nil {
		client = e.cli
	}
	response, _ := client.Do(reqest)
	data, err := ioutil.ReadAll(response.Body)
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
		Raw:          raw,
		UID:          at.Openid,
		AccessToken:  at.AccessToken,
		ExpiresIn:    int64(at.ExpiresIn),
		RefreshToken: at.RefreshToken,
	}, nil
}

//获取第三方用户信息
func (e *Provider) GetUserInfo(access_token string, openid string) (*goth.User, error) {
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + access_token + "&openid=" + openid
	reqest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	if e.cli != nil {
		client = e.cli
	}
	response, _ := client.Do(reqest)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	u := UserInfo{}
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	if u.Errcode != 0 {
		return nil, fmt.Errorf("errcode: %v, errmsg: %v", u.Errcode, u.Errmsg)
	}

	return &goth.User{
		Raw:       raw,
		UID:       u.Openid,
		Name:      u.Nickname,
		AvatarURL: u.Headimgurl,
		Location:  u.City + u.Country,
		Sex:       int8(u.Sex),
	}, nil
}
