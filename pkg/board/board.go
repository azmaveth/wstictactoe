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

package board

import "github.com/azmaveth/wstictactoe/pkg/player"

type Board struct {
	Cells [3][3]player.Player
}

func NewBoard(cells [3][3]player.Player) Board {
	return Board{cells}
}

func IsBoardFull(b Board) bool {
	var row = 0
	var col = 0
	var foundBlank = false
	for !foundBlank && row < 3 {
		for !foundBlank && col < 3 {
			if b.Cells[row][col] == player.Blank {
				foundBlank = true
			}
			col++
		}
		row++
	}
	return !foundBlank
}

func CheckForWinningPlayer(p player.Player, b Board) bool {
	return checkRows(p, b) ||
		checkColumns(p, b) ||
		checkDiagonals(p, b)
}

func checkRow(p player.Player, b Board, col int) bool {
	return b.Cells[0][col] == p &&
		b.Cells[1][col] == p &&
		b.Cells[2][col] == p
}

func checkRows(p player.Player, b Board) bool {
	return checkRow(p, b, 0) ||
		checkRow(p, b, 1) ||
		checkRow(p, b, 2)
}

func checkColumn(p player.Player, b Board, row int) bool {
	return b.Cells[row][0] == p &&
		b.Cells[row][1] == p &&
		b.Cells[row][2] == p
}

func checkColumns(p player.Player, b Board) bool {
	return checkColumn(p, b, 0) ||
		checkColumn(p, b, 1) ||
		checkColumn(p, b, 2)
}

func checkDiagonal1(p player.Player, b Board) bool {
	return b.Cells[0][0] == p &&
		b.Cells[1][1] == p &&
		b.Cells[2][2] == p
}

func checkDiagonal2(p player.Player, b Board) bool {
	return b.Cells[0][2] == p &&
		b.Cells[1][1] == p &&
		b.Cells[2][0] == p
}

func checkDiagonals(p player.Player, b Board) bool {
	return checkDiagonal1(p, b) ||
		checkDiagonal2(p, b)
}
