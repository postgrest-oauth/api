package main

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	numberBytes = "23456789"
	randNumbers = "numbers"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomString(randType string, length int) string {
	result := make([]byte, length)

	templateString := letterBytes
	if randType == randNumbers {
		templateString = numberBytes
	}

	for i := range result {
		result[i] = templateString[rand.Intn(len(templateString))]
	}

	return string(result)
}
