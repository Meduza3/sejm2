package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
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

type Room struct {
	ID              string
	Players         map[int]Player
	Clients         map[*websocket.Conn]bool
	Axes            [4]int
	Mu              sync.Mutex
	Niezrzeszeni    int
	NumerGlosowania int
}

var rooms = make(map[string]*Room)

func CreateRoom(id string) *Room {
	room := &Room{
		ID:      id,
		Players: make(map[int]Player),
		Clients: make(map[*websocket.Conn]bool),
	}
	rooms[id] = room
	room.NumerGlosowania = 1
	return room
}

func FindRoom(id string) (*Room, bool) {
	room, exists := rooms[id]
	//fmt.Printf("Finding a room with id %s\n", id)
	//fmt.Println(exists)
	return room, exists
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
	Action          string         `json:"action"`
	PlayerID        int            `json:"playerID,omitempty"`
	Ustawa          string         `json:"ustawa,omitempty"`
	Opinions        [4][4]int      `json:"opinions,omitempty"`
	SumaZa          int            `json:"sumaZa,omitempty"`
	SumaPrzeciw     int            `json:"sumaPrzeciw,omitempty"`
	SumaWstrzymal   int            `json:"sumaWstrzymal,omitempty"`
	NumerGlosowania int            `json:"numer",omitempty`
	Changes         map[string]int `json:"changes,omitempty"`
}

type PlayersMessage struct {
	Players []Player `json:"players"`
}

type IdMessage struct {
	Id int `json:"Id"`
}

var (
	numberOfPlayers int = 8
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

func handleConnections(w http.ResponseWriter, r *http.Request) {

	roomID := r.URL.Query().Get("roomID")
	//fmt.Printf("handleConnections called with roomID: %s\n", roomID)
	room, exists := FindRoom(roomID)
	if !exists {
		room = CreateRoom(roomID)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("WebSocket connection established")
	defer func() {
		fmt.Println("WebSocket connection closed")
		ws.Close()
	}()

	for {
		//fmt.Println("dupadupadupa")
		var msg WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			//room.Mu.Lock()
			delete(room.Clients, ws)
			//room.Mu.Unlock()
			break
		}

		//fmt.Printf("%s: Received message: %+v\n", room.ID, msg)

		switch msg.Action {
		case "join":
			fmt.Println("Test")
			handleJoin(room, ws)
		case "leave":
			handleLeave(room, msg.PlayerID)
		case "za":
			roomHandleZa(room, msg.PlayerID)
		case "przeciw":
			roomHandlePrzeciw(room, msg.PlayerID)
		case "wstrzymaj":
			roomHandleWstrzymaj(room, msg.PlayerID)
		case "ustawa":
			fmt.Println("Handling ustawa!")
			handleUstawa(room, msg.Ustawa)
		case "opinions":
			fmt.Println("Handling opinions!")
			handleOpinion(room, msg.PlayerID, msg.Opinions)
		}
	}
}

func handleOpinion(room *Room, playerID int, opinions [4][4]int) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		fmt.Println("Player does not exist in this room")
		return // Player must exist in room to modify opinions
	}

	player.Opinions = opinions      // Modify the field you want
	room.Players[playerID] = player // Put the modified struct back into the map

	// Convert your players map to a slice
	playersSlice := make([]Player, 0, len(room.Players))
	for _, player := range room.Players {
		playersSlice = append(playersSlice, player)
	}

	// Create an instance of PlayersMessage and set the Players field
	message := PlayersMessage{Players: playersSlice}

	// Broadcast the message
	broadcastToRoom(room, message)
}

func handleUstawa(room *Room, ustawa string) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	firstLetter := string(ustawa[0])
	secondNumber, err := strconv.Atoi(string(ustawa[1:3]))
	if err != nil {
		fmt.Printf("Dupa")
	}
	fmt.Printf("First letter of code: %s Second letter of code: %d\n", firstLetter, secondNumber)

	switch firstLetter {
	case "A":
		if room.Axes[0] == secondNumber {
			room.Axes[0] = 0
		} else {
			room.Axes[0] = secondNumber
		}

	case "B":
		if room.Axes[1] == secondNumber {
			room.Axes[1] = 0
		} else {
			room.Axes[1] = secondNumber
		}
	case "C":
		if room.Axes[2] == secondNumber {
			room.Axes[2] = 0
		} else {
			room.Axes[2] = secondNumber
		}
	case "D":
		if room.Axes[3] == secondNumber {
			room.Axes[3] = 0
		} else {
			room.Axes[3] = secondNumber
		}

	}
	axesData := map[string][4]int{"axes": room.Axes}
	broadcastToRoom(room, axesData)
}

func handleLeave(room *Room, playerID int) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	fmt.Printf("%s: Player %d LEFT\n", room.ID, playerID)
	delete(room.Players, playerID)
	playersSlice := make([]Player, 0, len(room.Players))
	for _, player := range room.Players {
		playersSlice = append(playersSlice, player)
	}
	message := PlayersMessage{Players: playersSlice}

	broadcastToRoom(room, message)
}

func roomHandleZa(room *Room, playerID int) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		fmt.Println("Player does not exist in this room")
		return
	}

	if player.Vote != FOR {
		fmt.Printf("%s: Player %d voted ZA\n", room.ID, playerID)
		player.Vote = FOR
		room.Players[playerID] = player
	} else {
		fmt.Printf("%s: Player %d cancelled their vote.\n", room.ID, playerID)
		player.Vote = NULL
		room.Players[playerID] = player
	}

	checkForEndOfRound(room)
}

func roomHandlePrzeciw(room *Room, playerID int) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		fmt.Println("Player does not exist in this room")
		return
	}

	if player.Vote != AGAINST {
		fmt.Printf("%s: Player %d voted PRZECIW\n", room.ID, playerID)
		player.Vote = AGAINST
		room.Players[playerID] = player
	} else {
		fmt.Printf("%s: Player %d cancelled their vote.\n", room.ID, playerID)
		player.Vote = NULL
		room.Players[playerID] = player
	}

	checkForEndOfRound(room)
}

func roomHandleWstrzymaj(room *Room, playerID int) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		fmt.Println("Player does not exist in this room")
		return
	}

	if player.Vote != ABSTAIN {
		fmt.Printf("%s: Player %d voted WSTRZYMAJ\n", room.ID, playerID)
		player.Vote = ABSTAIN
		room.Players[playerID] = player
	} else {
		fmt.Printf("%s: Player %d cancelled their vote.\n", room.ID, playerID)
		player.Vote = NULL
		room.Players[playerID] = player
	}

	checkForEndOfRound(room)
}

func checkForEndOfRound(room *Room) {

	var roundShouldEnd = true

	playersSlice := make([]Player, 0, len(room.Players))
	for _, player := range room.Players {
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
		calculateRound(room)
	}
}

func handleJoin(room *Room, ws *websocket.Conn) {
	//broadcastToRoom(room, WSMessage{Action: "someonejoins"})
	//fmt.Printf("An attempt to join room %s\n", room.ID)
	//room.Mu.Lock()
	//defer room.Mu.Unlock()

	room.Clients[ws] = true

	a := randInt()
	b := randInt()
	c := randInt()
	d := randInt()
	count := math.Floor(float64(460 / numberOfPlayers))
	playerID := len(room.Players) + 1
	player := Player{
		Id:    playerID,
		Count: int(count),
		Opinions: [4][4]int{
			{a, clampToFour(a + randMod()), clampToFour(a + randMod()), clampToFour(a + randMod())},
			{b, clampToFour(b + randMod()), clampToFour(b + randMod()), clampToFour(b + randMod())},
			{c, clampToFour(c + randMod()), clampToFour(c + randMod()), clampToFour(c + randMod())},
			{d, clampToFour(d + randMod()), clampToFour(d + randMod()), clampToFour(d + randMod())},
		},
		Vote: NULL,
	}
	fmt.Printf("%s: Player %d joined\n", room.ID, player.Id)
	room.Players[player.Id] = player

	// Convert your players map to a slice
	playersSlice := make([]Player, 0, len(room.Players))
	for _, player := range room.Players {
		playersSlice = append(playersSlice, player)
	}

	// Create an instance of PlayersMessage and set the Players field
	message := PlayersMessage{Players: playersSlice}
	ws.WriteJSON(IdMessage{Id: playerID})
	broadcastToRoom(room, message)
}

func broadcastToRoom(room *Room, message interface{}) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()
	for client := range room.Clients {
		err := client.WriteJSON(message)
		if err != nil {
			fmt.Printf("error broadcasting to room %s: %v", room.ID, err)
			client.Close()
			delete(room.Clients, client)
		}
	}
}

func calculateRound(room *Room) {
	//room.Mu.Lock()
	//defer room.Mu.Unlock()
	var playersTemp []Player

	for _, player := range room.Players { // Iterate through the map
		playersTemp = append(playersTemp, player)
		// Appending each player to playersTemp slice
	}

	sort.Slice(playersTemp, func(i, j int) bool {
		return playersTemp[i].Id < playersTemp[j].Id
	})

	sumaZa := 0
	sumaPrzeciw := 0
	sumaWstrzymal := 0
	for _, player := range room.Players {
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
		sumaPrzeciw += room.Niezrzeszeni
		niezrzeszeniVote = "PRZECIW"
	} else {
		sumaZa += room.Niezrzeszeni
		niezrzeszeniVote = "ZA"
	}

	fmt.Printf("%s: Here are the players:\n", room.ID)
	for _, tempPlayer := range playersTemp {
		fmt.Println(tempPlayer)
	}
	fmt.Printf("%s: Niezrzeszeni: %d zaglosowalo %s\n ", room.ID, room.Niezrzeszeni, niezrzeszeniVote)

	fmt.Printf("%s: Here is the legislation: ", room.ID)
	fmt.Println(room.Axes)

	for _, gracz := range playersTemp { //Dla kazdego gracza
		for i := 0; i <= 3; i++ { //Dla kazdej z osi
			if room.Axes[i] != 0 {
				var bloczki float64 = 0
				for j := 0; j <= 3; j++ { //Dla kazdego z klockow tego gracza
					if gracz.Vote == AGAINST {
						if isInLegislationArea(gracz.Opinions[i][j], room.Axes[i]) {
							bloczki += 25
						}

					} else if gracz.Vote == FOR {
						if !isInLegislationArea(gracz.Opinions[i][j], room.Axes[i]) {
							bloczki += 25
						}
					}
				}
				if bloczki != 0 {
					var odchodzacy = math.Ceil((bloczki / 100) * 0.2 * float64(gracz.Count)) // Mamy ilosc odchodzacych
					fmt.Printf("%s: Gracz %d, os %s: Wkurzylo sie %d bloczkow, co daje %d odchodzacych\n", room.ID, gracz.Id, numToAxis(i), int(bloczki/25), int(odchodzacy))
					naTejOsiWystajeNaPrawo := wystawalbyNaPrawo(gracz.Opinions[i][:], room.Axes[i])
					if naTejOsiWystajeNaPrawo {
						fmt.Printf("%s: Gracz %d, os %s: wystaje na prawo\n", room.ID, gracz.Id, numToAxis(i))
					} else {
						fmt.Printf("%s: Gracz %d, os %s: wystaje na lewo\n", room.ID, gracz.Id, numToAxis(i))
					}
					minDistance := 100
					var najblizszaPartia *Player

					najbardziejNaPrawo := max(gracz.Opinions[i][:])
					fmt.Printf("%s: Gracz %d, os %s: Najbardziej na prawo wysuniety klocek jest na %d\n", room.ID, gracz.Id, numToAxis(i), najbardziejNaPrawo)

					najbardziejNaLewo := min(gracz.Opinions[i][:])
					fmt.Printf("%s: Gracz %d, os %s: Najbardziej na lewo wysuniety klocek jest na %d\n", room.ID, gracz.Id, numToAxis(i), najbardziejNaLewo)
					for _, drugiGracz := range playersTemp {
						if drugiGracz.Vote != gracz.Vote && drugiGracz.Vote != ABSTAIN {
							if gracz.Vote == FOR && naTejOsiWystajeNaPrawo {
								//closest party to the right
								najblizszaOpiniaDrugiegoGracza := minGreaterThanThreshold(drugiGracz.Opinions[i][:], najbardziejNaPrawo)
								fmt.Printf("%s: Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", room.ID, gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
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
								fmt.Printf("%s: Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", room.ID, gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
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
								fmt.Printf("%s: Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", room.ID, gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
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
								fmt.Printf("%s: Gracz %d, os %s, drugi gracz %d: Najblizsza opinia drugiego gracza jest na %d\n", room.ID, gracz.Id, numToAxis(i), drugiGracz.Id, najblizszaOpiniaDrugiegoGracza)
								if najblizszaOpiniaDrugiegoGracza != 420 {
									distance := abs(najbardziejNaPrawo - najblizszaOpiniaDrugiegoGracza)
									if drugiGracz.Vote != AGAINST && distance < minDistance {
										minDistance = distance
										najblizszaPartia = &drugiGracz
									}
								}
							}
						} else {
							fmt.Printf("%s: Gracz %d, os %s: drugi gracz %d zaglosowal tak samo, pomijamy\n", room.ID, gracz.Id, numToAxis(i), drugiGracz.Id)
						}
					}
					if najblizszaPartia != nil {
						fmt.Printf("%s: Gracz %d, os %s: Poslowie przejda do partii gracza %d\n", room.ID, gracz.Id, numToAxis(i), najblizszaPartia.Id)
					} else {
						fmt.Printf("%s: Gracz %d, os %s: Poslowie przejda do niezrzeszonych\n", room.ID, gracz.Id, numToAxis(i))
					}
					player, exists := room.Players[gracz.Id]
					if exists {
						player.Count -= int(odchodzacy)
						room.Players[gracz.Id] = player
					} else {
						fmt.Println("Player does not exist.")
					}
					if najblizszaPartia != nil {
						secondPlayer := room.Players[najblizszaPartia.Id]
						secondPlayer.Count += int(odchodzacy)
						room.Players[najblizszaPartia.Id] = secondPlayer
					} else {
						room.Niezrzeszeni += int(odchodzacy)
					}
				} else {
					fmt.Printf("%s: Gracz %d, os %s: Nikogo nie wkurzyl, pomijamy\n", room.ID, gracz.Id, numToAxis(i))
				}

			} else {
				fmt.Printf("%s: Gracz %d, os %s: nie dotyczy tej ustawy, pomijamy\n", room.ID, gracz.Id, numToAxis(i))
			}
		}
		resetVotes(room)
	}

	fmt.Println("Here are the results:")

	fmt.Printf("%s: GLOSOWALO: %d\nZA: %d\nPRZECIW: %d\nWSTRZYMALO SIE: %d\n", room.ID, sumaZa+sumaPrzeciw+sumaWstrzymal, sumaZa, sumaPrzeciw, sumaWstrzymal)
	if sumaZa > sumaPrzeciw {
		fmt.Printf("%s: Ustawa PRZESZLA\n", room.ID)
	} else {
		fmt.Printf("%s: Ustawa ODRZUCONA\n", room.ID)
	}
	changes := make(map[string]int)
	for _, originalPlayer := range playersTemp {
		player := room.Players[originalPlayer.Id]
		changeCount := player.Count - originalPlayer.Count
		fmt.Printf("%s: Gracz %d ma teraz %d poslow. Zmiana: %d\n", room.ID, player.Id, player.Count, changeCount)
		changes[strconv.Itoa(player.Id)] = changeCount
	}
	fmt.Printf("Niezrzeszeni: %d\n", room.Niezrzeszeni)
	var message = WSMessage{Action: "results", SumaZa: sumaZa, SumaPrzeciw: sumaPrzeciw, SumaWstrzymal: sumaWstrzymal, NumerGlosowania: room.NumerGlosowania, Changes: changes}
	broadcastToRoom(room, message)
	room.NumerGlosowania++
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

func resetVotes(room *Room) {
	for id := range room.Players {
		player := room.Players[id]
		player.Vote = NULL
		room.Players[id] = player
	}

	message := WSMessage{Action: "resetVotes"}
	broadcastToRoom(room, message)
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

	playersStr := flag.String("players", "8", "The number of players")
	port := os.Getenv("PORT")

	// Parse the command-line flags
	flag.Parse()

	// Convert the number of players from string to int
	numberOfPlayers, _ = strconv.Atoi(*playersStr)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/ws", handleConnections)

	fmt.Printf("Server started on port %s\n!!!", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
