package main

import (
	"github.com/alia5/urlshort/cmd"
	"github.com/alia5/urlshort/serve"
	"github.com/alia5/urlshort/storinator"
	"github.com/alia5/urlshort/urlshort"
)

func main() {
	command := cmd.ParseCmd()
	storinator.ConnectDB(cmd.CLI.DB)
	switch command {
	// can I just say "FUCK `gofmt`s switch indentation! srlsy"?
	case cmd.CmdServe:
		urlshort.SetBaseUrl(cmd.CLI.BaseUrl)
		serve.Run(serve.ServeOptions{
			Port:  cmd.CLI.Port,
			Debug: cmd.CLI.Debug,
		})
	}
}
