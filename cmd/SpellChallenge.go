package main

import (
	"math/rand"
	"database/sql"
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	
)

func GetOrCreateDailySpell() (string,string,error) {
	today := time.Now().UTC().Format("2006-01-02")

	var championID string
	var spellSlot string

	err := db.QueryRow("SELECT champion_id, spell_slot FROM daily_spell WHERE day = ?", today).Scan(&championID,&spellSlot)

	if err == nil {
	return championID, spellSlot, nil
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
	if err != nil {
		return nil,err 
	}

	version := GetVersion()
	url := "https://ddragon.leagueoflegends.com/cdn/" + version + "/data/fr_FR/champion/" + championID + ".json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur api riot : %s", resp.Status)
	}

	var result ChampionSpellResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	champion, ok := result.Data[championID]
	if !ok {
		return nil, fmt.Errorf("champion %s introuvable dans la réponse", championID)
	}


	index := slotToIndex(spellSlot)
	if index == -1 {
		return nil, fmt.Errorf("slot invalide : %s", spellSlot)
	}

	if index >= len(champion.Spells) {
		return nil, fmt.Errorf("index spell hors limite")
	}

	spell := champion.Spells[index]

	card := &DailySpell{
		ChampionID:   champion.ID,
		ChampionName: champion.Name,
		SpellSlot:    spellSlot,
		SpellID:      spell.ID,
		SpellName:    spell.Name,
		Description:  spell.Description,
		ImageURL:     "https://ddragon.leagueoflegends.com/cdn/" + version + "/img/spell/" + spell.Image.Full,
	}

	return card, nil

}