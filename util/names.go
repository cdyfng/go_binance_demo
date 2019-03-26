package util

import (
		"math/rand"
		"time"
		"fmt"
)

func MakeRandomStrLower(l int) string {
	chars := "12345abcdefghijklmnopqrstuvwxyz"
	clen := float64(len(chars))
	res := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		rfi := int(clen * rand.Float64())
		res += fmt.Sprintf("%c", chars[rfi])
	}

	return res
}

