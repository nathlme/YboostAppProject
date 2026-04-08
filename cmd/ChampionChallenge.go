package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"
    "regexp"
	"unicode"

)

// Checks if for today’s date there is an associated champion and returns it, otherwise it randomly chooses one
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

    candidate := pickRandomChampionID()

    _, err = db.Exec("INSERT INTO daily_champion (day, champion_id) VALUES (?, ?)", today, candidate)
    if err == nil {
        return candidate, nil
    }

    err2 := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?", today).Scan(&championID)
    if err2 != nil {
        return "", err2
    }

    return championID, nil
}


// Choose a random name from a list of champions and return it
func pickRandomChampionID() string {
    if len(championCards) == 0 {
		return ""
	}
    
	rand.Seed(time.Now().UnixNano())
	return championCards[rand.Intn(len(championCards))].ID
}


// Check if the guessed champion is the same as the daily champion 
func SameChamp(guess string) bool {
	dailyID, err := GetOrCreateDailyChampion()
	if err != nil {
		log.Println("Erreur DB:", err)
		return false
	}

	normalizedGuess := normalizeChampionName(guess)

	for _, champ := range championCards {
		if normalizeChampionName(champ.Name) == normalizedGuess {
			return champ.ID == dailyID
		}
	}

	return false
}


// Replace special characters in a string
func normalizeChampionName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\'', '’', '-', ' ':
			continue
		case 'é', 'è', 'ê', 'ë':
			r = 'e'
		case 'à', 'â', 'ä':
			r = 'a'
		case 'î', 'ï':
			r = 'i'
		case 'ô', 'ö':
			r = 'o'
		case 'ù', 'û', 'ü':
			r = 'u'
		case 'ç':
			r = 'c'
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}


// Take a sentence and a list of banned word and replace them with "*" in a new sentence
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