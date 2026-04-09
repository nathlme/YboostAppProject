package main

import (
	"math/rand"
	"database/sql"
	"time"
	
)

func GetOrCreateDailySpell() (string,string,error) {
	today := time.Now().UTC().Format("2006-01-02")

	var championID string
	var spellSlot string

	err := db.QueryRow("SELECT champion_id, spell_slot FROM daily_spell WHERE day = ?", today).Scan(&championID,&spellSlot)

	if err != nil {
		return "", "", err 
	}

	if err != sql.ErrNoRows {
		return "", "", err
	}

	championID, err = pickRandomChampionID()
	if err != nil {
		return "", "", err
	}

	slots := []string{"Q","W","E","R"}
	spellSlot = slots[rand.Intn(len(slots))]

	_, err = db.Exec("INSERT INTO daily_spell (day, champion_id, spell_slot) VALUES (?,?,?)", today, championID, spellSlot)
	if err != nil {
		return "", "", err 
	}

	return championID, spellSlot, nil 
};



func slotToIndex(slot string) int {
	switch slot {
	case "Q":
		return 0
	case "W":
		return 1
	case "E":
		return 2
	case "R":
		return 3
	default:
		return -1
	}
}


func GetDailySpell() (*DailySpell, error) {
	championID, spellSlot, err := GetOrCreateDailySpell()

	url := 
}