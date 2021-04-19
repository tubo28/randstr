package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"flag"
	"math/rand"
)

var (
	digit, alnum, graph bool
	pattern             string
)

const (
	// O01Il are omitted
	digits  = "23456789"
	alphas  = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	symbols = `!#$%&*/?@`
)

func init() {
	flag.BoolVar(&digit, "digit", false, "Digits")
	flag.BoolVar(&alnum, "alnum", false, "Digits + Latin alphabet (default)")
	flag.BoolVar(&graph, "graph", false, "Digits + Latin alphabet + symbols")
	flag.StringVar(&pattern, "pattern", "XXXXXXXXXXXXXXXX", "Pattern. Each X is replaced with a random character.")

	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		panic("cannot get secure random seed")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func main() {
	flag.Parse()

	var subsets []string
	if digit || alnum || graph {
		subsets = append(subsets, digits)
	}
	if alnum || graph {
		subsets = append(subsets, alphas)
	}
	if graph {
		subsets = append(subsets, symbols)
	}
	if subsets == nil {
		subsets = append(subsets, digits, alphas)
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
