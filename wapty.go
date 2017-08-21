package main

import (
	"fmt"
	"os"

	"github.com/empijei/wapty/common"
	"github.com/empijei/wapty/decode"
	"github.com/empijei/wapty/help"
	"github.com/empijei/wapty/intercept"
	"github.com/empijei/wapty/mocksy"
	"github.com/empijei/wapty/ui"
)

const banner = `
                     _         
__      ____ _ _ __ | |_ _   _ 
\ \ /\ / / _' | '_ \| __| | | |
 \ V  V / (_| | |_) | |_| |_| |
  \_/\_/ \__,_| .__/ \__|\__, |
              |_|        |___/ 

`

var (
	//Version is taken by the build flags, represent current version as
	//<major>.<minor>.<patch>
	Version string

	//Commit is the output of `git rev-parse HEAD` at the moment of the build
	Commit string
)

var CmdProxy = &common.Command{
	Name:      "proxy",
	Run:       proxyMain,
	UsageLine: "proxy",
	Short:     "work as a proxy",
	Long:      "",
}

var CmdVersion = &common.Command{
	Name: "version",
	Run: func(_ ...string) {
		// Setup fallback version and commit in case wapty wasn't "properly" compiled
		if len(Version) == 0 {
			Version = "Unknown"
		}
		if len(Commit) == 0 {
			Commit = "Unknown"
		}
		fmt.Printf("Version: %s\nCommit: %s\n", Version, Commit)
	},
	UsageLine: "version",
	Short:     "print version and exit",
	Long:      "print version and exit",
}

func init() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	common.WaptyCommands = []*common.Command{
		decode.CmdDecode,
		CmdProxy,
		mocksy.CmdMocksy,
		CmdVersion,
		help.CmdHelp,
	}
}

func main() {
	fmt.Println(banner)
	if len(os.Args) > 1 {
		//read the first argument
		directive := os.Args[1]
		if len(os.Args) > 2 {
			//shift parameters left, but keep argv[0]
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			os.Args = os.Args[:1]
		}
		invokeMain(directive)
	} else {
		proxyMain()
	}
}

func proxyMain(_ ...string) {
	go ui.MainLoop()
	intercept.MainLoop()
}

func invokeMain(s string) {
	command, err := common.FindCommand(s)
	if err == nil {
		command.Flag.Usage = command.Usage
		//TODO handle this error
		_ = command.Flag.Parse(os.Args[1:])
		command.Run(command.Flag.Args()...)
		return
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		fmt.Fprintln(os.Stderr, "Available commands are:\n")
		for _, cmd := range common.WaptyCommands {
			fmt.Fprintln(os.Stderr, "\t"+cmd.Name+"\n\t\t"+cmd.Short)
		}
		fmt.Fprintln(os.Stderr, "\nDefault command is: proxy")
	}
}
