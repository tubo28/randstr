package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	quiet               bool
	digit, alnum, graph bool
	pattern             string
	length              int
)

const (
	// O01Il are omitted
	digits  = "23456789"
	lower   = "abcdefghijkmnopqrstuvwxyz"
	upper   = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	symbols = `!#$%&*/?@`
)

func init() {
	flag.BoolVar(&quiet, "q", true, "Do not show unimportant error messages")

	flag.BoolVar(&digit, "d", false, "Shorthand for digit")
	flag.BoolVar(&digit, "digit", false, "Digits")

	flag.BoolVar(&alnum, "a", false, "Shorthand for alnum")
	flag.BoolVar(&alnum, "alnum", false, "Digits + Latin alphabet (default)")

	flag.BoolVar(&graph, "g", false, "Shorthand for graph")
	flag.BoolVar(&graph, "graph", false, "Digits + Latin alphabet + symbols")

	flag.StringVar(&pattern, "p", "XXXXXXXXXXXXXXXX", "Shorthand for pattern")
	flag.StringVar(&pattern, "pattern", "XXXXXXXXXXXXXXXX", "Pattern. Each X is replaced with a random character.")

	flag.IntVar(&length, "l", 16, "Shorthand for length")
	flag.IntVar(&length, "length", 16, "Length of output")

	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		fail("cannot get secure random seed")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func main() {
	flag.Parse()

	var p, l bool
	flag.Visit(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, "p") {
			p = true
		}
		if strings.HasPrefix(f.Name, "l") {
			l = true
		}
	})
	if p && l {
		fail("at most one of pattern and length must be set")
	}
	if l {
		pattern = strings.Repeat("X", length)
	}

	var subsets []string
	if digit || alnum || graph {
		subsets = append(subsets, digits)
	}
	if alnum || graph {
		subsets = append(subsets, lower, upper)
	}
	if graph {
		subsets = append(subsets, symbols)
	}
	if subsets == nil {
		subsets = append(subsets, digits, lower, upper)
	}

	var chars []byte
	for _, s := range subsets {
		chars = append(chars, s...)
	}

	pat := []byte(pattern)
	gen := make([]byte, len(pat))
	done := false
	for i := 0; i < 100; i += 1 {
		for j, c := range pat {
			if c == 'X' {
				c = chars[rand.Intn(len(chars))]
			}
			gen[j] = c
		}
		ok := true
		for _, sub := range subsets {
			if !bytes.ContainsAny(gen, sub) {
				ok = false
				break
			}
		}
		if ok {
			done = true
			break
		}
	}
	if !done {
		fail("all generation attempts failed")
	}

	s := string(gen)
	fmt.Println(s)

	d, err := os.UserHomeDir()
	if err != nil {
		if quiet {
			return
		} else {
			fail("failed to detect home directory: " + err.Error())
		}
	}

	f, err := os.OpenFile(filepath.Join(d, ".randstr"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		if quiet {
			return
		} else {
			fail("failed to open generation log: " + err.Error())
		}
	}

	if _, err := fmt.Fprint(f, s, "\n"); err != nil {
		if quiet {
			return
		} else {
			fail("failed to write generation log: " + err.Error())
		}
	}
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
