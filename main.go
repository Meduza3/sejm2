package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type Player struct {
	id       int
	count    int
	IP       string
	opinions [4][4]int
}

type RequestBody struct {
	Code string `json:"code"`
}

var (
	players          = make(map[string]Player)
	player_count int = 0
	axisA        int = 0
	axisB        int = 0
	axisC        int = 0
	axisD        int = 0
)

func main() {

	var mu sync.Mutex

	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("Server started!")

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		player_count++
		ip := r.RemoteAddr
		player := Player{id: player_count, count: 100, IP: ip, opinions: [4][4]int{{-1, -1, 1, 1}, {-1, -1, 1, 1}, {-1, -1, 1, 1}, {-1, -1, 1, 1}}}
		players[ip] = player

		w.Header().Set("Content-Type", "application/json")
		response := fmt.Sprintf(`{"id": %d, "count": %d, "IP": "%s"}`, player.id, player.count, player.IP)
		w.Write([]byte(response))
	})

	http.HandleFunc("/ustawa", func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != "POST" {
			http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
			return
		}

		var reqBody RequestBody

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(reqBody.Code) > 0 {
			firstLetter := string(reqBody.Code[0])
			secondNumber, err := strconv.Atoi(string(reqBody.Code[1:3]))
			if err != nil {
				fmt.Printf("Dupa")
			}
			fmt.Printf("First letter of code: %s Second letter of code: %d\n", firstLetter, secondNumber)

			switch firstLetter {
			case "A":
				if axisA == secondNumber {
					axisA = 0
				} else {
					axisA = secondNumber
				}

			case "B":
				if axisB == secondNumber {
					axisB = 0
				} else {
					axisB = secondNumber
				}
			case "C":
				if axisC == secondNumber {
					axisC = 0
				} else {
					axisC = secondNumber
				}
			case "D":
				if axisD == secondNumber {
					axisD = 0
				} else {
					axisD = secondNumber
				}

			}

			// Respond to the client
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Received your request with code starting with: %s", firstLetter)))
		} else {
			http.Error(w, "Code is empty", http.StatusBadRequest)
		}
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

	http.HandleFunc("/aktualna_ustawa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d, %d, %d, %d", axisA, axisB, axisC, axisD)
	})

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
