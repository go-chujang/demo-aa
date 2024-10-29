package main

import (
	"strconv"
)

func progressGamble() (msg string, ok bool) {
	var err error
	switch SEQ_PROGRESS {
	case 0:
		return "1) rock-paper-scissors  2) zero-to-nine  3) exchange", true
	case 1:
		STORED_PARAM["gamble"] = LAST_INPUT
		switch LAST_INPUT {
		case "1":
			return "type 1) rock  2) paper  3) scissors", true
		case "2":
			return "type between 0 and 9", true
		case "3":
			return "type tokenId 1000 or 2000", true
		default:
			return "type 1 or 2 or 3", true
		}
	case 2:
		switch STORED_PARAM["gamble"] {
		case "1":
			var choice string
			switch LAST_INPUT {
			case "1":
				choice = "rock"
			case "2":
				choice = "paper"
			case "3":
				choice = "scissors"
			default:
				return "type 1 or 2 or 3", true
			}
			err = gambleRPS(choice)
		case "2":
			var guess uint64
			guess, err = strconv.ParseUint(LAST_INPUT, 10, 64)
			if err == nil {
				err = gambleZ2N(guess)
			}
		case "3":
			err = gambleExchange(LAST_INPUT)
		}
		DONE_PROGRESS = true
	}
	if err != nil {
		return err.Error(), false
	}
	return "type 'state' to check your balance", true
}

func gambleRPS(choice string) error {
	return postUserOp("/svc/v1/users/operations/gambles/rock-paper-scissors", map[string]interface{}{
		"choice": choice,
	})
}

func gambleZ2N(guess uint64) error {
	return postUserOp("/svc/v1/users/operations/gambles/zero-to-nine", map[string]interface{}{
		"guess": guess,
	})
}

func gambleExchange(tokenId string) error {
	return postUserOp("/svc/v1/users/operations/gambles/exchange", map[string]interface{}{
		"tokenId": tokenId,
	})
}
