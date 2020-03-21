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
	"github.com/azmaveth/wstictactoe/pkg/board"
	"github.com/azmaveth/wstictactoe/pkg/player"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

const (
	HOST = "ec2-3-94-79-118.compute-1.amazonaws.com"
	PORT = "42424"
	WS_ENDPOINT = "/ttt"
)

var addr = flag.String("addr", HOST + ":" + PORT, "HTTP Server Address")

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: WS_ENDPOINT}
	log.Printf("Attempting to connect to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket: ", err)
	}
	defer c.Close()

	// No winner yet, board not full
	var b = board.NewBoard([3][3]player.Player{
		{player.O, player.Blank, player.O},
		{player.O, player.X, player.X},
		{player.X, player.X, player.O}})
	sendBoard(b, c)

	// No winner, board full
	b = board.NewBoard([3][3]player.Player{
		{player.O, player.X, player.O},
		{player.O, player.X, player.X},
		{player.X, player.O, player.O}})
	sendBoard(b, c)

	// X wins
	b = board.NewBoard([3][3]player.Player{
		{player.O, player.X, player.O},
		{player.O, player.X, player.X},
		{player.X, player.X, player.O}})
	sendBoard(b, c)

	// O wins
	b = board.NewBoard([3][3]player.Player{
		{player.O, player.O, player.O},
		{player.O, player.X, player.X},
		{player.X, player.X, player.O}})
	sendBoard(b, c)
}

func sendBoard(b board.Board, c *websocket.Conn) {
	var bJson, _ = json.Marshal(b)

	log.Println("Sending: " + string(bJson))

	err := c.WriteMessage(websocket.TextMessage, bJson)
	if err != nil {
		log.Fatal("Error sending message: ", err)
	}

	_, responseMessage, err := c.ReadMessage()

	log.Println("Response: " + string(responseMessage))
}
