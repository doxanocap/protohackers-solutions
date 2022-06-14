package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "eldoseldos"
	dbname   = "grindset"
)

type Users struct {
	Id       int
	Email    string
	Password string
}

func main() {
	handleFunc()
}

func handleFunc() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/sign-in", signIn)
	r.HandleFunc("/welcome", welcomePage)
	r.HandleFunc("/wrong", wrongPassword)
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	_ = http.ListenAndServe(":2005", nil)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	temp, err := template.ParseFiles("templates/signin.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}

	email := string(r.FormValue("email"))
	password := string(r.FormValue("password"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	res, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var user Users
		err := res.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			panic(err)
		}
		if user.Password == password {
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/wrong", http.StatusSeeOther)
		}
	}
	_ = temp.ExecuteTemplate(w, "signin", user)
}

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	_ = temp.ExecuteTemplate(w, "index", nil)
}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/welcome.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	_ = temp.ExecuteTemplate(w, "index", user)
}

func wrongPassword(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "wrong password")
}
