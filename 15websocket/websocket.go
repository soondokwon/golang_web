package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Userid string `json:"userid"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
}

var (
	upgrader = &websocket.Upgrader{}
	count    int
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/index.html")
	})

	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("err : ", err)
					conn.Close()
					return
				}
				log.Printf("[%v] msg : [%v]", msgType, string(msg))
				conn.WriteMessage(msgType, msg)
			}
		}(conn)
	})

	// message back : 5초마다 클라이언트에 정보를 보낸다.
	http.HandleFunc("/v2/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("err : ", err)
					conn.Close()
					return
				}
				log.Printf("[%v] msg : [%v]", msgType, string(msg))
				conn.WriteMessage(msgType, msg)
			}
		}(conn)

		go func(conn *websocket.Conn) {
			ch := time.Tick(5 * time.Second)
			for range ch {
				count++

				conn.WriteJSON(User{
					Userid: "soondoe",
					Name:   "benny.kwon",
					Count:  count,
				})
			}

		}(conn)
	})

	http.ListenAndServe(":8000", nil)
}
