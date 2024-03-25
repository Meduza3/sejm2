package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
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
	niezrzeszeni int = 4
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
	player, exists := players[playerID]
	if !exists {
		fmt.Println("Player does not exist")
		return
	}
	if player.Vote != FOR {
		fmt.Printf("Player %d voted ZA\n", playerID)
		player.Vote = FOR
		players[playerID] = player
	} else {
		fmt.Printf("Player %d cancelled their vote.\n", playerID)
		player.Vote = NULL
		players[playerID] = player
	}

	checkForEndOfRound()
}

func handlePrzeciwVote(playerID int) {
	player, exists := players[playerID]
	if !exists {
		fmt.Println("Player does not exist")
		return
	}
	if player.Vote != AGAINST {
		fmt.Printf("Player %d voted PRZECIW\n", playerID)
		player.Vote = AGAINST
		players[playerID] = player
	} else {
		fmt.Printf("Player %d cancelled their vote.\n", playerID)
		player.Vote = NULL
		players[playerID] = player
	}

	checkForEndOfRound()
}

func handleWstrzymajVote(playerID int) {
	player, exists := players[playerID]
	if !exists {
		fmt.Println("Player does not exist")
		return
	}
	if player.Vote != ABSTAIN {
		fmt.Printf("Player %d voted WSTRZYMAJ\n", playerID)
		player.Vote = ABSTAIN
		players[playerID] = player
	} else {
		fmt.Printf("Player %d cancelled their vote.\n", playerID)
		player.Vote = NULL
		players[playerID] = player
	}

	checkForEndOfRound()
}

func checkForEndOfRound() {

	var roundShouldEnd = true

	playersSlice := make([]Player, 0, len(players))
	for _, player := range players {
		playersSlice = append(playersSlice, player)
	}

	for _, sliced := range playersSlice {
		if sliced.Vote == NULL {
			roundShouldEnd = false
			break
		}
	}

	if roundShouldEnd {
		fmt.Println("EVERYONE VOTED!")
		calculateRound()
	}
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
		Count: 57,
		Opinions: [4][4]int{
			{a, clampToFour(a + randMod()), clampToFour(a + randMod()), clampToFour(a + randMod())},
			{b, clampToFour(b + randMod()), clampToFour(b + randMod()), clampToFour(b + randMod())},
			{c, clampToFour(c + randMod()), clampToFour(c + randMod()), clampToFour(c + randMod())},
			{d, clampToFour(d + randMod()), clampToFour(d + randMod()), clampToFour(d + randMod())},
		},
		Vote: NULL,
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

func calculateRound() {
	var playersTemp []Player

	for _, player := range players { // Iterate through the map
		playersTemp = append(playersTemp, Player{Id: player.Id, Count: player.Count, Opinions: player.Opinions, Vote: player.Vote})
		// Appending each player to playersTemp slice
	}

	sort.Slice(playersTemp, func(i, j int) bool {
		return playersTemp[i].Id < playersTemp[j].Id
	})

	sumaZa := 0
	sumaPrzeciw := 0
	sumaWstrzymal := 0
	for _, player := range players {
		if player.Vote == FOR {
			sumaZa += player.Count
		} else if player.Vote == AGAINST {
			sumaPrzeciw += player.Count
		} else if player.Vote == ABSTAIN {
			sumaWstrzymal += player.Count
		}
	}
	niezrzeszeniZa := rand.Intn(2)
	var niezrzeszeniVote string
	if niezrzeszeniZa == 0 { // Generates either 0 or 1 randomly
		sumaPrzeciw += niezrzeszeni
		niezrzeszeniVote = "PRZECIW"
	} else {
		sumaZa += niezrzeszeni
		niezrzeszeniVote = "ZA"
	}

	fmt.Println("Here are the players:")
	for _, tempPlayer := range playersTemp {
		fmt.Println(tempPlayer)
	}
	fmt.Printf("Niezrzeszeni: %d zaglosowalo %s\n ", niezrzeszeni, niezrzeszeniVote)

	fmt.Println("Here is the legislation:")
	fmt.Println(axes)

	for _, gracz := range playersTemp { //Dla kazdego gracza
		for i := 0; i <= 3; i++ { //Dla kazdej z osi
			if axes[i] != 0 {
				var bloczki float64 = 0
				for j := 0; j <= 3; j++ { //Dla kazdego z klockow tego gracza
					if gracz.Vote == AGAINST {
						if isInLegislationArea(gracz.Opinions[i][j], axes[i]) {
							bloczki += 25
						}

					} else if gracz.Vote == FOR {
						if !isInLegislationArea(gracz.Opinions[i][j], axes[i]) {
							bloczki += 25
						}
					}
				}
				if bloczki != 0 {
					var odchodzacy = math.Ceil((bloczki / 100) * 0.2 * float64(gracz.Count)) // Mamy ilosc odchodzacych
					fmt.Printf("Gracz %d, os %s: Wkurzylo sie %d bloczkow, co daje %d odchodzacych\n", gracz.Id, numToAxis(i), int(bloczki/25), int(odchodzacy))
					naTejOsiWystajeNaPrawo := wystawalbyNaPrawo(gracz.Opinions[i][:], axes[i])
					if naTejOsiWystajeNaPrawo {
						fmt.Printf("Gracz %d, os %s: wystaje na prawo\n", gracz.Id, numToAxis(i))
					} else {
						fmt.Printf("Gracz %d, os %s: wystaje na lewo\n", gracz.Id, numToAxis(i))
					}
					minDistance := 100
					var najblizszaPartia *Player

					najbardziejNaPrawo := max(gracz.Opinions[i][:])
					fmt.Printf("Gracz %d, os %s: Najbardziej na prawo wysuniety klocek jest na %d\n", gracz.Id, numToAxis(i), najbardziejNaPrawo)

					najbardziejNaLewo := min(gracz.Opinions[i][:])
					fmt.Printf("Gracz %d, os %s: Najbardziej na lewo wysuniety klocek jest na %d\n", gracz.Id, numToAxis(i), najbardziejNaLewo)
					for _, drugiGracz := range playersTemp {
						if drugiGracz.Vote != gracz.Vote {
							if gracz.Vote == FOR && naTejOsiWystajeNaPrawo {
								//closest party to the right
								najblizszaOpiniaDrugiegoGracza := minGreaterThanThreshold(drugiGracz.Opinions[i][:], najbardziejNaPrawo)
								fmt.Printf("Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
								if najblizszaOpiniaDrugiegoGracza != 420 {
									distance := abs(najbardziejNaPrawo - najblizszaOpiniaDrugiegoGracza)
									if drugiGracz.Vote != FOR && distance < minDistance {
										minDistance = distance
										najblizszaPartia = &drugiGracz
									}
								}

							} else if gracz.Vote == FOR && !naTejOsiWystajeNaPrawo {
								//Closest party to the left
								najblizszaOpiniaDrugiegoGracza := maxSmallerThanThreshold(drugiGracz.Opinions[i][:], najbardziejNaLewo)
								fmt.Printf("Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
								if najblizszaOpiniaDrugiegoGracza != -420 {
									distance := abs(najbardziejNaLewo - najblizszaOpiniaDrugiegoGracza)
									if drugiGracz.Vote != FOR && distance < minDistance {
										minDistance = distance
										najblizszaPartia = &drugiGracz
									}
								}

							} else if gracz.Vote == AGAINST && naTejOsiWystajeNaPrawo {
								//Closest party to the left
								najblizszaOpiniaDrugiegoGracza := maxSmallerThanThreshold(drugiGracz.Opinions[i][:], najbardziejNaLewo)
								fmt.Printf("Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
								if najblizszaOpiniaDrugiegoGracza != -420 {
									distance := abs(najbardziejNaLewo - najblizszaOpiniaDrugiegoGracza)
									if drugiGracz.Vote != AGAINST && distance < minDistance {
										minDistance = distance
										najblizszaPartia = &drugiGracz
									}
								}

							} else if gracz.Vote == AGAINST && !naTejOsiWystajeNaPrawo {
								//Closest party to the right
								najblizszaOpiniaDrugiegoGracza := minGreaterThanThreshold(drugiGracz.Opinions[i][:], najbardziejNaPrawo)
								fmt.Printf("Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
								if najblizszaOpiniaDrugiegoGracza != 420 {
									distance := abs(najbardziejNaPrawo - najblizszaOpiniaDrugiegoGracza)
									if drugiGracz.Vote != AGAINST && distance < minDistance {
										minDistance = distance
										najblizszaPartia = &drugiGracz
									}
								}
							}
						} else {
							fmt.Printf("Gracz %d, os %s: drugi gracz %d zaglosowal tak samo, pomijamy\n", gracz.Id, numToAxis(i), drugiGracz.Id)
						}
					}
					if najblizszaPartia != nil {
						fmt.Printf("Gracz %d, os %s: Poslowie przejda do partii gracza %d\n", gracz.Id, numToAxis(i), najblizszaPartia.Id)
					} else {
						fmt.Printf("Gracz %d, os %s: Poslowie przejda do niezrzeszonych\n", gracz.Id, numToAxis(i))
					}
					player, exists := players[gracz.Id]
					if exists {
						player.Count -= int(odchodzacy)
						players[gracz.Id] = player
					} else {
						fmt.Println("Player does not exist.")
					}
					if najblizszaPartia != nil {
						secondPlayer := players[najblizszaPartia.Id]
						secondPlayer.Count += int(odchodzacy)
						players[najblizszaPartia.Id] = secondPlayer
					} else {
						niezrzeszeni += int(odchodzacy)
					}
				} else {
					fmt.Printf("Gracz %d, os %s: Nikogo nie wkurzyl, pomijamy\n", gracz.Id, numToAxis(i))
				}

			} else {
				fmt.Printf("Gracz %d, os %s: nie dotyczy tej ustawy, pomijamy\n", gracz.Id, numToAxis(i))
			}
		}
		resetVotes()
	}

	fmt.Println("Here are the results:")

	fmt.Printf("GLOSOWALO: %d\nZA: %d\nPRZECIW: %d\nWSTRZYMALO SIE: %d\n", sumaZa+sumaPrzeciw+sumaWstrzymal, sumaZa, sumaPrzeciw, sumaWstrzymal)
	if sumaZa > sumaPrzeciw {
		fmt.Printf("Ustawa PRZESZLA\n")
	} else {
		fmt.Printf("Ustawa ODRZUCONA\n")
	}
	for _, originalPlayer := range playersTemp {
		player := players[originalPlayer.Id]
		fmt.Printf("Gracz %d ma teraz %d poslow. Zmiana: %d\n", player.Id, player.Count, player.Count-originalPlayer.Count)
	}
	fmt.Printf("Niezrzeszeni: %d\n", niezrzeszeni)

}

func wystawalbyNaPrawo(opinion []int, legislation int) bool {
	maxDistance := 0
	for _, v := range opinion {
		distance := legislation - v
		if abs(distance) > abs(maxDistance) {
			maxDistance = distance
		}
	}
	return maxDistance < 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(nums []int) int {
	minValue := 420
	for _, num := range nums {
		if num < minValue {
			minValue = num
		}
	}
	return minValue
}

func max(nums []int) int {
	maxValue := -420
	for _, num := range nums {
		if num > maxValue {
			maxValue = num
		}
	}
	return maxValue
}

func minGreaterThanThreshold(nums []int, threshold int) int {
	minValue := 420
	for _, num := range nums {
		if num >= threshold && num < minValue {
			minValue = num
		}
	}
	return minValue
}

func maxSmallerThanThreshold(nums []int, threshold int) int {
	maxValue := -420
	for _, num := range nums {
		if num <= threshold && num > maxValue {
			maxValue = num
		}
	}
	return maxValue
}

func numToAxis(num int) string {
	switch num {
	case 0:
		return "A"
	case 1:
		return "B"
	case 2:
		return "C"
	case 3:
		return "D"
	}
	return "ugabuga"
}

func resetVotes() {
	for id := range players {
		player := players[id]
		player.Vote = NULL
		players[id] = player
	}

	message := WSMessage{Action: "resetVotes"}
	for socket := range clients {
		socket.WriteJSON(message)
	}
}

func isInLegislationArea(opinion int, legislation int) bool {
	switch legislation {
	case -4, -3, -2, 2, 3, 4:
		if opinion == legislation || opinion == legislation-1 || opinion == legislation+1 {
			return true
		}
	case -1:
		if opinion == -1 || opinion == -2 || opinion == 1 {
			return true
		}
	case 1:
		if opinion == 1 || opinion == 2 || opinion == -1 {
			return true
		}
	}
	return false
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
