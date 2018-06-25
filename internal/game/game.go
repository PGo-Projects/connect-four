package game

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/PGo-Projects/connect-four/internal/board"
	"github.com/PGo-Projects/connect-four/internal/player"
	"github.com/PGo-Projects/connect-four/internal/userio"
	"github.com/PGo-Projects/connect-four/internal/utils"
	term "github.com/buger/goterm"
)

type Game struct {
	board   *board.Board
	players []player.Player
	turn    int
}

func New() *Game {
	return &Game{board: board.New()}
}

func (g *Game) Start() {
	term.Clear()
	term.MoveCursor(1, 1)
	PrintInstructions()

	g.determinePlayers()
	g.turn = 0
	g.drawBoard()
	g.play()
}

func (g *Game) drawBoard() {
	g.board.Print()
}

func (g *Game) play() {
	if g.board.IsOver() {
		term.Println("Game Over!")
		term.Flush()
		return
	}

	playErr := g.getCurrentPlayer().PlayMove(g.board)
	g.drawBoard()
	if playErr != nil {
		term.Println(playErr)
		term.Flush()
	} else {
		g.nextTurn()
	}

	g.play()
}

func (g *Game) getCurrentPlayer() player.Player {
	return g.players[g.turn]
}

func (g *Game) nextTurn() {
	g.turn = 1 - g.turn
}

func (g *Game) determinePlayers() {
	userWishToPlayAgainstAI := userio.PromptUser(&userio.PromptUserInfo{
		AddressMsg:                 "",
		PromptMsg:                  "Want to play against an AI?",
		UserResponseIsValidPattern: "([yY][eE][sS])|([nN][oO])",
		ErrMsg: "Please answer yes or no.",
	})
	if match, err := regexp.MatchString("[yY][eE][sS]", userWishToPlayAgainstAI); err == nil && match {
		playerWishToGoesFirst := userio.PromptUser(&userio.PromptUserInfo{
			AddressMsg:                 "",
			PromptMsg:                  "Do you want to go first?",
			UserResponseIsValidPattern: "([yY][eE][sS])|([nN][oO])",
			ErrMsg: "Please answer yes or no.",
		})
		if match, err := regexp.MatchString("[yY][eE][sS]", playerWishToGoesFirst); err == nil && match {
			playerToken := userio.PromptUser(&userio.PromptUserInfo{
				AddressMsg:                 "",
				PromptMsg:                  "What token do you want: X or O?",
				UserResponseIsValidPattern: "X|O",
				ErrMsg: "Token must be a X or a O, please try again!",
			})
			g.players = []player.Player{
				player.New(player.HUMAN, playerToken),
				player.New(player.AI, utils.GetOtherToken(playerToken)),
			}
		} else {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < 50 {
				g.players = []player.Player{player.New(player.AI, "X"), player.New(player.HUMAN, utils.GetOtherToken("X"))}
			} else {
				g.players = []player.Player{player.New(player.AI, "O"), player.New(player.HUMAN, utils.GetOtherToken("O"))}
			}
		}
	} else {
		firstPlayerToken := userio.PromptUser(&userio.PromptUserInfo{
			AddressMsg:                 "",
			PromptMsg:                  "What token should the first player use?",
			UserResponseIsValidPattern: "X|O",
			ErrMsg: "Token must be a X or a O, please try again!",
		})
		g.players = []player.Player{
			player.New(player.HUMAN, firstPlayerToken),
			player.New(player.HUMAN, utils.GetOtherToken(firstPlayerToken)),
		}
	}
}

func PrintInstructions() {
	term.Println("This is a game of connect four!")
	term.Flush()
}
