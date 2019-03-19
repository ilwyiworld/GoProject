package main

import (
	"fmt"
	"crypto/sha1"
)

func main() {
	s := "sha1 this string"
	h := sha1.New()
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte slice.
	// The argument to Sum can be used to append to an existing byte slice: it usually isn’t needed.
	bs := h.Sum(nil)
	fmt.Println(s)

	// SHA1 values are often printed in hex,
	// for example in git commits.
	// Use the %x format verb to convert a hash results to a hex string.
	fmt.Printf("%x\n", bs)
}