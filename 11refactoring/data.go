package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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

	db.Exec("create table if not exists users(userid text, pw text, firstname text, lastname text, email text)")

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users(userid, pw, firstname, lastname, email) values(?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Id, u.Pw, u.Fname, u.Lname, u.Email)
	tx.Commit()

	return err
}

func userExists(id string, pw string) (*User, bool) {
	u := &User{}
	result := false

	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	//rows, err := db.Query("select userid, pw, firstname, lastname, email from users where userid='" + id + "' and pw='" + pw + "'")
	rows, err := db.Query("select userid, pw, firstname, lastname, email from users where userid = $1 and pw= $2", id, pw)
	if err != nil {
		return nil, false
	}

	for rows.Next() {
		rows.Scan(&u.Id, &u.Pw, &u.Fname, &u.Lname, &u.Email)
	}

	if id == u.Id && pw == u.Pw {
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
