package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Student adalah model data mahasiswa
type Student struct {
	ID   string // NIP
	Name string // Nama
}

// DatabaseHandler menampung fungsi fungsi terkait query
type DatabaseHandler struct {
	db *sql.DB
}

// View menampung atribut-atribut untuk keperluan rendering html
type View struct {
	Students     []*Student
	Student      Student
	SearchByName string
	ErrorMessage string
}

// Untuk mempermudah penilaian, semua logic ada di file ini
func main() {

	dbHandler := &DatabaseHandler{}
	dbHandler.init()

	templates := template.Must(template.ParseFiles("templates/student.html"))

	// include folder static untuk menyimpan css dan js
	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) //Go looks in the relative static directory first, then matches it to a

	// Tugas ini hanya memerlukan path "/" yang akan menerima post request dan mengembalikan hasilnya.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view := &View{}

		if id := r.FormValue("id"); id != "" {
			view.Student.ID = strings.TrimSpace(id)

			if _, err := strconv.Atoi(view.Student.ID); err != nil {
				view.ErrorMessage = "NIP harus berupa angka.  "
			}
		}
		if name := r.FormValue("name"); name != "" {
			view.Student.Name = strings.TrimSpace(name)

			if len(view.Student.Name) < 3 || len(view.Student.Name) > 100 {
				view.ErrorMessage += "Nama harus terdiri dari minimal 3 karakter dan maksimal 100 karakter (termasuk spasi)."
			}
		}
		if view.Student.ID != "" && view.Student.Name != "" && view.ErrorMessage == "" {
			dbHandler.save(&view.Student)
			view.Student = Student{}
		}

		if searchByName := r.FormValue("search_by_name"); searchByName != "" {
			view.SearchByName = searchByName
			view.Students = dbHandler.search(searchByName)
		} else {
			view.Students = dbHandler.getList()
		}

		if err := templates.ExecuteTemplate(w, "student.html", view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Listening on PORT 8080")
	//Start the web server, browse ke localhost:8080 untuk melihat hasilnya
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func (dh *DatabaseHandler) init() {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Panic("DB_HOST is not set")
	}
	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Panic("DB_PORT is not set")
	}
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Panic("DB_USER is not set")
	}
	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		log.Panic("DB_PASS is not set")
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Panic("DB_NAME is not set")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	var err error
	dh.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}

	// test koneksi db
	err = dh.db.Ping()
	if err != nil {
		log.Panic(err)
	}

	sql := `CREATE TABLE students (
		id VARCHAR(15) PRIMARY KEY,
		name VARCHAR(100)
	);`
	_, err = dh.db.Exec(sql)
	if err != nil {
		log.Print("table students already exist")
	}
}

func (dh *DatabaseHandler) save(s *Student) {
	sql := `INSERT INTO students (id, name) VALUES ($1, $2)`
	_, err := dh.db.Exec(sql, s.ID, s.Name)
	if err != nil {
		log.Panic(err)
	}
}

func (dh *DatabaseHandler) getList() []*Student {
	sql := `SELECT * FROM students`
	students := []*Student{}
	rows, err := dh.db.Query(sql)
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		s := &Student{}
		err := rows.Scan(&s.ID, &s.Name)
		if err != nil {
			log.Fatal(err)
		}

		students = append(students, s)
	}

	return students
}

func (dh *DatabaseHandler) search(name string) []*Student {
	sql := `SELECT * FROM students where name ilike '%' || $1 || '%'`
	students := []*Student{}
	rows, err := dh.db.Query(sql, name)
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		s := &Student{}
		err := rows.Scan(&s.ID, &s.Name)
		if err != nil {
			log.Fatal(err)
		}

		students = append(students, s)
	}

	return students
}
