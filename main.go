package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Player struct {
	id    int
	count int
	IP    string
}

func main() {

	var player_count int = 0
	var players = make(map[string]Player)
	var mu sync.Mutex

	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("Server started!")

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		player_count++
		ip := r.RemoteAddr
		player := Player{id: player_count, count: 10, IP: ip}
		players[ip] = player
		fmt.Printf("JOIN Player %d from IP %s joined, count: %d\n", player.id, player.IP, player_count)
		w.Write([]byte(fmt.Sprintf("You joined as player %d!", player.id)))
	})

	http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		ip := r.RemoteAddr
		player := players[ip]
		player_count--
		fmt.Printf("LEAVE Player %d from IP %s left, count: %d\n", player.id, player.IP, player_count)
		delete(players, ip)

	})

	http.HandleFunc("/za", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Correctly delivered your ZA vote"))
		fmt.Println("Player voted ZA")
	})

	http.HandleFunc("/przeciw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Correctly delivered your PRZECIW vote"))
		fmt.Println("Player voted PRZECIW")
	})

	http.HandleFunc("/wstrzymaj", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Correctly delivered your WSTRZYMAJ vote"))
		fmt.Println("Player voted WSTRZYMAJ")
	})

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
