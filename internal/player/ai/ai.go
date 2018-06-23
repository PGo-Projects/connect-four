package ai

import (
	"errors"
	"log"
	"math"

	"github.com/PGo-Projects/connect-four/internal/board"
	"github.com/PGo-Projects/connect-four/internal/utils"
	term "github.com/buger/goterm"
)

const (
	TYPE                   = "ai"
	OUT_OF_BOUND_ERROR_MSG = "Out of bound!"
)

type Ai struct {
	token string
}

func New(token string) *Ai {
	return &Ai{token: token}
}

func (a *Ai) GetToken() string {
	return a.token
}

func (a *Ai) GetType() string {
	return TYPE
}

func (a *Ai) PlayMove(b *board.Board) error {
	term.Println(term.Color("Please wait a moment while I think about my next move...", term.GREEN))
	term.Flush()
	if b.NumOfEmptySpots() > 27 {
		bestMove, _ := a.minimax(b, 0, math.MinInt32, math.MaxInt32, 1, 10)
		return b.Put(bestMove, a.token)
	} else {
		bestMove, _ := a.minimax(b, 0, math.MinInt32, math.MaxInt32, 1, 15)
		return b.Put(bestMove, a.token)
	}
}

func (a *Ai) minimax(b *board.Board, turn int, alpha int, beta int, depth int, searchLimit int) (int, int) {
	if b.IsOver() || depth > searchLimit {
		return evaluate(b, depth, turn)
	}

	availableMoves := b.AvailableMoves()
	moves := make([]int, 0)
	scores := make([]int, 0)
	bestMove := -1
	for _, col := range availableMoves {
		if turn == 0 {
			b.Put(col, a.token)
		} else {
			b.Put(col, utils.GetOtherToken(a.token))
		}
		moves = append(moves, col)
		_, score := a.minimax(b, 1-turn, alpha, beta, depth+1, searchLimit)
		scores = append(scores, score)
		if turn == 0 {
			bestScore, index := getBestScore(scores, 0)
			bestMove = moves[index]
			if bestScore > alpha {
				alpha = bestScore
			}
		} else {
			bestScore, _ := getBestScore(scores, 1)
			if bestScore < beta {
				beta = bestScore
			}
		}
		b.Remove(col)
		if alpha >= beta {
			break
		}
	}

	if turn == 0 {
		return bestMove, alpha
	} else {
		return bestMove, beta
	}
}

func evaluate(b *board.Board, depth int, turn int) (int, int) {
	if b.SomeoneWon() {
		if turn == 0 {
			return -1, depth - math.MaxInt32
		} else {
			return -1, math.MaxInt32 - depth
		}
	}
	if b.IsFilled() {
		return -1, 0
	}
	return -1, linkedRowStrength(b, turn)
}

func getLength(b *board.Board, position []int, dirOffset []int, token string, visited *[][]bool, count int) int {
	row := position[0] + count*dirOffset[0]
	col := position[1] + count*dirOffset[1]

	tokenOnBoard, err := b.Get(row, col)
	for err != nil && tokenOnBoard == token {
		(*visited)[row][col] = true
		count++
		row += dirOffset[0]
		col += dirOffset[1]
		tokenOnBoard, err = b.Get(row, col)
	}
	return count
}

func getLengthPair(b *board.Board, position []int, dirOffset []int, token string, visited *[][]bool) (int, int) {
	position = []int{position[0] + dirOffset[0], position[1] + dirOffset[1]}
	playerLength := getLength(b, position, dirOffset, token, visited, 0)
	possibleLength := getLength(b, position, dirOffset, " ", visited, playerLength)
	return playerLength, possibleLength
}

func getRating(b *board.Board, position []int, directions [][]int, token string, visited *[][]bool, turn int) int {
	rating := 0
	lineWeights := []int{0, 10, 100}
	for i := 0; i < 4; i++ {
		posDirection := []int{directions[0][i], directions[1][i]}
		posPlayerLength, posPossibleLength := getLengthPair(b, position, posDirection, token, visited)
		negDirection := []int{-directions[0][i], -directions[1][i]}
		negPlayerLength, negPossibleLength := getLengthPair(b, position, negDirection, token, visited)
		if posPossibleLength+negPossibleLength >= 3 {
			rating += -turn * lineWeights[posPlayerLength+negPlayerLength]
		}
	}
	return rating
}

func linkedRowStrength(b *board.Board, turn int) int {
	total := 0
	directions := [][]int{
		{1, 1, 0, -1},
		{0, 1, 1, 1},
	}
	visited := [][]bool{
		{false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false},
	}

	for row := 0; row < board.HEIGHT; row++ {
		for col := 0; col < board.WIDTH; col++ {
			if visited[row][col] {
				continue
			}
			token, err := b.Get(row, col)
			if err != nil {
				log.Fatalf("Invalid row {} or column {}", row, col)
			}
			if token != " " {
				total += getRating(b, []int{row, col}, directions, token, &visited, turn)
			}
			visited[row][col] = true
		}
	}
	return total
}

func getBestScore(scores []int, turn int) (int, int) {
	if turn == 0 {
		scoreIndex, err := getMaxIndexOfSlice(scores)
		if err == nil {
			return scores[scoreIndex], scoreIndex
		}
	} else {
		scoreIndex, err := getMinIndexOfSlice(scores)
		if err == nil {
			return scores[scoreIndex], scoreIndex
		}
	}
	return -1, -1
}

func getMinIndexOfSlice(slice []int) (minIndex int, err error) {
	if len(slice) == 0 {
		return -1, errors.New(OUT_OF_BOUND_ERROR_MSG)
	}
	currentMin := slice[0]
	for index, num := range slice {
		if num < currentMin {
			currentMin = num
			minIndex = index
		}
	}
	return minIndex, nil
}

func getMaxIndexOfSlice(slice []int) (maxIndex int, err error) {
	if len(slice) == 0 {
		return -1, errors.New(OUT_OF_BOUND_ERROR_MSG)
	}
	currentMax := slice[0]
	for index, num := range slice {
		if num > currentMax {
			currentMax = num
			maxIndex = index
		}
	}
	return maxIndex, nil
}
