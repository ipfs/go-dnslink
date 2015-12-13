package main

import (
	"fmt"
	"os"

	dnslink "github.com/jbenet/go-dnslink"
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
	if len(args) < 1 {
		fmt.Print(usage)
		return nil
	}

	many := len(args) > 1
	for _, domain := range args {
		if many {
			fmt.Print(domain, ": ")
		}

		result, err := dnslink.Resolve(domain)
		if result != "" {
			fmt.Print(result)
		}
		if err != nil {
			if !many {
				return err
			} else {
				fmt.Print("error: ", err)
			}
		}
		fmt.Println()
	}
	return nil
}
