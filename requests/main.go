package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const MAX_CONNECTION_POOL_SIZE = 1000

// Struct to hold the JSON request body
type RequestBody struct {
	Number int `json:"Number"`
}

type ResponseBody struct {
	Status   int    `json:"status"`
	Response string `json:"response"`
}

func makePostRequestToServer(url string, requestId int) {
	// Generate a random number between 1000000000 and 5000000000
	randomNum := rand.Intn(4000000000) + 1000000000

	// Create request body
	body := RequestBody{Number: randomNum}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	newBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Print or process the response body
	fmt.Printf("Response Body: %s\n", string(newBody))

	resp.Body.Close()
	fmt.Printf("Request %d sent with Number: %d\n", requestId+1, randomNum)
}

func postCheckPrimeNumRequest(url string, requestId int, wg *sync.WaitGroup) {
	defer wg.Done()

	makePostRequestToServer(url, requestId)
}

func sendSerializedRequests() {
	seed := rand.NewSource(time.Now().UnixNano())
	rand.New(seed)

	url := "http://localhost:8080/prime"

	for i := 0; i < 100; i++ {
		makePostRequestToServer(url, i)
	}
}

func sendRequestsInParallel() {
	seed := rand.NewSource(time.Now().UnixNano())
	rand.New(seed)

	url := "http://localhost:8080/prime"

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go postCheckPrimeNumRequest(url, i, &wg)
	}

	wg.Wait()
	fmt.Println("All POST requests are sent to server in parallel")
}

func sendPostRequestViaConnPool(url string, requestId int, wg *sync.WaitGroup, connectionPoolChan chan struct{}) {
	defer wg.Done()
	connectionPoolChan <- struct{}{}
	makePostRequestToServer(url, requestId)
	<-connectionPoolChan
}

func sendRequestsUsingConnectionPool() {
	seed := rand.NewSource(time.Now().UnixNano())
	rand.New(seed)

	url := "http://localhost:8080/prime"
	var wg sync.WaitGroup
	var connectionPoolChan = make(chan struct{}, MAX_CONNECTION_POOL_SIZE)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go sendPostRequestViaConnPool(url, i, &wg, connectionPoolChan)
	}

	wg.Wait()
	fmt.Println("Done sending all POST requests via Client Connection Pool!")
}

func main() {
	// sendSerializedRequests is a function to make requests to the server in a serial way one after the other
	// sendSerializedRequests()

	// sendRequestsInParallel function uses go-routines to send all requests at once to the server.
	// sendRequestsInParallel()

	// sendRequestsUsingConnectionPool function uses connectionPool + go-routines to limit the number of connections
	//  that should be made to the server at all times, this helps to have more control over the traffic being sent on backend server from frontend
	sendRequestsUsingConnectionPool()
}
