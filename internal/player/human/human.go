package human

import (
	"fmt"
	"strconv"

	"github.com/PGo-Projects/connect-four/internal/board"
	"github.com/PGo-Projects/connect-four/internal/userio"
)

const (
	TYPE                    = "human"
	MOVE_VALIDATION_PATTERN = "[1234567]"
)

type Human struct {
	token string
}

func New(token string) *Human {
	return &Human{token: token}
}

func (h *Human) GetToken() string {
	return h.token
}

func (h *Human) GetType() string {
	return TYPE
}

func (h *Human) PlayMove(b *board.Board) error {
	addressMsg := fmt.Sprintf("To the player with token %s:", h.token)
	colPrompt := "Which column do you want to put your token?"
	colErrMsg := "Not a valid column, please try again!"

	col, colErr := strconv.Atoi(userio.PromptUser(&userio.PromptUserInfo{
		AddressMsg:                 addressMsg,
		PromptMsg:                  colPrompt,
		UserResponseIsValidPattern: MOVE_VALIDATION_PATTERN,
		ErrMsg: colErrMsg,
	}))
	if colErr == nil {
		playErr := b.Put(col-1, h.token)
		return playErr
	}
	panic("This should not happen!")
}
