package main

import (
	"github.com/xiaoyan648/goth"
	"github.com/xiaoyan648/goth/providers/github"
)

func main() {
	goth.UseProviders(
		github.NewProvider("1dc4c54e64e3dd1cd79c", "8555221c37aec8215c2df49b36c5cd3481db67d3"),
	)

	//"xxx/auth/github/callback", {code->token->user}
}
