package dnslink

import (
	"fmt"
	"testing"
)

type mockDNS struct {
	entries map[string][]string
}

func (m *mockDNS) lookupTXT(name string) (txt []string, err error) {
	txt, ok := m.entries[name]
	if !ok {
		return nil, fmt.Errorf("No TXT entry for %s", name)
	}
	return txt, nil
}

func TestDNSLink(t *testing.T) {
	entries := make(map[string][]string)
	entries["ipfs.io"] = []string{"ipfs.io"}
	entries["dnslink"] = []string{"_dnslink.libp2p.io"}
	m := &mockDNS{
		entries: entries,
	}
	if _, err := m.lookupTXT("ipfs.io"); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve("ipfs.io"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.lookupTXT("dnslink"); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve("_dnslink.libp2p.io"); err != nil {
		t.Fatal(err)
	}
}

func TestDnsLinkParsing(t *testing.T) {
	goodEntries := [][]string{
		[]string{"/dns/foo.com", "foo.com", ""},
		[]string{"/dns/foo.com/bar/baz", "foo.com", "/bar/baz"},
		[]string{"/dns/bar.com", "bar.com", ""},
		[]string{"/dns/baz.test.it/bar/baz", "baz.test.it", "/bar/baz"},
	}

	badEntries := []string{
		"/foo/foo.com",
		"/baz/foo.com/bar/baz",
		"/foo.com/bar/baz",
		"foo.com/bar/baz",
		"QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"QmYhE8xgFCjGcz6PHgnvJz5NOTCORRECT",
		"/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"/ipns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/bar",
	}

	for _, e := range goodEntries {
		a, b, err := ParseLinkDomain(e[0])
		if err != nil {
			t.Fatal("expected entry to parse correctly:", e, "got:", err)
		}
		if a != e[1] {
			t.Fatal("expected entry to parse domain correctly:", e[0], e[1], "got:", a)
		}
		if b != e[2] {
			t.Fatal("expected entry to parse rest correctly:", e[0], e[2], "got:", b)
		}
	}

	for _, e := range badEntries {
		_, _, err := ParseLinkDomain(e)
		if err == nil {
			t.Fatal("expected entry parse to fail:", e)
		}
	}
}

func TestDnsEntryParsing(t *testing.T) {
	goodEntries := []string{
		"dnslink=/dns/foo.com",
		"dnslink=/dns/foo.com/bar/baz",
		"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/ipns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo",
		"dnslink=/ipns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/bar",
		"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo/bar/baz",
		"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo",
	}

	badEntries := []string{
		"/dns/foo.com",
		"/dns/foo.com/bar/baz",
		"foo.com",
		"QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", // we dont support this one here.
		"QmYhE8xgFCjGcz6PHgnvJz5NOTCORRECT",
		"quux=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=",
		"dnslink=ipns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/bar",
	}

	for _, e := range goodEntries {
		_, err := ParseTXT(e)
		if err != nil {
			t.Fatal("expected entry to parse correctly:", e, "got:", err)
		}
	}

	for _, e := range badEntries {
		_, err := ParseTXT(e)
		if err == nil {
			t.Fatal("expected entry parse to fail:", e)
		}
	}
}

func newMockDNS() *mockDNS {
	return &mockDNS{
		entries: map[string][]string{
			"foo.com":            []string{"dnslink=/dns/bar.com/foo/f/o/o"},
			"bar.com":            []string{"dnslink=/dns/test.it.baz.com/bar/b/a/r"},
			"test.it.baz.com":    []string{"dnslink=/baz/b/a/z"},
			"ipfs.example.com":   []string{"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD"},
			"dns1.example.com":   []string{"dnslink=/dns/ipfs.example.com"},
			"dns2.example.com":   []string{"dnslink=/dns/dns1.example.com"},
			"equals.example.com": []string{"dnslink=/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/=equals"},
			"loop1.example.com":  []string{"dnslink=/dns/loop2.example.com"},
			"loop2.example.com":  []string{"dnslink=/dns/loop1.example.com"},
			"bad.example.com":    []string{"dnslink="},
			"multi.example.com": []string{
				"some stuff",
				"dnslink=/dns/dns1.example.com",
				"masked dnslink=/dns/example.invalid",
			},
		},
	}
}

func TestResolution(t *testing.T) {
	mock := newMockDNS()
	r := &Resolver{lookupTXT: mock.lookupTXT}
	testResolution(t, r, "foo.com", DefaultDepthLimit, "/baz/b/a/z/bar/b/a/r/foo/f/o/o", nil)
	testResolution(t, r, "bar.com", DefaultDepthLimit, "/baz/b/a/z/bar/b/a/r", nil)
	testResolution(t, r, "test.it.baz.com", 1, "/baz/b/a/z", nil)
	testResolution(t, r, "foo.com", 0, "/dns/foo.com", ErrResolveLimit)
	testResolution(t, r, "foo.com", 1, "/dns/bar.com/foo/f/o/o", ErrResolveLimit)
	testResolution(t, r, "foo.com", 2, "/dns/test.it.baz.com/bar/b/a/r/foo/f/o/o", ErrResolveLimit)
	testResolution(t, r, "bar.com", 0, "/dns/bar.com", ErrResolveLimit)
	testResolution(t, r, "bar.com", 1, "/dns/test.it.baz.com/bar/b/a/r", ErrResolveLimit)
	testResolution(t, r, "test.it.baz.com", 0, "/dns/test.it.baz.com", ErrResolveLimit)
	testResolution(t, r, "ipfs.example.com", DefaultDepthLimit, "/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns1.example.com", DefaultDepthLimit, "/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns1.example.com", 1, "/dns/ipfs.example.com", ErrResolveLimit)
	testResolution(t, r, "dns2.example.com", DefaultDepthLimit, "/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns2.example.com", 1, "/dns/dns1.example.com", ErrResolveLimit)
	testResolution(t, r, "dns2.example.com", 2, "/dns/ipfs.example.com", ErrResolveLimit)
	testResolution(t, r, "multi.example.com", DefaultDepthLimit, "/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "multi.example.com", 1, "/dns/dns1.example.com", ErrResolveLimit)
	testResolution(t, r, "multi.example.com", 2, "/dns/ipfs.example.com", ErrResolveLimit)
	testResolution(t, r, "equals.example.com", DefaultDepthLimit, "/ipfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/=equals", nil)
	testResolution(t, r, "loop1.example.com", 1, "/dns/loop2.example.com", ErrResolveLimit)
	testResolution(t, r, "loop1.example.com", 2, "/dns/loop1.example.com", ErrResolveLimit)
	testResolution(t, r, "loop1.example.com", 3, "/dns/loop2.example.com", ErrResolveLimit)
	testResolution(t, r, "loop1.example.com", DefaultDepthLimit, "/dns/loop1.example.com", ErrResolveLimit)
	testResolution(t, r, "bad.example.com", DefaultDepthLimit, "", ErrInvalidDnslink)
}

func testResolution(t *testing.T, r *Resolver, name string, depth int, expected string, expError error) {
	p, err := r.ResolveN(name, depth)
	if err != expError {
		t.Fatal(fmt.Errorf(
			"Expected %s with a depth of %d to have a '%s' error, but got '%s'",
			name, depth, expError, err))
	}
	if p != expected {
		t.Fatal(fmt.Errorf(
			"%s with depth %d resolved to %s != %s",
			name, depth, p, expected))
	}
}
