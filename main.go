package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id       int
	Count    int
	Opinions [4][4]int
	Vote     Vote
}

type Vote string

const (
	FOR     Vote = "FOR"
	AGAINST Vote = "AGAINST"
	ABSTAIN Vote = "ABSTAIN"
	NULL    Vote = "NULL"
)

type RequestBody struct {
	Code string `json:"code"`
}

type WSMessage struct {
	Action   string    `json:"action"`
	PlayerID int       `json:"playerID,omitempty"`
	Ustawa   string    `json:"ustawa,omitempty"`
	Opinions [4][4]int `json:"opinions,omitempty"`
}

type PlayersMessage struct {
	Players []Player `json:"players"`
}

type IdMessage struct {
	Id int `json:"Id"`
}

var (
	players          = make(map[int]Player)
	player_count int = 0
	axes         [4]int
)

func randInt() int {
	num := rand.Intn(8) - 4 // Generates a number between -4 and 3
	if num >= 0 {
		num++ // Adjusts the range to -4 to 4, excluding 0
	}
	return num
}

func randMod() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(3) - 1
}

func clampToFour(val int) int {
	if val < -4 {
		return -4
	} else if val > 4 {
		return 4
	} else if val == 0 {
		if rand.Intn(2) == 0 {
			return -1
		} else {
			return 1
		}
	}
	return val
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		var msg WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error1: %v\n", err)
			delete(clients, ws)
			player_count--
			break
		}

		switch msg.Action {
		case "join":
			handleJoin(ws)
		case "leave":
			handleLeave(msg.PlayerID)
		case "za":
			handleZaVote(msg.PlayerID)
		case "przeciw":
			handlePrzeciwVote(msg.PlayerID)
		case "wstrzymaj":
			handleWstrzymajVote(msg.PlayerID)
		case "ustawa":
			fmt.Println("Handling ustawa!")
			handleUstawa(msg.Ustawa)
		case "opinions":
			fmt.Println("Handling opinions!")
			handleOpinion(msg.PlayerID, msg.Opinions)
		}
	}
}

func handleOpinion(playerID int, opinions [4][4]int) {
	player := players[playerID] // Extract the player struct from the map
	player.Opinions = opinions  // Modify the field you want
	players[playerID] = player  // Put the modified struct back into the map

	// Convert your players map to a slice
	playersSlice := make([]Player, 0, len(players))
	for _, player := range players {
		playersSlice = append(playersSlice, player)
	}

	// Create an instance of PlayersMessage and set the Players field
	message := PlayersMessage{Players: playersSlice}

	// Broadcast the message
	broadcastToClients(message)
}

func handleUstawa(ustawa string) {
	firstLetter := string(ustawa[0])
	secondNumber, err := strconv.Atoi(string(ustawa[1:3]))
	if err != nil {
		fmt.Printf("Dupa")
	}
	fmt.Printf("First letter of code: %s Second letter of code: %d\n", firstLetter, secondNumber)

	switch firstLetter {
	case "A":
		if axes[0] == secondNumber {
			axes[0] = 0
		} else {
			axes[0] = secondNumber
		}

	case "B":
		if axes[1] == secondNumber {
			axes[1] = 0
		} else {
			axes[1] = secondNumber
		}
	case "C":
		if axes[2] == secondNumber {
			axes[2] = 0
		} else {
			axes[2] = secondNumber
		}
	case "D":
		if axes[3] == secondNumber {
			axes[3] = 0
		} else {
			axes[3] = secondNumber
		}

	}
	axesData := map[string][4]int{"axes": axes}
	broadcastToClients(axesData)
}

func handleLeave(playerID int) {
	fmt.Printf("Player %d LEFT\n", playerID)
	delete(players, playerID)
	playersSlice := make([]Player, 0, len(players))
	for _, player := range players {
		playersSlice = append(playersSlice, player)
	}
	message := PlayersMessage{Players: playersSlice}

	// Broadcast the message
	for socket := range clients {
		socket.WriteJSON(message)
	}
}

func handleZaVote(playerID int) {
	fmt.Printf("Player %d voted ZA\n", playerID)
}

func handlePrzeciwVote(playerID int) {
	fmt.Printf("Player %d voted PRZECIW\n", playerID)
}

func handleWstrzymajVote(playerID int) {
	fmt.Printf("Player %d voted WSTRZYMAJ\n", playerID)
}

func handleJoin(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	player_count++
	a := randInt()
	b := randInt()
	c := randInt()
	d := randInt()
	player := Player{
		Id:    len(players) + 1,
		Count: 100,
		Opinions: [4][4]int{
			{a, clampToFour(a + randMod()), clampToFour(a + randMod()), clampToFour(a + randMod())},
			{b, clampToFour(b + randMod()), clampToFour(b + randMod()), clampToFour(b + randMod())},
			{c, clampToFour(c + randMod()), clampToFour(c + randMod()), clampToFour(c + randMod())},
			{d, clampToFour(d + randMod()), clampToFour(d + randMod()), clampToFour(d + randMod())},
		},
	}
	fmt.Printf("Player %d joined\n", player.Id)
	players[player.Id] = player

	// Convert your players map to a slice
	playersSlice := make([]Player, 0, len(players))
	for _, player := range players {
		playersSlice = append(playersSlice, player)
	}

	// Create an instance of PlayersMessage and set the Players field
	message := PlayersMessage{Players: playersSlice}

	// Broadcast the message
	ws.WriteJSON(IdMessage{Id: player.Id})
	for socket := range clients {
		socket.WriteJSON(message)
	}
}

func broadcastToClients(message interface{}) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			fmt.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Server started!")

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
