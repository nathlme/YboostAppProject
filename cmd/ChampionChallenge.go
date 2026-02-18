package main

import (
	
	"database/sql"
	"time"
)

func GetOrCreateDailyChampion() (string,error) {
	today := time.Now().UTC().Format("2006-01-02")
	var championID string 
	row := db.QueryRow("SELECT champion_id FROM daily_champion WHERE day = ?",today)

	err := row.Scan(&championID)

	if err == sql.ErrNoRows {
		championID = "Ahri"
		_, erre := db.Exec("INSERT INTO daily_champion (day,champion_id) VALUES (?,?)",today,championID)
		
		if erre != nil {
			return "",erre
		}


			
	}else if err != nil {
		return "",err
	}

	return championID,nil
}