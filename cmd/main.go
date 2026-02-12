package main

import  (
	
	"net/http"
	"os"
	"log"
	"fmt"
)


func main () {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
			w.Write([]byte("Hello Yboost"))
		})
	

	// Handlers
	http.HandleFunc("/register", handlerRegister)
	http.HandleFunc("/logged",handlerLogged)
	http.HandleFunc("/Profil",handlerProfil)
	http.HandleFunc("/Champs",handlerChamps)
	http.HandleFunc("/Items",handlerItems)
	http.HandleFunc("/Search",handlerSearch)
	http.HandleFunc("/News",handlerNews)
	http.HandleFunc("/List",handlerList)


	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Println("Server running on port", port)
	fmt.Println("http://localhost:8080/logged")
	fmt.Println("http://localhost:8080/register")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}




