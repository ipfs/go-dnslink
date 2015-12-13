# dnslink resolution in go-ipfs

Package dnslink implements a dns link resolver. dnslink is a basic
standard for placing traversable links in dns itself. See dnslink.info

A dnslink is a path link in a dns TXT record, like this:

```
dnslink=/ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy
```

For example:

```
> dig TXT ipfs.io
ipfs.io.  120   IN  TXT  dnslink=/ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy
```

This package eases resolving and working with thse dns links. For example:

```go
import (
  dnslink "github.com/jbenet/go-dnslink"
)

link, err := dnslink.Resolve("ipfs.io")
// link = "/ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy"
```

It even supports recursive resolution. Suppose you have three domains with
dnslink records like these:

```
> dig TXT foo.com
foo.com.  120   IN  TXT  dnslink=/dns/bar.com/f/o/o
> dig TXT bar.com
bar.com.  120   IN  TXT  dnslink=/dns/long.test.baz.it/b/a/r
> dig TXT long.test.baz.it
long.test.baz.it.  120   IN  TXT  dnslink=/b/a/z
```

Expect these resolutions:

```go
dnslink.ResolveN("long.test.baz.it", 0) // "/dns/long.test.baz.it"
dnslink.Resolve("long.test.baz.it")     // "/b/a/z"

dnslink.ResolveN("bar.com", 1)          // "/dns/long.test.baz.it/b/a/r"
dnslink.Resolve("bar.com")              // "/b/a/z/b/a/r"

dnslink.ResolveN("foo.com", 1)          // "/dns/bar.com/f/o/o/"
dnslink.ResolveN("foo.com", 2)          // "/dns/long.test.baz.it/b/a/r/f/o/o/"
dnslink.Resolve("foo.com")              // "/b/a/z/b/a/r/f/o/o"
```

## Usage

### As a library

```go
import (
  log
  fmt

  dnslink "github.com/jbenet/go-dnslink"
)

func main() {
  link, err := dnslink.Resolve("ipfs.io")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(link) // string path
}
```

### As a commandline tool

Check out [the commandline tool](dnslink/), which works like this:

```sh
> dnslink ipfs.io
/ipfs/QmR7tiySn6vFHcEjBeZNtYGAFh735PJHfEMdVEycj9jAPy
```
