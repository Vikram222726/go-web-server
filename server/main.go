package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const MAX_QUEUE_SIZE = 10000
const MAX_WORKERS = 20

var workerPool = make(chan struct{}, MAX_WORKERS)
var numChannel = make(chan int, MAX_QUEUE_SIZE)

type RequestData struct {
	Number int `json:"number"`
}

type ResponseData struct {
	Status   int    `json:"status"`
	Response string `json:"response"`
}

func checkIsPrime(num int) bool {
	// if num <= 1 {
	// 	fmt.Println("Number 1 is not prime")
	// 	return false
	// }
	// if num <= 3 {
	// 	fmt.Println("Number is prime", num)
	// 	return true
	// }
	// if num%2 == 0 || num%3 == 0 {
	// 	fmt.Println("Number is not prime", num)
	// 	return false
	// }
	// for i := 5; i*i <= num; i += 6 {
	// 	if num%i == 0 || num%(i+2) == 0 {
	// 		fmt.Println("Number is not prime", num)
	// 		return false
	// 	}
	// }
	for i := num - 1; i > 1; i-- {
		if num%i == 0 {
			fmt.Println("Number is not prime", num)
			return false
		}
	}
	fmt.Println("Number is prime", num)
	return true
}

func processPrimeNums() {
	for {
		select {
		case num := <-numChannel:
			fmt.Println("Next Num in Channel is: ", num)
			<-workerPool
			checkIsPrime(num)
			workerPool <- struct{}{}
		}
	}
}

func init() {
	for i := 0; i < MAX_WORKERS; i++ {
		workerPool <- struct{}{}
		go processPrimeNums()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func primeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only accepting POST request on this endpoint", http.StatusMethodNotAllowed)
		return
	}

	var newData RequestData
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, "Invalid Body String!", http.StatusBadRequest)
	}

	if newData.Number > 5000000000 {
		http.Error(w, "Number out of Bound!", http.StatusBadRequest)
		return
	}

	fmt.Println("Received new number: ", newData.Number)

	// Add this new number inside the numChannel to be processed by goroutines
	numChannel <- newData.Number

	response := ResponseData{
		Status:   200,
		Response: "Request Acknowledged, sent for processing!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/prime", primeHandler)

	fmt.Println("Starting server on Port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server Failed!", err)
	}
}
