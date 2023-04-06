package main

import (
	"fmt"
	"net/http"

	"github.com/xiaoyan648/goth"
	"github.com/xiaoyan648/goth/providers/github"
	"github.com/xiaoyan648/goth/providers/wechat"
	wechatapplet "github.com/xiaoyan648/goth/providers/wechat-applet"
)

func main() {
	goth.UseProviders(
		github.NewProvider("1dc4c54e64e3dd1cd79c", "8555221c37aec8215c2df49b36c5cd3481db67d3"),
		wechat.NewProvider("xxx", "xxx", nil),
		wechatapplet.NewProvider("xxx", "xxx", http.DefaultClient),
	)
	p, _ := goth.GetProvider("github")
	token, err := p.Token("testcode")
	if err != nil {
		panic(err)
	}
	user, err := p.GetUserInfo(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", user)
}
