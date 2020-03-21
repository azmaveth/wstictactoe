/*
 * MIT License
 *
 * Copyright (c) 2020.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/azmaveth/wstictactoe/pkg/board"
	"github.com/azmaveth/wstictactoe/pkg/player"
)

const (
	HOST = "ec2-3-94-79-118.compute-1.amazonaws.com"
	PORT = "42424"
	HOME_PAGE = "/"
	WS_ENDPOINT = "/ttt"
)

var addr = flag.String("addr", ":" + PORT, "HTTP Server Address")

var upgrader = websocket.Upgrader{}

func defaultHomePage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	
	if r.URL.Path != HOME_PAGE {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	fmt.Fprint(w, "Welcome to the Tic Tac Toe WebSocket Server\n\n")
	fmt.Fprint(w, "Please connect to ws://" + HOST + ":" + PORT + WS_ENDPOINT + " to reach the WebSocket.\n")
}

func websocketPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // TODO: Implement proper CORS validation
	}
	
	if r.URL.Path != WS_ENDPOINT {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during WebSocket upgrade: ", err)
		return
	}
	
	log.Println("Client connected from: " + r.RemoteAddr)
	defer c.Close()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket: ", err)
			break
		}
		
		log.Printf("Received from WebSocket: %s", message)

		results := map[string][3][3]player.Player{}
		err = json.Unmarshal(message, &results)
		if err != nil {
			log.Println("Error unmarshalling JSON: ", err)
			break
		}
		var b = board.NewBoard(results["Cells"])

		var replyMessage = checkGame(b)

		err = c.WriteMessage(messageType, []byte(replyMessage))
		if err != nil {
			log.Println("Error writing to WebSocket: ", err)
			break
		}
	}
}

func routeHttpRequests() {
	http.HandleFunc(HOME_PAGE, defaultHomePage)
	http.HandleFunc(WS_ENDPOINT, websocketPage)
}

func main() {
	flag.Parse()
	routeHttpRequests()
	log.Println("Listening for client...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func checkGame(b board.Board) string {
	if board.CheckForWinningPlayer(player.X, b) {
		return "X wins!"
	} else if board.CheckForWinningPlayer(player.O, b) {
		return "O wins!"
	} else if board.IsBoardFull(b) {
		return "No winner: cat's game."
	} else {
		return "Game not yet complete."
	}
}
