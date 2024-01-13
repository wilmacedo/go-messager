package util

import (
	"crypto/rand"
	"math/big"
)

const lettersAndNumbers = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateRandomID() (string, error) {
	size := len(lettersAndNumbers)
	var result string

	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(size)))
		if err != nil {
			return "", err
		}

		result += string(lettersAndNumbers[num.Int64()])
	}

	return result, nil
}
