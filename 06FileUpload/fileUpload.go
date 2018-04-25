package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Title string
	Body  []byte
}

// global variables
var (
	gTplView   = template.Must(template.New("view").ParseFiles("./html/base.html", "./html/view.html", "./html/header.html", "./html/footer.html"))
	gTplEdit   = template.Must(template.New("edit").ParseFiles("./html/base.html", "./html/edit.html", "./html/header.html", "./html/footer.html"))
	gTplUpload = template.Must(template.New("upload").ParseFiles("./html/base.html", "./html/upload.html", "./html/header.html", "./html/footer.html"))
	gDb, _     = sql.Open("sqlite3", "cache/web.db")
	gCreateDb  = "create table if not exists pages(title text, body blob, timestamp text)"
)

func (this *Page) saveCache() error {
	//-------------------------------------------
	// 1. file
	f := "cache/" + this.Title + ".txt"
	ioutil.WriteFile(f, this.Body, 0600)

	//-------------------------------------------
	// 2. db
	timestamp := strconv.FormatInt(time.Now().Unix(), 10) // time to string
	gDb.Exec(gCreateDb)

	tx, _ := gDb.Begin()
	stmt, _ := tx.Prepare("insert into pages(title, body, timestamp) values(?, ?, ?)")
	_, err := stmt.Exec(this.Title, this.Body, timestamp)
	tx.Commit()

	return err
}

func load(title string) (*Page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func loadFromDb(whereTitle string) (*Page, error) {
	var title string
	var body []byte

	q, err := gDb.Query("select title, body from pages where title='" + whereTitle + "' order by timestamp Desc limit 1")
	if err != nil {
		return nil, err
	}

	for q.Next() {
		q.Scan(&title, &body)
	}

	return &Page{Title: title, Body: body}, nil
}

func view(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]

	p, err := loadFromDb(title)
	if err != nil {
		p, _ = load(title)
	}

	if p.Title == "" {
		p, _ = load(title)
	}

	gTplView.ExecuteTemplate(w, "base", p)
}

func edit(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]

	p, err := loadFromDb(title)
	if err != nil {
		p, _ = load(title)
	}

	if p.Title == "" {
		p, _ = load(title)
	}

	gTplEdit.ExecuteTemplate(w, "base", p)
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.saveCache()

	// redirect : 다른 페이지로 이동
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		title := "Upload"
		p := &Page{Title: title}
		gTplUpload.ExecuteTemplate(w, "base", p)

	case "POST":
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		m := r.MultipartForm
		files := m.File["myfiles"]
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			dest, err := os.Create("./files/" + files[i].Filename)
			defer dest.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dest, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/files/"+files[i].Filename, http.StatusFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]

	if len(title) <= 0 {
		fmt.Fprintf(w, "This is index page.")
		return
	}

	p := &Page{Title: title, Body: []byte("고 언어 재미 있어요. from index page.")}
	p.saveCache()

	fmt.Fprintf(w, "[%s.txt] generated...", title)
}

func main() {
	// resource 파일 서버 위치
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("res"))))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	// 일반 web handler 등록
	http.HandleFunc("/", index)
	http.HandleFunc("/view/", view)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/save/", save)
	http.HandleFunc("/upload/", upload)

	// http 리스너 시작
	http.ListenAndServe(":8000", nil)
}
