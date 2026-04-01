package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"
	"fmt"
)

// Return a random item name form the list itemName
func pickRandomItemName () (string, error) {
	if len(itemCards) == 0 {
		return "", fmt.Errorf("itemCards est vide")
	}
	rand.Seed(time.Now().UnixNano())
	return itemCards[rand.Intn(len(itemCards))].Name, nil
}

// Return the daily item of the DB if exist or create it if he doesn't
func GetOrCreateDailyItem () (string, error) {
	today := time.Now().UTC().Format("2006-01-02")

	var itemName string 

	err := db.QueryRow("SELECT item_name FROM daily_item WHERE day = ?", today).Scan(&itemName)
	if err == nil {
		return itemName,nil 
	}
	if err != sql.ErrNoRows {
		return "", err
	}


	candidate, err := pickRandomItemName()
	_, err = db.Exec("INSERT INTO daily_item (day, item_name) VALUES (?,?)",  today, candidate)
	if err == nil {
		return candidate, nil
	}

	err2 := db.QueryRow("SELECT item_name FROM daily_item WHERE day = ?", today).Scan(&itemName)
	if err2 != nil {
		return "", err 
	}
	return itemName, nil 
}

// Check if the guessed item is the same as the daily item
func SameItem (guess string) bool {
	daily, err := GetOrCreateDailyItem()
	if err != nil {
		log.Print("Erreur DB")
	}
	return strings.EqualFold(guess, daily)
}


func IsValidChallengeItem(item Item) bool {
	// 1. Doit être achetable
	if !item.Gold.Purchasable {
		return false
	}

	// 2. Doit exister sur la map SR (map 11)
	if !item.Maps["30"] {
		return false
	}

	// 3. Optionnel : éviter les entrées bizarres sans vrai contenu
	if item.Name == "" {
		return false
	}

	return true
}