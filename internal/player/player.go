package player

import (
	"github.com/PGo-Projects/connect-four/internal/board"
	"github.com/PGo-Projects/connect-four/internal/player/ai"
	"github.com/PGo-Projects/connect-four/internal/player/human"
)

const (
	AI    = ai.TYPE
	HUMAN = human.TYPE
)

type Player interface {
	GetToken() string
	GetType() string
	PlayMove(*board.Board) error
}

func New(playerType string, token string) Player {
	if playerType == HUMAN {
		return human.New(token)
	}
	return ai.New(token)
}
