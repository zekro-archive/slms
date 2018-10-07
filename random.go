package main

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	fullsetsize   = 95
	setstartindex = 32
)

func getRandChar(charset string, pattern *regexp.Regexp) string {
	filteredset := pattern.FindAllString(charset, len(charset))
	rnd := int(rand.Float32() * float32(len(filteredset)))
	return filteredset[rnd]
}

func createcharset() string {
	var set = ""
	for i := 0; i < fullsetsize; i++ {
		set += string(setstartindex + i)
	}
	return set
}

func RandomString(len int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	res := ""
	pattern := regexp.MustCompile(`[a-zA-Z0-9]`)
	charset := createcharset()
	for i := 0; i < len; i++ {
		res += getRandChar(charset, pattern)
	}
	return res
}
