package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-sql-driver/mysql"
)

type AuthPageData struct {
	ErrorMessage string
	Email        string
	Pseudo       string
}

func renderRegisterPage(w http.ResponseWriter, errorMessage, email, pseudo string) {
	t, err := template.ParseFiles("templates/RegisterPage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := AuthPageData{
		ErrorMessage: errorMessage,
		Email:        email,
		Pseudo:       pseudo,
	}

	t.Execute(w, data)
}


func handlerRegister(w http.ResponseWriter, r*http.Request ){
	switch r.Method {
		case http.MethodGet :

			t, err := template.ParseFiles("templates/RegisterPage.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			
			t.Execute(w, AuthPageData{})

		case http.MethodPost :

			err := r.ParseForm()
			if err != nil {
				renderRegisterPage(w, "Erreur de formulaire", "", "")	
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")
			pseudo := r.FormValue("pseudo")
	
			message := validateRegister(email, pseudo, password)
			if message != "" {
				renderRegisterPage(w, message, email, pseudo)
				return
			}

			hash, err := PasswordHash([]byte(password))
			if err != nil {
				renderRegisterPage(w, "Erreur lors du hash du mot de passe", email, pseudo)
				return
			}

			_, err = db.Exec(
				"INSERT INTO users (email, pseudo, password_hash) VALUES (?,?,?)",
				email, pseudo, hash,
			)
			if err != nil {
				if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
					renderRegisterPage(w, "Email ou pseudo déjà utilisé", email, pseudo)
					return
				}

				renderRegisterPage(w, "Impossible de créer le compte", email, pseudo)
				return
			}	

			http.Redirect(w, r, "/login", http.StatusSeeOther)
				
		default :
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	} 
}


func renderLoginPage(w http.ResponseWriter, errorMessage, email string) {
	t, err := template.ParseFiles("templates/LoginPage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := AuthPageData{
		ErrorMessage: errorMessage,
		Email:        email,
	}

	t.Execute(w, data)
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("templates/LoginPage.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, AuthPageData{})

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			renderLoginPage(w, "Erreur de formulaire", "")
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		row := db.QueryRow("SELECT id, password_hash, pseudo FROM users WHERE email = ?", email)

		var id int
		var hash string
		var pseudo string

		err = row.Scan(&id, &hash, &pseudo)

		if err == sql.ErrNoRows {
			renderLoginPage(w, "Identifiants incorrects", email)
			return
		}

		if err != nil {
			renderLoginPage(w, "Erreur base de données", email)
			return
		}

		res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if res != nil {
			renderLoginPage(w, "Identifiants incorrects", email)
			return
		}

		sessionID, err := generateSessionID()
		if err != nil {
			renderLoginPage(w, "Erreur de session", email)
			return
		}

		sessions[sessionID] = id

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/logged", http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}



func handlerLogged(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			 
			userID, pseudo, err := GetCurrentUser(r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			t, err := template.ParseFiles("templates/HomePage.html")

			if err!=nil {
				http.Error(w,"Internal Server Error",500)
				return
			}

			t.Execute(w, struct {
				Pseudo string
				UserID int 
			}{
				Pseudo: pseudo,
				UserID: userID,
			})

	default :
			fmt.Fprint(w,"erreur 405")
	}
}


func handlerProfil(w http.ResponseWriter, r*http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID, pseudo, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates/ProfilPage.html")
	if err != nil {
		http.Error(w, "ERR_PARSE_TEMPLATE", http.StatusInternalServerError)
		return
	}
		
	err2 := ResetStreak(userID)
	if err2 != nil{
		http.Error(w, "Erreur serveur", 500)
		return
	}

	var streak 		int 
	var bestStreak 	int 
	var lastPlayed 	sql.NullString
	var dayPlayed	int
	var inscriptionDate string

	err3 := db.QueryRow(
		"SELECT created_at, streak, best_streak, day_played, last_played FROM users WHERE id = ?", userID, ).Scan(&inscriptionDate, &streak, &bestStreak, &dayPlayed, &lastPlayed ) 
	if err3 != nil {
		log.Println("Erreur SQL profil :", err3)
		http.Error(w, "Erreur serveur", 500)
		return
	}

	lp := ""
	if lastPlayed.Valid {
		lp = lastPlayed.String
	}

	data := struct {
		Pseudo string
		UserID int 
		Streak int 
		BestStreak int 
		LastPlayed string
		DayPlayed int
		InscriptionDate string
	}{
		Pseudo: pseudo,
		UserID: userID, 
		Streak: streak,
		BestStreak: bestStreak,
		LastPlayed: lp,
		DayPlayed: dayPlayed,
		InscriptionDate: inscriptionDate,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println("jsp :", err)
		return
	}	

}


func handlerChamps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	
	userID, pseudo, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	championID, err := GetOrCreateDailyChampion("daily_champion")
	if err != nil {
		http.Error(w, "ERR_GET_DAILY", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/ChampPage.html")
	if err != nil {
		http.Error(w, "ERR_PARSE_TEMPLATE", http.StatusInternalServerError)
		return
	}

	championsJSON, err := json.Marshal(championCards)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	name, title, lore, err := GetFullLore()
	if err != nil {
		log.Println("json error:", err)
		return
	}

	bannedWord := []string{name, title}
	HiddenLore := HideName(lore, bannedWord)


	data := struct {
		Pseudo string
		UserID int
		ChampionID string
		ChampionsJSON template.JS
		ChampLore string
		FullLore string
	}{
		Pseudo: pseudo,
		UserID: userID,
		ChampionID: championID,
		ChampionsJSON: template.JS(championsJSON),
		ChampLore: HiddenLore,
		FullLore: lore,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println("jsp :", err)
		return
	}	
}



func handlerItems(w http.ResponseWriter, r*http.Request){ 
	
	if r.Method != http.MethodGet {
		http.Error(w,"Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID, pseudo, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	itemName, err := GetOrCreateDailyItem()
	if err != nil {
		http.Error(w, "ERR_GET_DAILY_ITEM", http.StatusInternalServerError)
		return 
	}

	t, err := template.ParseFiles("templates/ItemsPage.html")
	if err != nil {
		http.Error(w, "ERR_PARSE_TEMPLATE_ITEM", http.StatusInternalServerError)
		return 
	}

	itemJSON, err := json.Marshal(itemCards)
	if err != nil {
		log.Println("json error:", err)
		return 
	}

	var dailyCard ItemCard
	finded := false 

	for _, card := range itemCards {
		if card .Name == itemName {
			dailyCard.Name = card.Name 
			dailyCard.ImageId = card.ImageId
			dailyCard.ImageUrl = card.ImageUrl
			dailyCard.Stats = card.Stats
			dailyCard.Price = card.Price
			dailyCard.SpecialRecipe = card.SpecialRecipe
			dailyCard.From = card.From

			finded = true 
		}
	}
	
	if !finded {
	http.Error(w, "ERR_DAILY_ITEM_NOT_FOUND_IN_ITEMCARDS", http.StatusInternalServerError)
	return
}

	itemStatsJSON, err := json.Marshal(dailyCard.Stats)
	if err != nil {
		http.Error(w, "ERR_MARSHAL_ITEM_STATS", http.StatusInternalServerError)
		return
	}

	dico, err := GetItem()
	if err != nil {
		http.Error(w, "ERR_GET_ITEMS_DICT", http.StatusInternalServerError)
		return
	}

	dailyComponents := GetItemComponents(dailyCard.From, dico)

	componentsJSON, err := json.Marshal(dailyComponents)
	if err != nil {
		http.Error(w, "ERR_MARSHAL_COMPONENTS", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pseudo string
		UserID int
		ItemName      string
		ItemStats     map[string]float64
		ItemPrice     int
		ItemJSON      template.JS
		ItemStatsJSON template.JS
		ItemComponentsJSON template.JS
	}{
		Pseudo: pseudo,
		UserID: userID,
		ItemName:      dailyCard.Name,
		ItemStats:     dailyCard.Stats,
		ItemPrice:     dailyCard.Price,
		ItemJSON:      template.JS(itemJSON),
		ItemStatsJSON: template.JS(itemStatsJSON),
		ItemComponentsJSON: template.JS(componentsJSON),
	}

	if err := t.Execute(w, data); err != nil {
		log.Println("jsp item :", err) 
		return 
	}
}


func handlerSpell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID, pseudo, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	spell, err := GetDailySpell()
	if err != nil {
		http.Error(w, "Erreur chargement spell", http.StatusInternalServerError)
		log.Println("spell error:", err)
		return
	}
	 
	championsJSON, err := json.Marshal(championCards)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	bannedWord := []string{spell.ChampionName}
	correctDescription := CleanTooltip(spell.Description)
	description := HideName(correctDescription, bannedWord) 

	t, err := template.ParseFiles("templates/SpellPage.html")
	if err != nil {
		http.Error(w, "ERR_PARSE_TEMPLATE_SPELL", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pseudo string
		UserID int
		Description string
		ChampionsJSON template.JS
		SpellSlot	string
		SpellImg 	string 
		SpellName    string
	}{
		Pseudo: pseudo,
		UserID: userID,
		Description: description,
		ChampionsJSON: template.JS(championsJSON),
		SpellSlot: spell.SpellSlot,
		SpellImg: spell.ImageURL,
		SpellName: spell.SpellName,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println("template spell:", err)
		return
	}
	
}



func handlerGuess(w http.ResponseWriter, r*http.Request) {
	switch r.Method {
		case http.MethodPost :
			var req GuessRequest
			var check GuessResponse 

			userID,_,err := GetCurrentUser(r) 
			if err != nil {
				http.Error(w, "Non autorisé", http.StatusUnauthorized)
				return
			}

			err2 := json.NewDecoder(r.Body).Decode(&req)

			if err2 != nil {
				http.Error(w,"JSON invalide", 400)
				return
			}
			 
			check.Same = SameChamp(req.Champion,"daily_champion") 
			if check.Same {
				UpdateStreak(userID)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(check)
			
	default :
			http.Error(w,"Erreur",http.StatusMethodNotAllowed)
	}
}

func handlerGuessItem(w http.ResponseWriter, r*http.Request) {
	switch r.Method {
		case http.MethodPost :
			var req ItemGuessRequest
			var check GuessResponse 

			userID,_,err := GetCurrentUser(r) 
			if err != nil {
				http.Error(w, "Non autorisé", http.StatusUnauthorized)
				return
			}

			err2 := json.NewDecoder(r.Body).Decode(&req)

			if err2 != nil {
				http.Error(w,"JSON invalide", 400)
				return
			}
			 
			check.Same = SameItem(req.Item) 

			if check.Same {
				UpdateStreak(userID)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(check)
			
	default :
			http.Error(w,"Erreur",http.StatusMethodNotAllowed)
	}
}

func handlerGuessSpell(w http.ResponseWriter, r*http.Request) {
	switch r.Method {
		case http.MethodPost :
			var req GuessRequest
			var check GuessResponse 

			userID,_,err := GetCurrentUser(r)
			if err != nil {
				http.Error(w, "Non autorisé", http.StatusUnauthorized)
				return
			}

			err2:= json.NewDecoder(r.Body).Decode(&req)

			if err2 != nil {
				http.Error(w,"JSON invalide", 400)
				return
			}
			 
			check.Same = SameChamp(req.Champion, "daily_spell") 

			if check.Same {
				UpdateStreak(userID)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(check)
			
	default :
			http.Error(w,"Erreur",http.StatusMethodNotAllowed)
	}
}


func handlerLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	
	c, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
 
 
