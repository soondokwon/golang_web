package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

func encode(s []byte) string {
	return base64.URLEncoding.EncodeToString(s)
}

func decode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

func setFlashMsg(w http.ResponseWriter, name string, msg []byte) {
	c := &http.Cookie{Name: name, Value: encode(msg)}
	http.SetCookie(w, c) // browser에 쿠키를 설정
}

func getFlashMsg(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil { // 오류가 발생하면...
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	val, err := decode(c.Value)
	if err != nil { // 오류가 발생하면...
		return nil, err
	}

	// 현재 사용한 cookie를 삭제
	deleteCookie := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(0, 1)}
	http.SetCookie(w, deleteCookie) // browser에 쿠키를 설정

	return val, nil

}
