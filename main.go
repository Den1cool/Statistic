package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type statistics struct {
	DateT   string
	ViewsT  string
	ClicksT string
	CostT   string
}

type Form struct {
	Alert string
	Date1 string
	Date2 string
}

func save(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=root dbname=db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r.ParseForm()
	http.ServeFile(w, r, "save.html")
	date := r.FormValue("date")
	if date != "" {
		views := r.FormValue("views")
		if views == "" {
			views = "0"
		}
		clicks := r.FormValue("clicks")
		if clicks == "" {
			clicks = "0"
		}
		cost := r.FormValue("cost")
		if cost == "" {
			cost = "0"
		}
		_, err = db.Exec("insert into static VALUES($1, $2, $3, $4)", date, views, clicks, cost)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(date)
		fmt.Println(views)
		fmt.Println(clicks)
		fmt.Println(cost)
	}
}
func show(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=root dbname=db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tmpl, _ := template.ParseFiles("show.html")
	r.ParseForm()
	dat1 := r.FormValue("date1")
	dat2 := r.FormValue("date2")
	fmt.Println(dat1)
	fmt.Println(dat2)
	rows, err := db.Query("select * from static where date >= $1  AND date <=  $2", dat1, dat2)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	Stat := []statistics{}
	for rows.Next() {
		s := statistics{}
		err := rows.Scan(&s.DateT, &s.ViewsT, &s.ClicksT, &s.CostT)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		Stat = append(Stat, s)
	}
	for idx := range Stat {
		Stat[idx].DateT = Stat[idx].DateT[:10]
	}
	tmpl.Execute(w, Stat)
}

func drop(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=root dbname=db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("TRUNCATE static")
	if err != nil {
		log.Fatalln(err)
		fmt.Fprintf(w, "%s", "Произошла ошибка попробуйте удалить данные позже")
	}
	fmt.Fprintf(w, "%s", "Данные успешно удалены")
}
func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func find(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "find.html")
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/find", find)
	mux.HandleFunc("/save", save)
	mux.HandleFunc("/drop", drop)
	mux.HandleFunc("/show", show)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":9000", mux)
}
