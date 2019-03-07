# dnslink - resolve dns links in TXT records

This is a simple commandline tool to resolve dnslink records. It is built with the [go-dnslink](../) package.

For more information about dnslink, see

- This note: https://github.com/jbenet/random-ideas/issues/28

## Install

Compile with Go

```sh
go get -u github.com/ipfs/go-dnslink/dnslink
```

## Usage

```
> dnslink --help
dnslink - resolve dns links in TXT records

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
```

## Examples

Resolve a single domain

```sh
> dnslink blog.ipfs.io
/ipns/ipfs.io/blog
```

Resolve multiple domains

```sh
> dnslink ipfs.io blog.ipfs.io
ipfs.io: /ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy
blog.ipfs.io: /ipns/ipfs.io/blog
```

Error handling

```sh
> dnslink foo.bar
error: lookup foo.bar on 10.0.1.1:53: no such host
```
