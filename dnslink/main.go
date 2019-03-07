package main

import (
	"fmt"
	"os"

	dnslink "github.com/ipfs/go-dnslink"
)

var usage = `dnslink - resolve dns links in TXT records

USAGE
    dnslink <domain>

EXAMPLE
    > dnslink blog.ipfs.io
    /ipns/ipfs.io/blog

    > dnslink ipfs.io blog.ipfs.io
    ipfs.io: /ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy
    blog.ipfs.io: /ipns/ipfs.io/blog

    > dnslink foo.bar
    error: lookup foo.bar on 10.0.1.1:53: no such host
`

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 || hasHelp(args) {
		fmt.Print(usage)
		return nil
	}

	if len(args) == 1 {
		return printLink(args[0])
	}
	return printLinks(args)
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

// print a single link
func printLink(domain string) error {
	link, err := dnslink.Resolve(domain)
	if err != nil {
		return err
	}
	fmt.Println(link)
	return nil
}

// print multiple links.
// errors printed as output, and do not fail the entire process.
func printLinks(domains []string) error {
	for _, domain := range domains {
		fmt.Print(domain, ": ")

		result, err := dnslink.Resolve(domain)
		if result != "" {
			fmt.Print(result)
		}
		if err != nil {
			fmt.Print("error: ", err)
		}
		fmt.Println()
	}
	return nil
}
