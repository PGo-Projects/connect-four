package utils

func GetOtherToken(chosenToken string) string {
	if chosenToken == "X" {
		return "O"
	} else {
		return "X"
	}
}
