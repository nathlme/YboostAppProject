package main

import (
	
	"database/sql"
	"math/rand"
	"time"
)

func GetOrCreateDailyChampion() (string, error) {
    today := time.Now().UTC().Format("2006-01-02")

    var championID string
    err := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?", today).Scan(&championID)
    if err == nil {
        return championID, nil
    }
    if err != sql.ErrNoRows {
        return "", err
    }

    candidate := pickRandomChampionName()
    _, err = db.Exec("INSERT INTO daily_champion (day, champion_id) VALUES (?, ?)", today, candidate)
    if err == nil {
        return candidate, nil
    }

    err2 := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?", today).Scan(&championID)
    if err2 != nil {
        return "", err 
    }
    return championID, nil
}

func pickRandomChampionName() string {
	rand.Seed(time.Now().UnixNano())
	return championNames[rand.Intn(len(championNames))]
}