package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xiaoyan648/goth"
)

type Provider struct {
	ClientId     string
	ClientSecret string //GitHub里所获取
	// RedirectUrl  string //重定向URL
}

func NewProvider(clientId, secret string) goth.Provider {
	return &Provider{
		ClientId:     clientId,
		ClientSecret: secret,
		// RedirectUrl:  redirectUrl,
	}
}

func (p *Provider) Name() string {
	return "github"
}

func (p *Provider) Token(code string) (*goth.Token, error) {
	url := getTokenAuthUrl(p.ClientId, p.ClientSecret, code)
	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	// 将响应体解析为 token，并返回
	var token goth.Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (p *Provider) GetUserInfo(token *goth.Token) (*goth.User, error) {
	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}
	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &goth.User{
		Raw: userInfo,
	}, nil
}

//获取地址
func getTokenAuthUrl(cid, secret, code string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		cid, secret, code,
	)
}
