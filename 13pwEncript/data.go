package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid   string
	Id     string
	Pw     string
	Fname  string
	Lname  string
	Email  string
	Errors map[string]string
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	db.Exec("create table if not exists users(uuid text, userid text, pw text, firstname text, lastname text, email text)")

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users(uuid, userid, pw, firstname, lastname, email) values(?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Uuid, u.Id, encPass(u.Pw), u.Fname, u.Lname, u.Email)
	tx.Commit()

	return err
}

func userExists(id string, pw string) (*User, bool) {
	u := &User{}
	result := false

	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	rows, err := db.Query("select userid, pw, firstname, lastname, email from users where userid = $1", id)
	if err != nil {
		return nil, false
	}

	for rows.Next() {
		rows.Scan(&u.Id, &u.Pw, &u.Fname, &u.Lname, &u.Email)
	}

	pwResult := bcrypt.CompareHashAndPassword([]byte(u.Pw), []byte(pw))
	if id == u.Id && pwResult == nil {
		result = true
	}

	return u, result
}

func makeUserList() (*[]User, error) {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	rows, err := db.Query("select userid, pw, email from users")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Pw, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return &users, nil
}

func encPass(pw string) string {
	//log.Println("pw : ", pw)

	pass := []byte(pw)
	hashPw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	//log.Println("hashPw : ", hashPw)

	return string(hashPw)
}

func genUUID() (id string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return
	}
	id = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return id
}
