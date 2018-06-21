package userio

import (
	"fmt"
	"regexp"

	term "github.com/buger/goterm"
)

type PromptUserInfo struct {
	AddressMsg                 string
	PromptMsg                  string
	UserResponseIsValidPattern string
	ErrMsg                     string
}

func PromptUser(info *PromptUserInfo) string {
	if info.AddressMsg != "" {
		term.Println(info.AddressMsg)
	}
	term.Println(info.PromptMsg)
	term.MoveCursorUp(1)
	term.Flush()
	userResponse := ""
	match, err := regexp.MatchString(info.UserResponseIsValidPattern, userResponse)
	for err != nil || !match {
		if userResponse != "" {
			term.Println(info.ErrMsg)
			term.Flush()
			userResponse = ""
		}
		for userResponse == "" {
			fmt.Scanln(&userResponse)
		}
		match, err = regexp.MatchString(info.UserResponseIsValidPattern, userResponse)
	}
	term.Print("\n")
	term.Flush()
	return userResponse
}
