package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func confirm() bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure? [y/N]: ")
	t, err := r.ReadString('\n')
	if err != nil {
		return false
	}

	t = strings.ToLower(strings.TrimSpace(t))
	if t != "y" {
		return false
	}

	return true
}

const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.,"

// RandString generates a random string suitable for passwords.
func RandString(size int) string {
	var s strings.Builder
	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(validChars))))
		if err != nil {
			return s.String()
		}

		s.WriteByte(validChars[n.Int64()])
	}
	return s.String()
}
