package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type App struct {
	Port   string
	Router *http.ServeMux
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("error in connection with database %s", err)
	}

	log.Println("Connected to database")

	app := App{
		Port:   ":8080",
		Router: NewRouterWithDeps(db),
	}

	log.Println("Rock Paper Scissors running on Port ", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, app.Router))
}
