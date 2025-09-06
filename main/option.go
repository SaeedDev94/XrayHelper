package main

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/commands"
	"XrayHelper/main/log"
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var Option struct {
	ConfigFilePath   string `short:"c" long:"config" description:"specify configuration file"`
	VerboseFlag      bool   `short:"v" long:"verbose" description:"show verbose debug information"`
	VersionFlag      bool   `short:"V" long:"version" description:"show current version"`

	Service commands.ServiceCommand `command:"service" description:"control core service"`
	Proxy   commands.ProxyCommand   `command:"proxy" description:"control system proxy"`
}

// LoadOption load Option, the program entry
func LoadOption() (exitCode int) {
	// if no args, show Intro
	if len(os.Args) == 1 {
		fmt.Println(builds.VersionStatement())
		fmt.Println(builds.IntroStatement())
		return
	}
	log.Verbose = &Option.VerboseFlag
	builds.ConfigFilePath = &Option.ConfigFilePath
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		var flagsError *flags.Error
		if errors.As(err, &flagsError) {
			if errors.Is((*flagsError).Type, flags.ErrCommandRequired) {
				if Option.VersionFlag {
					fmt.Println(builds.Version())
					err = nil
				} else {
					exitCode = 127
				}
			} else if errors.Is((*flagsError).Type, flags.ErrHelp) {
				fmt.Println(builds.VersionStatement())
				fmt.Println(err.Error())
				err = nil
			} else {
				exitCode = 126
			}
			log.HandleError(err)
		} else {
			log.HandleError(err)
			exitCode = 1
		}
	}
	return
}
