package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Player struct {
	Id       int
	Count    int
	IP       string
	Opinions [4][4]int
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

func main() {

	var mu sync.Mutex

	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("Server started!")

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		player_count++
		ip := r.RemoteAddr
		a := randInt()
		b := randInt()
		c := randInt()
		d := randInt()
		player := Player{Id: len(players) + 1, Count: 100, IP: ip, Opinions: [4][4]int{
			{clampToFour(a), clampToFour(a + int(randMod())), clampToFour(a + int(randMod())), clampToFour(a + int(randMod()))},
			{clampToFour(b), clampToFour(b + int(randMod())), clampToFour(b + int(randMod())), clampToFour(b + int(randMod()))},
			{clampToFour(c), clampToFour(c + int(randMod())), clampToFour(c + int(randMod())), clampToFour(c + int(randMod()))},
			{clampToFour(d), clampToFour(d + int(randMod())), clampToFour(d + int(randMod())), clampToFour(d + int(randMod()))},
		}}
		players[ip] = player

		w.Header().Set("Content-Type", "application/json")

		// Marshal player struct to JSON
		jsonResponse, err := json.Marshal(player)
		if err != nil {
			// Handle error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// The player struct, including opinions, is now converted into a JSON string and sent in the response
		w.Write([]byte(jsonResponse))
	})

	http.HandleFunc("/ustawa", func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests

		switch r.Method {
		case "GET":
			resp := struct {
				AxisA int `json:"axisA"`
				AxisB int `json:"axisB"`
				AxisC int `json:"axisC"`
				AxisD int `json:"axisD"`
			}{axisA, axisB, axisC, axisD}

			respJSON, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(respJSON)

		case "POST":
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
		}

	})

	http.HandleFunc("/gracze", func(w http.ResponseWriter, r *http.Request) {
		var allOpinions []map[string][4][4]int
		mu.Lock()
		defer mu.Unlock()

		for _, player := range players {
			allOpinions = append(allOpinions, map[string][4][4]int{fmt.Sprintf("Player%d", player.Id): player.Opinions})
		}

		w.Header().Set("Content-Type", "application/json")

		// Convert the slice of opinions to JSON
		jsonResponse, err := json.Marshal(allOpinions)
		if err != nil {
			// Handle error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.Write(jsonResponse)
	})

	http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		ip := r.RemoteAddr
		player := players[ip]
		player_count--
		fmt.Printf("LEAVE Player %d from IP %s left, count: %d\n", player.Id, player.IP, player_count)
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
