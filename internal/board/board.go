package board

import (
	"errors"
	"fmt"

	term "github.com/buger/goterm"
)

const (
	HEIGHT = 6
	WIDTH  = 7
)

type Board struct {
	board                [][]string
	bottomMostEmptySpots []int
	lastMove             []int
	emptySpotsCount      int
}

func New() *Board {
	return &Board{
		board: [][]string{
			{" ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " "},
		},
		bottomMostEmptySpots: []int{5, 5, 5, 5, 5, 5, 5},
		lastMove:             []int{-1, -1},
		emptySpotsCount:      HEIGHT * WIDTH,
	}
}

func (b *Board) Get(row int, col int) (string, error) {
	if row >= 0 && row < HEIGHT && col >= 0 && col < WIDTH {
		return b.board[row][col], nil
	}
	return "", errors.New("Not valid indices!")
}

func (b *Board) Put(col int, token string) error {
	row := b.bottomMostEmptySpots[col]
	if row == -1 {
		return errors.New("The column is full!")
	}
	b.board[row][col] = token
	b.lastMove = []int{row, col}
	b.bottomMostEmptySpots[col] -= 1
	b.emptySpotsCount -= 1
	return nil
}

func (b *Board) Remove(col int) error {
	row := b.bottomMostEmptySpots[col] + 1
	if row >= HEIGHT {
		return errors.New("No token in this column!")
	}
	b.board[row][col] = " "
	b.bottomMostEmptySpots[col]++
	b.emptySpotsCount++
	return nil
}

func (b *Board) AvailableMoves() []int {
	availableMoves := make([]int, 0)
	bottomMostEmptySpotsLength := len(b.bottomMostEmptySpots)
	start := bottomMostEmptySpotsLength / 2

	if b.bottomMostEmptySpots[start] > -1 {
		availableMoves = append(availableMoves, start)
	}
	for offset := 1; offset < (bottomMostEmptySpotsLength/2)+1; offset++ {
		if start-offset < bottomMostEmptySpotsLength && b.bottomMostEmptySpots[start-offset] > -1 {
			availableMoves = append(availableMoves, start-offset)
		}
		if start+offset < bottomMostEmptySpotsLength && b.bottomMostEmptySpots[start+offset] > -1 {
			availableMoves = append(availableMoves, start+offset)
		}
	}
	return availableMoves
}

func (b *Board) NumOfEmptySpots() int {
	return b.emptySpotsCount
}

func (b *Board) IsOver() bool {
	return b.IsFilled() || b.SomeoneWon()
}

func (b *Board) IsFilled() bool {
	return b.emptySpotsCount == 0
}

func (b *Board) SomeoneWon() bool {
	if b.lastMove[0] == -1 || b.lastMove[1] == -1 {
		return false
	}

	lastMoveRow := b.lastMove[0]
	lastMoveCol := b.lastMove[1]
	token := b.board[lastMoveRow][lastMoveCol]

	return b.wonByRow(lastMoveRow, lastMoveCol, token) || b.wonByCol(lastMoveRow, lastMoveCol, token) ||
		b.wonByMajorDiagonal(lastMoveRow, lastMoveCol, token) || b.wonByMinorDiagonal(lastMoveRow, lastMoveCol, token)
}

func (b *Board) wonByRow(row int, col int, token string) bool {
	return b.countByDirection(row, col, 0, -1, token)+b.countByDirection(row, col, 0, 1, token) >= 3
}

func (b *Board) wonByCol(row int, col int, token string) bool {
	return b.countByDirection(row, col, -1, 0, token)+b.countByDirection(row, col, 1, 0, token) >= 3
}

func (b *Board) wonByMajorDiagonal(row int, col int, token string) bool {
	return b.countByDirection(row, col, -1, -1, token)+b.countByDirection(row, col, 1, 1, token) >= 3
}

func (b *Board) wonByMinorDiagonal(row int, col int, token string) bool {
	return b.countByDirection(row, col, -1, 1, token)+b.countByDirection(row, col, 1, -1, token) >= 3
}

func (b *Board) countByDirection(row int, col int, rowOffset int, colOffset int, token string) int {
	isSameCount := 0
	for i := 1; i < 5; i++ {
		modifiedRow := row + (rowOffset * i)
		modifiedCol := col + (colOffset * i)
		if modifiedRow < 0 || modifiedRow >= HEIGHT || modifiedCol < 0 || modifiedCol >= WIDTH ||
			b.board[modifiedRow][modifiedCol] != token {
			break
		} else {
			isSameCount++
		}
	}
	return isSameCount
}

func (b *Board) Print() {
	term.Clear()
	term.MoveCursor(1, 1)
	term.Println(b)
	term.Flush()
}

func (b *Board) String() string {
	strBoard := ""

	for rowIndex, row := range b.board {
		for colIndex, token := range row {
			if rowIndex == b.lastMove[0] && colIndex == b.lastMove[1] {
				strBoard += fmt.Sprintf(" %s ", term.Color(token, term.GREEN))
			} else if token == " " {
				strBoard += "   "
			} else if token == "X" {
				strBoard += fmt.Sprintf(" %s ", term.Color(token, term.BLUE))
			} else if token == "O" {
				strBoard += fmt.Sprintf(" %s ", term.Color(token, term.RED))
			}
			if colIndex < WIDTH-1 {
				strBoard += "|"
			}
		}
		strBoard += "\n"
		if rowIndex < HEIGHT-1 {
			strBoard += "---------------------------"
			strBoard += "\n"
		}
	}
	return strBoard
}
