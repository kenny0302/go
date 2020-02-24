package main

import (
	"log"
	"net/http"
	"token"
	"github.com/gorilla/websocket"
	"time"
)

func main() {
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			log.Println(time.Now())
			return
		}
		defer func() {
			log.Println("disconnect !!")
			
			c.Close()
		}()
		
		err = c.WriteMessage(websocket.TextMessage, token.Token())
		if err != nil {
			log.Println(err)
			return
		}
		
	})
	log.Println("server start at :8899")
	log.Fatal(http.ListenAndServe(":8899", nil))
}
