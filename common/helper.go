package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var characters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func randSequence(number int) string {
	b := make([]rune, number)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := range b {
		b[i] = characters[r1.Intn(999999)%len(characters)]
	}

	return string(b)
}

func GenerateSalt(length int) string {
	if length < 0 {
		length = 50
	}

	return randSequence(length)
}

func GenerateSecretKey(length int) string {
	if length < 0 {
		length = 30
	}

	return randSequence(length)
}

func GenerateRandomNumberString() string {
	min := 100
	max := 999
	randomNumber := rand.Intn(max-min+1) + min

	return fmt.Sprintf("%d", time.Now().UnixNano()) + strconv.Itoa(randomNumber)
}

func FindElement(slice []string, value string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}

	return false
}
