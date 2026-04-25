package main 


import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"fmt"
	"time"
	"database/sql"
 
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


func GetCurrentUser(r *http.Request) (int, string, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return 0, "", err
	}

	userID, ok := sessions[c.Value]
	if !ok {
		return 0, "", fmt.Errorf("session invalide")
	}

	var pseudo string
	err = db.QueryRow("SELECT pseudo FROM users WHERE id = ?", userID).Scan(&pseudo)
	if err != nil {
		return 0, "", err
	}


	return userID, pseudo, nil
}


func UpdateStreak(userID int) error {
	now := time.Now()
	today := now.Format("2006-01-02")

	var streak int
	var bestStreak int
	var dayPlayed int
	var lastPlayed sql.NullString

	err := db.QueryRow(
		"SELECT streak, best_streak, day_played, last_played FROM users WHERE id = ?", userID, ).Scan(&streak, &bestStreak, &dayPlayed, &lastPlayed)

	if err != nil {
		return err
	}

 	if !lastPlayed.Valid {
		_, err = db.Exec(`UPDATE users SET streak = 1, best_streak = 1, day_played = 1, last_played = ? WHERE id = ?`, today, userID,)
		return err
	}

	lastDate, _ := time.Parse("2006-01-02", lastPlayed.String)
	yesterday := now.AddDate(0, 0, -1)

 
	if lastDate.Format("2006-01-02") == today {
		return nil
	}

 
	if lastDate.Format("2006-01-02") == yesterday.Format("2006-01-02") {
		streak++
	} else {
		streak = 1
	}

	dayPlayed++

	if streak > bestStreak {
		bestStreak = streak
	}

	_, err = db.Exec(
		`UPDATE users SET streak = ?, best_streak = ?, day_played = ?, last_played = ?  WHERE id = ?`, streak, bestStreak, dayPlayed, today, userID, 
	)

	return err
}