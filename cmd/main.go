package main

import  (
	
	"net/http"
	"os"
	"log"
)


func main () {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello Yboost 🚀"))
		})
	
	http.HandleFunc("/register", handlerRegister)
		
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}




