package main

import (
	"net/http"
)

func setSession(u *User, w http.ResponseWriter) {
	// JSON 형태
	val := map[string]string{
		"id":    u.Id,
		"pw":    u.Pw,
		"email": u.Email,
		"fname": u.Fname,
		"lname": u.Lname,
	}

	if encoded, err := cookieHandler.Encode("session", val); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserInfo(r *http.Request, key string) (result string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			result = cookieValue[key]
		}
	}
	return result
}

func clearSession(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

//------------------------------------------------------
func setFlashMsg(w http.ResponseWriter, name string, msg string) {
	// JSON 형태
	val := map[string]string{
		name: msg,
	}

	if encoded, err := cookieHandler.Encode(name, val); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getFlashMsg(w http.ResponseWriter, r *http.Request, name string) (msg string) {
	if cookie, err := r.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			clearSession(w, name)
		}
	}
	return msg
}
