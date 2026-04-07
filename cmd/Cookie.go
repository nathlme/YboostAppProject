package main 


import (
	"crypto/rand"
	"encoding/hex"
)


// Generates a secure random session identifier 
func generateSessionID() (string, error) {
	B := make([]byte, 32)

	_, err := rand.Read(B)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(B), nil
}


var sessions = map[string]int{}