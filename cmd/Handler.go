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

			c, err := r.Cookie("session_id")

			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userID, ok := sessions[c.Value] 

			if !ok {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			log.Println("userID:", userID)

			row := db.QueryRow("SELECT pseudo FROM users WHERE id = ?", userID)

			var pseudo string
			
			err = row.Scan(&pseudo)

			if err == sql.ErrNoRows {
				http.Error(w,"identifiants invalide",401)
				return 
			}

			t, err := template.ParseFiles("templates/HomePage.html")

			if err!=nil {
				http.Error(w,"Internal Server Error",500)
				return
			}					
			log.Print(GetOrCreateDailyChampion())

			t.Execute(w, struct {Pseudo string}{Pseudo: pseudo})

	default :
			fmt.Fprint(w,"erreur 405")
	}
}


func handlerProfil(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/ProfilPage.html")

				if err!=nil {
					http.Error(w,"Internal Server Error",500)
					return
				}

			w.Header().Set("Content-Type", "text/html")
			w.Write(contenu)

	default :
			fmt.Fprint(w,"erreur 405")
	}
}


func handlerChamps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	championID, err := GetOrCreateDailyChampion()
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
	lore = HideName(lore, bannedWord)


	data := struct {
		ChampionID string
		ChampionsJSON template.JS
		ChampLore string
	}{
		ChampionID: championID,
		ChampionsJSON: template.JS(championsJSON),
		ChampLore: lore,
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
		ItemName      string
		ItemStats     map[string]float64
		ItemPrice     int
		ItemJSON      template.JS
		ItemStatsJSON template.JS
		ItemComponentsJSON template.JS
	}{
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



func handlerSpell(w http.ResponseWriter, r*http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w,"Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	t, err := template.ParseFiles("templates/SpellPage.html")
	if err != nil {
		http.Error(w, "ERR_PARSE_TEMPLATE_ITEM", http.StatusInternalServerError)
		return 
	}

	data := struct {

	}{

	}

	if err := t.Execute(w, data); err != nil {
		log.Println("jsp item :", err) 
		return 
	}
}




func handlerGuess(w http.ResponseWriter, r*http.Request) {
	switch r.Method {
		case http.MethodPost :
			var req GuessRequest
			var check GuessResponse 

			err := json.NewDecoder(r.Body).Decode(&req)

			if err != nil {
				http.Error(w,"JSON invalide", 400)
				return
			}
			 
			check.Same = SameChamp(req.Champion) 
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

			err := json.NewDecoder(r.Body).Decode(&req)

			if err != nil {
				http.Error(w,"JSON invalide", 400)
				return
			}
			 
			check.Same = SameItem(req.Item) 
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(check)
			
	default :
			http.Error(w,"Erreur",http.StatusMethodNotAllowed)
	}
}