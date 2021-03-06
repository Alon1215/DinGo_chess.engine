package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var tell = mainTell
var trim = strings.TrimSpace
var low = strings.ToLower

var saveBm = ""

func uci(input chan string) {
	toEng, frEng := Engine()
	var cmd string
	var bm string
	quit := false
	for !quit {
		select {
		case cmd = <-input:
			//			tell("info string uci got ", cmd, "\n")
		case bm = <-frEng:
			handleBm(bm)
			continue
		}
		words := strings.Split(cmd, " ")
		words[0] = trim(low(words[0]))

		switch words[0] {
		case "uci":
			handleUci()
		case "setoption":
			handleSetOption(words)
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleNewgame()
		case "position":
			handlePosition(cmd)
		case "debug":
			handleDebug(words)
		case "register":
			handleRegister(words)
		case "go":
			handleGo(toEng, words)
		case "ponderhit":
			handlePonderhit()
		case "stop":
			handleStop()
		case "quit", "q":
			handleQuit()
			quit = true
			continue
		//////// Added commands, not part of the uci protocol ////////
		case "perft":
			if len(words) > 1 {
				depth, err := strconv.Atoi(words[1])
				if err != nil {
					tell(err.Error())
				} else {
					startPerft(depth, &board)
				}
			}
		case "pb":
			board.Print()
		case "pbb":
			board.printAllBB()
		case "pm":
			board.printAllLegals()
		case "eval":
			fmt.Println("eval =", evaluate(&board))
		case "pos":
			handleMyPositions(words)
		case "moves":
			handleMyMoves(words)
		case "key":
			fmt.Printf("key = %x, fullkey=%x\n", board.key, board.fullKey())
			index := board.fullKey() & uint64(trans.mask)
			lock := trans.lock(board.fullKey())
			fmt.Printf("index = %x, lock=%x\n", index, lock)
		case "see":
			fr, to := empty, empty
			if len(words[1]) == 2 && len(words[2]) == 2 {
				fr = fen2Sq[words[1]]
				to = fen2Sq[words[2]]
			} else if len(words[1]) == 4 {
				fr = fen2Sq[words[1][0:2]]
				to = fen2Sq[words[1][2:]]
			} else {
				fmt.Println("error in fr/to")
				continue
			}

			fmt.Println("see = ", see(fr, to, &board))
		case "qs":
			fmt.Println("qs =", qs(maxEval, &board))
		case "hist":
			history.print(50)
		case "moveval": // all moves and values
			//handleMoveVal() // NOT IMPLEMNTED, (just print command)
		case "limits":
			fmt.Println(limits)
		default:
			tell("info string unknown cmd ", cmd)
		}
	}

	tell("info string leaving uci()")
}

func handleUci() {
	tell("id name DinGo")
	tell("id author Alon Michaeli")

	tell("option name Hash type spin default 128 min 16 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}

func handleNewgame() {
	board.newGame()
	history.clear()
}

func handlePosition(cmd string) {
	// position [fen <fenstring> | startpos ]  moves <move1> .... <movei>
	board.newGame()
	cmd = trim(strings.TrimPrefix(cmd, "position"))
	parts := strings.Split(cmd, "moves")
	if len(cmd) == 0 || len(parts) > 2 {
		err := fmt.Errorf("%v wrong length=%v", parts, len(parts))
		tell("info string Error", fmt.Sprint(err))
		return
	}

	alt := strings.Split(parts[0], " ")
	alt[0] = trim(alt[0])
	if alt[0] == "startpos" {
		parts[0] = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	} else if alt[0] == "fen" {
		parts[0] = trim(strings.TrimPrefix(parts[0], "fen"))
	} else {
		err := fmt.Errorf("%#v must be %#v or %#v", alt[0], "fen", "startpos")
		tell("info string Error", err.Error())
		return
	}
	// Now parts[0] is the fen-string only

	// start the parsing
	//fmt.Printf("info string parse %#v\n", parts[0])
	parseFEN(parts[0])

	if len(parts) == 2 {
		parts[1] = low(trim(parts[1]))
		//fmt.Printf("info string parse %#v\n", parts[1])
		parseMvs(parts[1])
	}
}

func handleStop() {
	if limits.Infinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}

		limits.setInfinite(false)
	}
	limits.setStop(true)
}

// handleQuit not really necessary
func handleQuit() {

}

func handleBm(bm string) {
	if limits.Infinite {
		saveBm = bm
		return
	}
	tell(bm)
}

func handleSetOption(words []string) {
	// setoption name Hash value 256
	if len(words) < 5 {
		tell("info string Don't have this option " + strings.Join(words[:], " "))
	}
	if low(trim(words[1])) != "name" {
		tell("info string 'name' is missing in this option " + strings.Join(words[:], " "))
	}
	switch low(trim(words[2])) {
	case "hash":
		if trim(low(words[3])) != "value" {
			tell("info string 'value' is missing in this option " + strings.Join(words[:], " "))
		}
		if val, err := strconv.Atoi(trim(words[4])); err == nil {
			if err = trans.new(val); err != nil {
				tell(err.Error())
			}
		} else {
			tell("info string The Hash value is not numeric " + strings.Join(words[:], " "))
		}
	default:
		tell("info string Don't have this option " + strings.Join(words[:], " "))
	}
}

// go  searchmoves <move1-moveii>/ponder/wtime <ms>/ btime <ms>/winc <ms>/binc <ms>/movestogo <x>/
//     depth <x>/nodes <x>/movetime <ms>/mate <x>/infinite
func handleGo(toEng chan bool, words []string) {

	ok := limits.ParseLimits(words[1:])
	if !ok {
		tell("go command invalid")
	} else {
		toEng <- true
	}

}

func handleMyPositions(words []string) {
	if len(words) < 2 {
		tell("info string not correct pos command " + strings.Join(words[:], " "))
	} else {

		words[1] = trim(low(words[1]))
		handleSetOption(strings.Split("setoption name hash value 256", " "))
		switch words[1] {
		case "london":
			handlePosition("position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 c7c5 b1d2 b8c6 c2c3 e7e6 f1d3 f8d6")
		case "phil":
			handlePosition("position startpos moves e2e4 d7d6 d2d4 e7e5 d4e5 d6e5 d1d8 e8d8 g1f3 f7f6 b1c3 c7c6 f1c4")
		case "english":
			handlePosition("position startpos moves c2c4 e7e5 g2g3 b8c6 f1g2 g7g6 b1c3 f8g7 e2e4 d7d6 g1e2 g8f6")
		default:
			tell("info string not correct pos command " + words[1] + " doesn't exist. " + strings.Join(words[:], " "))
		}
	}
}

// not implemented uci commands
func handlePonderhit() {
	tell("info string ponderhit not implemented")
}

func handleDebug(words []string) {
	// debug [ on | off ]
	tell("info string debug not implemented")
}

func handleRegister(words []string) {
	// register later/name <x>/code <y>
	tell("info string register not implemented")
}

func handleMyMoves(words []string) {
	mvString := strings.Join(words[1:], " ")
	parseMvs(mvString)
}

//------------------------------------------------------
func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	var reader *bufio.Reader
	reader = bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}
