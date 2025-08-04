package main

import (
	"bookmanager/cmd/bookmanager/api"
	"bookmanager/cmd/bookmanager/commands"
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	apiURL := flag.String("api-url", "http://localhost:8080/api/v1/", "API server URL")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("bookmanager v%s\n", version)
		os.Exit(0)
	}

	client := api.NewAPIClient(*apiURL, *verbose)

	args := flag.Args()
	if len(args) < 1 {
		printHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "book":
		commands.HandleBookCommand(client, args[1:])
	case "collection":
		commands.HandleCollectionCommand(client, args[1:])
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
    fmt.Println(`Usage: bookmanager [global options] <command> [command options]

    Global options:
    --api-url    URL of the API server (default: http://localhost:8080/api/v1)
    --verbose    Enable verbose output
    --version    Show version and exit
    --help       Show help

    Commands:
    book        Manage books
    collection  Manage collections
    help        Shows this help message

    Use 'bookmanager <command> --help' for more information about a command.`)
}
