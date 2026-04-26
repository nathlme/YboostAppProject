package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"golang.org/x/crypto/bcrypt"
)



func handlerRegister(w http.ResponseWriter, r*http.Request ){
	switch r.Method {
		case http.MethodGet :

			contenu, err := os.ReadFile("templates/RegisterPage.html")
			if err!=nil {
				http.Error(w,"Internal Server Error",500)
				return
			}
			
			w.Header().Set("Content-Type", "text/html")
			w.Write(contenu)
			
		case http.MethodPost :

			err := r.ParseForm()
			if err != nil {
				http.Error(w,"Parse error",400)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")
			pseudo := r.FormValue("pseudo")
	
			message := validateRegister(email,pseudo,password)

			if message != "" {
				http.Error(w,message,400)
				return
			}

			hash, err := PasswordHash([]byte(password))
			if err != nil {
				http.Error(w,"Erreur de Hash", 400)
				return
			}

			rest, err := db.Exec("INSERT INTO users (email, pseudo, password_hash) VALUES (?,?,?) ",email, pseudo, hash)

			id, _ := rest.LastInsertId()
			log.Println("Nouvel ID:", id)

			if err != nil {
				http.Error(w, "DB error", 500)
				return 
			}	

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			

			log.Println("Hash:", hash)

			log.Println("Nouvelle inscription :", email,len(password),pseudo)
			
			return
			
		default :
			fmt.Fprint(w,"erreur 405")
	} 
}



func handlerLogin(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/LoginPage.html")

				if err!=nil {
					http.Error(w,"Internal Server Error",500)
					return
				}

			w.Header().Set("Content-Type", "text/html")
			w.Write(contenu)

		case http.MethodPost :

			err := r.ParseForm()
			if err != nil {
				http.Error(w,"Parse error", 400)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")
			row := db.QueryRow("SELECT id, password_hash, pseudo FROM users WHERE email = ?",email)
			
			var id int
			var hash string
			var pseudo string

			err = row.Scan(&id, &hash, &pseudo)

			if err == sql.ErrNoRows {
				http.Error(w,"identifiants invalide",401)
				return 
			}
		
			if err != nil {
				http.Error(w, "DB error", 500)
				return 
			}	
			
			res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

			if res != nil {
				http.Error(w,"identifiants incorrect !",401)
				return

			}else {
				log.Println("email : ",email,", pseudo : ",pseudo)

				sessionID, err := generateSessionID()

				if err != nil {
					http.Error(w,"Cookie issues",500)
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
			
				http.SetCookie(w,&cookie)

				http.Redirect(w, r, "/logged", http.StatusSeeOther)
				return 
			}			
		default :
				fmt.Fprint(w,"erreur 405")
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
			log.Print(GetOrCreateDailyChampion("daily_champion"))

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
	
	var streak 		int 
	var bestStreak 	int 
	var lastPlayed 	string
	var dayPlayed	int
	var inscriptionDate string

	err2 := db.QueryRow(
		"SELECT created_at, streak, best_streak, day_played, last_played FROM users WHERE id = ?", userID, ).Scan(&inscriptionDate, &streak, &bestStreak, &dayPlayed, &lastPlayed ) 
	if err2 != nil {
		return 
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
		LastPlayed: lastPlayed,
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
				log.Println("correct")
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
				log.Println("correct")
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
				log.Println("correct")
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
 
 
