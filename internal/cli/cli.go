package cli

import (
	"fmt"
	"os"

	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/tessera/internal/buildinfo"
)

func printUsage() {
	fmt.Printf(`tessera %s

Usage:
  tessera <command> [flags]

Commands:
  version                        Print the current version
  help                           Print this help message
  service-file                   Generate a systemd service file

`, buildinfo.Version)
}

func DispatchCommands(args []string) {
	if len(args) == 1 {
		return
	}
	cmd := args[1]
	switch cmd {
	case "help":
		printUsage()
		os.Exit(0)
	case "setup-env":
		os.Exit(setupEnv())
	case "version":
		fmt.Println(buildinfo.AppName, buildinfo.Version)
		os.Exit(0)
	case "service-file":
		goutils.GenerateServiceFile(" Centralized username/password auth API with a Go client library and admin dashboard ")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", cmd)
		printUsage()
		os.Exit(1)
	}
}
