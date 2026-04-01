package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"
    "regexp"

)

func GetOrCreateDailyChampion() (string, error) {
    today := time.Now().UTC().Format("2006-01-02")

    var championName string

    err := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?", today).Scan(&championName)
    if err == nil {
        return championName, nil
    }
    if err != sql.ErrNoRows {
        return "", err
    }

    candidate := pickRandomChampionName()
    _, err = db.Exec("INSERT INTO daily_champion (day, champion_id) VALUES (?, ?)", today, candidate)
    if err == nil {
        return candidate, nil
    }

    err2 := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?", today).Scan(&championName)
    if err2 != nil {
        return "", err 
    }
    return championName, nil
}

func pickRandomChampionName() string {
	rand.Seed(time.Now().UnixNano())
	return championCards[rand.Intn(len(championCards))].Name
}

func SameChamp(guess string) bool {
    dayli, err := GetOrCreateDailyChampion()
    if err != nil {
        log.Print("Erreur DB")
    }
    return strings.EqualFold(guess, dayli) 
}


    


func HideName(lore string, banned []string) string {
	for _, word := range banned {
		if word == "" {
			continue
		}

		mask := ""
		for i := 0; i < len(word); i++ {
			if word[i] == ' ' || word[i] == '-' || word[i] == '\'' {
				mask += string(word[i])
			} else {
				mask += "*"
			}
		}

		pattern := "(?i)" + regexp.QuoteMeta(word)
		re := regexp.MustCompile(pattern)
		lore = re.ReplaceAllString(lore, mask)
	}

	return lore
}