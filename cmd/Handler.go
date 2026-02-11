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
			fmt.Fprint(w, "OK")

			
		default :
			fmt.Fprint(w,"erreur 405")
	} 

	
}