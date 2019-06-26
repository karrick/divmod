package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		panic("not enough arguments")
	}

	n64, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}

	d64, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}

	n := uint(n64)
	d := uint(d64)

	q, r := divmod_inline(n, d)
	fmt.Println(q, r)

	q, r = divmod_noinline(n, d)
	fmt.Println(q, r)
}

func divmod_inline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}

//go:noinline
func divmod_noinline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}
