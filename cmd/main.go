package main

import  (
	
	"net/http"
	"os"
	"log"
	"fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
)

var db *sql.DB


func main () {
	if err := LoadChampionCards(); err != nil {
		log.Fatal(err)
	}
	
	LoadItemCard()
	GetVersion()
	raw := os.Getenv("SCALINGO_MYSQL_URL")
	fmt.Println("RAW =", raw)

	u,err := url.Parse(os.Getenv("SCALINGO_MYSQL_URL"))
	
	if err != nil {
		log.Fatal(err)
		return
	}
	if u.User == nil {
		log.Println("URL sans identifiants")
		return 
	}

	user := u.User.Username()
	pass,_ := u.User.Password()
	tunnel := os.Getenv("MYSQL_TUNNEL_HOSTPORT")

	var host string

	if tunnel != "" {
		host = tunnel
	} else {
		host = u.Host
	}
	
	dbName := strings.TrimPrefix(u.Path,"/") 


	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?",
    user,
    pass,
    host,
    dbName,
	)
	

	db, err = sql.Open("mysql",dsn)
	if err != nil {
		log.Fatal(err)
		return 
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("DB connected")



	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
			w.Write([]byte("Hello Yboost"))
		})
	

	// Handlers
	http.HandleFunc("/register", handlerRegister)
	http.HandleFunc("/login",handlerLogin)
	http.HandleFunc("/logged",handlerLogged)
	http.HandleFunc("/Profil",handlerProfil)
	http.HandleFunc("/Champs",handlerChamps)
	http.HandleFunc("/Items",handlerItems)
	http.HandleFunc("/Spells",handlerSpell)
	http.HandleFunc("/guess", handlerGuess)
	http.HandleFunc("/guess-item", handlerGuessItem)


	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}
	
	
	fmt.Println("http://localhost:8080/logged")
	fmt.Println("http://localhost:8080/register")
	fmt.Println("http://localhost:8080/login")
	log.Println(len(championCards))
	log.Fatal(http.ListenAndServe(":"+port, nil))

	
}




