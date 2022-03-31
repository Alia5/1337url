package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/alia5/urlshort/storinator"
)

var CLI struct {
	// TODO move to serve command
	Port    int16  `short:"p" long:"port" help:"Port to listen on" default:"7080"`
	BaseUrl string `short:"b" long:"baseurl" help:"Base url for short urls" default:"https://localhost:7080"`
	Debug   bool   `long:"debug" help:"Enable debug mode"`

	//
	DB storinator.DBSettings `prefix:"db." embed:""`
}

type Cmd uint

const (
	CmdServe Cmd = iota
	CmdAdd
	CmdDelete
	CmdList
	CmdClear
)

func ParseCmd() Cmd {
	ctx := kong.Parse(&CLI, kong.Configuration(kong.JSON, "config.json", "config.local.json"))
	switch ctx.Command() {
	// TODO: commands
	// serve
	// list
	// clear
	// delete
	// add
	default:
		return CmdServe
		// panic("no command")
	}

}
