package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
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

			log.Println("Hash:", hash)

			log.Println("Nouvelle inscription :", email,len(password),pseudo)
			http.Redirect(w, r, "/logged", http.StatusSeeOther)
			return
			
		default :
			fmt.Fprint(w,"erreur 405")
	} 

	
}





func handlerLogged(w http.ResponseWriter, r*http.Request){
	switch r.Method {
		case http.MethodGet :
			contenu, err := os.ReadFile("templates/HomePage.html")

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
