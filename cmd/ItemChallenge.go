package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"
	"fmt"
)


// Choose a random name from a list of items and return it
func pickRandomItemName () (string, error) {
	if len(itemCards) == 0 {
		return "", fmt.Errorf("itemCards est vide")
	}

	rand.Seed(time.Now().UnixNano())
	return itemCards[rand.Intn(len(itemCards))].Name, nil
}


// Checks if for today’s date there is an associated item in the daily_item DB and returns it, otherwise it randomly chooses one
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
	if err != nil {
		return "", err
	}

	_, err = db.Exec("INSERT INTO daily_item (day, item_name) VALUES (?,?)", today, candidate)
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
		return false 
	}

	return strings.EqualFold(guess, daily)
}

// Check if the an item is in the league of legends SR and not only in others mode 
func IsValidChallengeItem(item Item) bool {
	if !item.Gold.Purchasable {
		return false
	}

	if item.Maps == nil || !item.Maps["11"] {
		return false
	}

	blacklist := map[string]bool{
		"Lame du gardien":  true,
		"Marteau du gardien": true,
		"Corne du gardien":   true,
		"Orbe du gardien":    true,
		"Couronne du Roi-démon": true, 
		"Poro-Snax":         true,
	}

	if blacklist[item.Name] {
		return false
	}

	if strings.Contains(item.Name, "Bonus") ||
		strings.Contains(item.Name, "Test") ||
		strings.Contains(item.Name, "Dummy") {
		return false
	}

	if item.Gold.Total < 2000 {
		return false
	}

	if item.SpecialRecipe != 0 {
		return false
	}

	return true
}