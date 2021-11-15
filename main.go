package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"flag"
	"math/rand"
	"strings"
)

var (
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
	flag.BoolVar(&digit, "digit", false, "Digits")
	flag.BoolVar(&alnum, "alnum", false, "Digits + Latin alphabet (default)")
	flag.BoolVar(&graph, "graph", false, "Digits + Latin alphabet + symbols")
	flag.StringVar(&pattern, "pattern", "XXXXXXXXXXXXXXXX", "Pattern. Each X is replaced with a random character.")

	flag.IntVar(&length, "length", 16, "Length of output")
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		panic("cannot get secure random seed")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func main() {
	flag.Parse()

	var p, l bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "p" || f.Name == "pattern" {
			p = true
		}
		if f.Name == "l" || f.Name == "length" {
			l = true
		}
	})
	if p && l {
		panic("at most one of pattern and length must be set")
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
		panic("all generation attempts failed")
	}
	println(string(gen))
}
