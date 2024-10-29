package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-chujang/demo-aa/common/net/httpf"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
)

const (
	// URL           = "http://localhost:5000"
	URL           = "http://172.29.126.229:80"
	privateKeyHex = "0x52a01f1560d12c932a72dfa1730de107be0ff1c9c22e42d89b9a2a78d9b2c467"

	PROGRESS_LOGIN    = "login"
	PROGRESS_JOIN     = "join"
	PROGRESS_STATE    = "state"
	PROGRESS_FAUCET   = "faucet"
	PROGRESS_TRANSFER = "transfer"
	PROGRESS_GAMBLE   = "gamble"
)

var (
	PRIVATE_KEY, pverr = ethutil.PvHex2Key(privateKeyHex)
	ADDRESS_EOA, _     = ethutil.PvKey2Address(PRIVATE_KEY)
	LOGGED_IN          bool
	IN_PROGRESS        string
	SEQ_PROGRESS       int
	DONE_PROGRESS      bool
	LAST_INPUT         string
	STORED_PARAM       = make(map[string]string)
)

func main() {
	if pverr != nil {
		panic(pverr)
	}
	httpf.SetTimeout(time.Second * 100)
	fmt.Println("demo-client 'exit' to quit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		var message string
		switch {
		case IN_PROGRESS == "" && !LOGGED_IN:
			message = typeGuide(PROGRESS_LOGIN, PROGRESS_JOIN)
		case IN_PROGRESS == "" && LOGGED_IN:
			message = typeGuide(PROGRESS_STATE, PROGRESS_FAUCET, PROGRESS_TRANSFER, PROGRESS_GAMBLE)
		default:
			// progress selector
			// progress handler return false => not fail, reset progress
			var ok bool
			switch IN_PROGRESS {
			case PROGRESS_LOGIN, PROGRESS_JOIN:
				message, ok = progressLoginOrJoin()
			case PROGRESS_STATE:
				message, ok = state()
			case PROGRESS_FAUCET:
				message, ok = faucet()
			case PROGRESS_TRANSFER:
				message, ok = progressTransfer()
			case PROGRESS_GAMBLE:
				message, ok = progressGamble()
			}
			if DONE_PROGRESS || !ok {
				DONE_PROGRESS = false
				fmt.Println(message)
				resetProgress()
				continue
			}
			if ok {
				SEQ_PROGRESS++
			}
		}
		waitForInput(message)

		input, err := reader.ReadString('\n')
		if err != nil || !parseInput(&input) {
			resetProgress()
			continue
		}
		if input == "exit" {
			return
		}
		LAST_INPUT = input

		if IN_PROGRESS == "" {
			switch input {
			case PROGRESS_JOIN, PROGRESS_LOGIN:
			case PROGRESS_STATE, PROGRESS_FAUCET, PROGRESS_TRANSFER, PROGRESS_GAMBLE:
				if !LOGGED_IN {
					resetProgress()
				}
			default:
				continue
			}
			IN_PROGRESS = input
		}
	}
}

/////////////////////////////////////////////////////////
//	utils

func resetProgress() {
	IN_PROGRESS = ""
	SEQ_PROGRESS = 0
	LAST_INPUT = ""
}

func typeGuide(list ...string) string {
	withQuotations := make([]string, 0, len(list))
	for _, v := range list {
		withQuotations = append(withQuotations, fmt.Sprintf("'%s'", v))
	}
	return fmt.Sprintf("type: %s", strings.Join(withQuotations, " or "))
}

func waitForInput(msg string) {
	if msg != "" {
		fmt.Println(msg)
		fmt.Print("> ")
	}
}

func parseInput(input *string) bool {
	parsed := conv.TrimSpacePrefix(*input)
	if runtime.GOOS == "windows" {
		parsed = strings.TrimSuffix(parsed, "\r\n")
	} else {
		parsed = strings.TrimSuffix(parsed, "\n")
	}
	*input = parsed
	return parsed != ""
}
