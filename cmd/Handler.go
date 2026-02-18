package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"


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

			http.Redirect(w, r, "/Login", http.StatusSeeOther)
			

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
			if err != nil {
				http.Error(w, "DB error", 500)	
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


func handlerChamps(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/ChampPage.html")

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



func handlerItems(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/ItemsPage.html")

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


func handlerSearch(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/SearchPage.html")

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


func handlerNews(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/NewsPage.html")

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

func handlerList(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/ListPage.html")

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

