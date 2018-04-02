package main

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	numberBytes = "1234567890"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomString(length int) string {
	return generateRandom(letterBytes, length)
}

func generateRandomNumbers(length int) string {
	return generateRandom(numberBytes, length)
}

func generateRandom(template string, length int) string {
	result := make([]byte, length)

	for i := range result {
		result[i] = template[rand.Intn(len(template))]
	}

	return string(result)
}
