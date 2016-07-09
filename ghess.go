/*
Go Chess Engine - Ghess
Fenimore Love 2016
GPLv3

TODO: Search and Evaluation
TODO: Fen PGN reading
TODO: Fen output
*/
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// The chessboard type
type Board struct {
	board []byte // piece position
	// Game Variables
	castle    []byte // castle possibility KQkq or ----
	empassant int    // square vulnerable to empassant
	score     string
	toMove    string // Next move is w or b
	moves     int    // the count of moves
	check     bool
	// Map for display grid
	pgnMap map[string]int // the pgn format
	//pieceMap map[int]string    // coord to standard notation
	pieces map[string]string // the unicode fonts
	// Game Positions
	fen     string         // Game position
	pgn     string         // Game history
	headers string         // Pgn format
	pattern *regexp.Regexp // For parsing PGN

}

// Create a new Board in the starting position
// Initialize starting values
func NewBoard() Board {
	b := make([]byte, 120)
	fmt.Println("Initializing new Chess game\n")

	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)

	// Printed Board Notations
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'

	// Map of PGN notation
	m := make(map[string]int)
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"] = 18, 17, 16, 15, 14, 13, 12, 11
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"] = 28, 27, 26, 25, 24, 23, 22, 21
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"] = 38, 37, 36, 35, 34, 33, 32, 31
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"] = 48, 47, 46, 45, 44, 43, 42, 41
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"] = 58, 57, 56, 55, 54, 53, 52, 51
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"] = 68, 67, 66, 65, 64, 63, 62, 61
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"] = 78, 77, 76, 75, 74, 73, 72, 71
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"] = 88, 87, 86, 85, 84, 83, 82, 81

	// Todo make map for pieceMap[]

	// Map of unicode fonts
	r := make(map[string]string)
	r["p"], r["P"] = "\u2659", "\u265F"
	r["b"], r["B"] = "\u2657", "\u265D"
	r["n"], r["N"] = "\u2658", "\u265E"
	r["r"], r["R"] = "\u2656", "\u265C"
	r["q"], r["Q"] = "\u2655", "\u265B"
	r["k"], r["K"] = "\u2654", "\u265A"
	r["."] = "\u00B7"

	// Regex Pattern for matching pgn moves
	pattern, _ := regexp.Compile(`([PNBRQK]?[a-h]?[1-8]?)x?([a-h][1-8])([\+\?\!]?)|O(-?O){1,2}`)
	return Board{
		board:   b,
		castle:  []byte(`KQkq`),
		pgnMap:  m,
		pieces:  r,
		toMove:  "w",
		score:   "*",
		moves:   1,
		pattern: pattern,
	}
}

// Return PNG String
func (b *Board) PgnString() string {
	return b.headers + b.pgn
}

// Create printable board
func (b *Board) String() string {
	// TODO Rotate Board
	game := b.board
	p := b.pieces
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i]) + ": "
			continue
		}
		printBoard += "|" + p[string(game[i])] + "|"
	}

	printBoard += "\n"
	printBoard += "   :a::b::c::d::e::f::g::h:\n"
	return printBoard
}

// Wrapper in for standard notation positions
func (b *Board) standardWrapper(orig, dest string) error {
	// TODO: use two coordinates, include piece value
	//e2e4
	e := b.Move(b.pgnMap[orig], b.pgnMap[dest])
	if e != nil {
		return e
	}
	return nil
}

// Validate move
// Change byte values to new values
func (b *Board) Move(orig, dest int) error {
	val := b.board[orig]
	var o byte         // supposed starting square
	var d byte         // supposed destination
	var empassant bool //refactor?
	var isCastle bool
	if b.toMove == "w" {
		// check that orig is Upper
		// and dest is Enemy or Empty
		o = []byte(bytes.ToUpper(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToLower(b.board[dest : dest+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		// and dest is Enemy or Empty
		o = []byte(bytes.ToLower(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToUpper(b.board[dest : dest+1]))[0]
	}
	// Check for Castle
	if orig == 14 {
		isCastle = b.board[dest] == 'R'
	} else if orig == 84 {
		isCastle = b.board[dest] == 'r'
	}

	err := b.basicValidation(orig, dest, o, d, isCastle)
	if err != nil {
		return err
	}

	p := string(bytes.ToUpper(b.board[orig : orig+1]))
	switch {
	case p == "P":
		e := b.validPawn(orig, dest, d)
		if e != nil {
			return e
		}
		emp := dest - orig
		if emp > 11 || emp < -11 {
			empassant = true
		}
	case p == "N":
		e := b.validKnight(orig, dest)
		if e != nil {
			return e
		}
	case p == "B":
		e := b.validBishop(orig, dest)
		if e != nil {
			return e
		}
	case p == "R":
		e := b.validRook(orig, dest)
		if e != nil {
			return e
		}
	case p == "Q":
		e := b.validQueen(orig, dest)
		if e != nil {
			return e
		}
	case p == "K": // is castle?
		if !isCastle {
			e := b.validKing(orig, dest, false)
			if e != nil {
				return e
			}
			if orig == 14 || orig == 84 { // starting pos
				switch {
				case o == 'K':
					b.castle[0], b.castle[1] = '-', '-'
				case o == 'k':
					b.castle[2], b.castle[3] = '-', '-'
				}
			}
		} else {
			e := b.validKing(orig, dest, true)
			if e != nil {
				return e
			}
		}
	}
	// Make sure new position doesn't put in check
	isWhite := b.toMove == "w"
	possible := *b                 // slices are  still pointing...
	boardCopy := make([]byte, 120) // b.board is Pointer
	castleCopy := make([]byte, 4)
	copy(boardCopy, b.board)
	copy(castleCopy, b.castle)
	possible.board = boardCopy
	possible.castle = castleCopy
	// Check possibilities
	possible.updateBoard(orig, dest, val, empassant, isCastle)
	// find mover's king
	var king int
	for idx, val := range possible.board {
		if isWhite && val == 'K' {
			king = idx
			break
		} else if !isWhite && val == 'k' {
			king = idx
			break
		}
	}
	isCheck := possible.isInCheck(king)
	if isCheck {
		return errors.New("Cannot move into Check")
	}
	if isCastle {
		copy2 := make([]byte, 120)
		copy(copy2, b.board)
		possible.board = copy2
		if isWhite && dest < orig {
			//King side, 13
			possible.updateBoard(orig, 13, 'K',
				false, false)
			king = 13
		} else if isWhite && dest > orig {
			possible.updateBoard(orig, 15, 'K',
				false, false)
			// Queen side, 15
			king = 15
		} else if !isWhite && dest < orig {
			possible.updateBoard(orig, 83, 'k',
				false, false)
			// King 83
			king = 83
		} else if !isWhite && dest > orig {
			possible.updateBoard(orig, 85, 'k',
				false, false)
			// Queen 85
			king = 85
		}
		isCheck = possible.isInCheck(king)
		if isCheck {
			return errors.New("Cannot Castle through check")
		}
	}
	// update real board
	b.updateBoard(orig, dest, val, empassant, isCastle)
	return nil
}

// Updates board, useless without Move() validation
func (b *Board) updateBoard(orig, dest int,
	val byte, empassant, isCastle bool) {
	isWhite := b.toMove == "w"
	// Check for Promotion
	isPromotion := false
	if b.board[orig] == 'p' && dest < 20 {
		// is promotion
		isPromotion = true
	} else if b.board[orig] == 'P' && dest > 80 {
		isPromotion = true
	}
	// Check for castle deactivation
	if b.board[orig] == 'r' || b.board[orig] == 'R' {
		switch { // Castle
		case orig == b.pgnMap["a1"]:
			b.castle[1] = '-'
		case orig == b.pgnMap["a8"]:
			b.castle[3] = '-'
		case orig == b.pgnMap["h1"]:
			b.castle[0] = '-'
		case orig == b.pgnMap["h8"]:
			b.castle[2] = '-'
		}
	} else if isCastle {
		kingSide  := orig > dest
		queenSide := orig < dest
		switch {
		case isWhite && kingSide:
			b.castle[0], b.castle[1] = '-', '-'
		case isWhite && queenSide:
			b.castle[0], b.castle[1] = '-', '-'
		case !isWhite && kingSide:
			b.castle[2], b.castle[3] = '-', '-'
		case !isWhite && queenSide:
			b.castle[2], b.castle[3] = '-', '-'
		}
	}
	// Set origin
	b.board[orig] = '.'
	// Set destination
	if isCastle {
		if dest > orig { // queen side
			b.board[dest-2],
				b.board[dest-3] = val, b.board[dest]
		} else { // king side
			b.board[dest+1],
				b.board[dest+2] = val, b.board[dest]
		}
		b.board[dest] = '.'
	} else if isPromotion {
		switch {
		case dest < 20:
			b.board[dest] = 'q'
		case dest > 80:
			b.board[dest] = 'Q'
		}
	} else { // Normal Move/Capture
		b.board[dest] = val
	}
	// TODO check for Check
	// Update Game variables
	if b.toMove == "w" {
		b.toMove = "b"
	} else {
		b.moves++ // add one to move count
		b.toMove = "w"
	}
	if empassant {
		b.empassant = dest
	} else {
		b.empassant = 0
	}

	// Check if move put other player in Check
	isCheck := b.isPlayerInCheck()
	if isCheck {
		b.check = true
	} else {
		b.check = false
	}
}

// Check if current player is in Check
func (b *Board) isPlayerInCheck() bool {
	isWhite := b.toMove == "w"
	for idx, val := range b.board {
		if val == 'K' && b.isUpper(idx) && isWhite {
			return b.isInCheck(idx)
		}
		if val == 'k' && !b.isUpper(idx) && !isWhite {
			return b.isInCheck(idx)
		}
	}
	return false
}

// Check if target is in Check
func (b *Board) isInCheck(target int) bool {
	isWhite := b.isUpper(target)
	k := b.board[target]

	// store all the orig of the opponents pieces
	attackers := make([]int, 0, 16)

	for idx, val := range b.board {
		matchWhite, _ := regexp.MatchString(`[PNBRQK]`,
			string(val))
		matchBlack, _ := regexp.MatchString(`[pnbrqk]`,
			string(val))
		if isWhite && matchBlack {
			attackers = append(attackers, idx)
		} else if !isWhite && matchWhite { // black
			attackers = append(attackers, idx)
		}
	}
	//fmt.Println("white ", isWhite, "attackers ", attackers, "king", k)
	// check for valid attacks
	for _, val := range attackers {
		p := string(bytes.ToUpper(b.board[val : val+1]))
		switch {
		case p == "P":
			e := b.validPawn(val, target, k)
			if e == nil {
				fmt.Println("Pawn check")
				return true
			}
		case p == "N":
			e := b.validKnight(val, target)
			if e == nil {
				fmt.Println("Knight check")
				return true
			}
		case p == "B":
			e := b.validBishop(val, target)
			if e == nil {
				fmt.Println("Bishop check")
				return true
			}
		case p == "R":
			e := b.validRook(val, target)
			if e == nil {
				fmt.Println("Rook check")
				return true
			}
		case p == "Q":
			e := b.validQueen(val, target)
			if e == nil {
				return true
			}
		case p == "K":
			e := b.validKing(val, target, false)
			if e == nil {
				return true
			}
		}
	}
	// if nothing was valid, return false
	return false
}

// Check: right-color, origin-empty, attack-enemy
func (b *Board) basicValidation(orig, dest int, o, d byte, isCastle bool) error {
	// Check if it is the right turn
	if b.board[orig] != o {
		return errors.New("Not your turn")
	}
	// Check if Origin is Empty
	if o == '.' {
		return errors.New("Empty square")
	}
	// Check if destination is Enemy
	if b.board[dest] != d && !isCastle { //
		return errors.New("Can't attack your own piece")
	}
	return nil
}

// validate Pawn Move
func (b *Board) validPawn(orig int, dest int, d byte) error {
	err := errors.New("Illegal Pawn Move")
	var remainder int
	var empOffset int
	var empTarget byte
	// Whose turn
	if b.toMove == "w" {
		remainder = dest - orig
		empOffset = -10 // where the empassant piece should be
		empTarget = 'p'
	} else if b.toMove == "b" {
		remainder = orig - dest
		empOffset = 10
		empTarget = 'P'
	}
	// What sort of move
	if remainder == 10 {
		// regular move
		if b.board[dest] != '.' {
			return err
		}
	} else if remainder == 20 { // two spaces
		// double starter move
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		} else if orig < 70 && b.toMove == "b" {
			return err
		}
	} else if remainder == 9 || remainder == 11 {
		// Attack vector
		// check if b.board[dest+10] == '.'
		if b.board[dest] == d && d != '.' {
			// Proper attack
		} else if b.board[dest] == d && dest+empOffset == b.empassant {
			// Empassant attack
			if b.board[dest+empOffset] == empTarget {
				// is the right case
				// TODO move this to UpdateBoard
				b.board[b.empassant] = '.'
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("Not valid Pawn move.")
	}
	return nil
}

// Validate Knight move.
func (b *Board) validKnight(orig int, dest int) error {
	var possibilities [8]int
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+21,
		orig+19, orig+12, orig+8, orig-8,
		orig-12, orig-19, orig-21
	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	return errors.New("Illegal Knight Move")
}

// Validate Bishop move.
func (b *Board) validBishop(orig int, dest int) error {
	// Check if other pieces are in the way
	err := errors.New("Illegal Bishop Move")
	trajectory := orig - dest
	a1h8 := trajectory % 11 // if 0 remainder...
	a8h1 := trajectory % 9
	// Check which slope
	if a1h8 == 0 {
		if dest > orig { // go to bottom right
			for i := orig + 11; i <= dest-11; i += 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if dest < orig { // go to top left
			for i := orig - 11; i >= dest+11; i -= 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if a8h1 == 0 {
		if dest > orig { // go to bottem left
			for i := orig + 9; i <= dest-9; i += 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if orig > dest { // go to top right
			for i := orig - 9; i >= dest+9; i -= 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		return errors.New("Bishop move not valid")
	}
	return nil
}

// Validate rook move.
func (b *Board) validRook(orig int, dest int) error {
	// Check if pieces are in the way
	err := errors.New("Illegal Rook Move")
	remainder := dest - orig
	if remainder < 10 && remainder > -10 {
		// Horizontal
		if remainder < 0 {
			for i := orig - 1; i > dest; i-- {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 1; i < dest; i++ {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		if remainder%10 != 0 {
			return err
		}
		// Vertical
		if remainder < 0 { // descends
			for i := orig - 10; i > dest; i -= 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 10; i < dest; i += 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	}
	return nil
}

// Validate queen move.
func (b *Board) validQueen(orig int, dest int) error {
	err := errors.New("Illegal Queen Move")
	remainder := dest - orig
	vertical := remainder%10 == 0
	horizontal := remainder < 9 && remainder > -9 // Horizontal
	diagA8 := remainder%9 == 0                    // Diag a8h1
	diagA1 := remainder%11 == 0                   // Diag a1h8
	// Check if moves through not-empty squares
	if horizontal { // 1st
		if remainder < 0 {
			for i := orig - 1; i > dest; i-- {
				if b.board[i] != '.' {
					return err
				}
			}
		} else { // go right
			for i := orig + 1; i < dest; i++ {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if vertical {
		if remainder < 0 {
			for i := orig - 10; i > dest; i -= 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 10; i < dest; i += 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if diagA8 {
		if dest > orig { // go to bottem left
			for i := orig + 9; i <= dest-9; i += 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if orig > dest { // go to top right
			for i := orig - 9; i >= dest+9; i -= 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if diagA1 {
		if dest > orig { // go to bottom right
			for i := orig + 11; i <= dest-11; i += 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if dest < orig { // go to top left
			for i := orig - 11; i >= dest+11; i -= 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		return errors.New("Illegal Queen Move")
	}
	// check if anything is inbetween

	return nil
}

// Validate king move.
// Check for castle
func (b *Board) validKing(orig int, dest int, castle bool) error {
	castlerr := errors.New("Something is in your way")
	noCastle := errors.New("Castle on this side is foutu")
	var possibilities [8]int
	g := b.board // g for gameboard
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+10,
		orig+11, orig+1, orig+9, orig-10,
		orig-11, orig-1, orig-9
	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	if castle {
		queenSideCastle := !(g[orig+1] != '.' || g[orig+2] != '.' || g[orig+3] != '.')
		kingSideCastle := !(g[orig-1] != '.' || g[orig-2] != '.')
		if dest > orig { // Queen side
			if !queenSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[1] != 'Q' {
					return noCastle
				}

			} else { // b
				if b.castle[3] != 'q' {
					return noCastle
				}
			}
		} else if orig > dest {
			if !kingSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[0] != 'K' {
					return noCastle
				}
			} else {
				if b.castle[2] != 'k' {
					return noCastle
				}
			}
		}

	} else {
		return errors.New("Illegal King Move")
	}
	return nil
}

// Parse a pgn move
// Infer the origin piece
func (b *Board) ParseMove(move string) error {
	move = strings.TrimRight(move, "\r\n") // prepare for input
	// Variables
	var piece string    // find move piece
	var orig int        // find origin coord of move
	var square string   // find pgnMap key of move
	var attacker string // left of x
	//var precise string // for multiple possibilities
	var target byte // the piece to move, in proper case

	// Status
	isCastle := false
	isWhite := b.toMove == "w"
	isCapture, _ := regexp.MatchString(`x`, move)

	res := b.pattern.FindStringSubmatch(move)
	if res == nil && move != "O-O" && move != "O-O-O" {
		return errors.New("invalid input")
	} else if move == "O-O" || move == "O-O-O" {
		// be nice
		isCastle = true
	}

	// Either is catpure or not
	if isCapture {
		attacker = res[1]
		if attacker == strings.ToLower(attacker) {
			piece = "P"
		} else { // if  upper case, forcement a piece
			piece = res[1]
		}
		square = res[2]
	} else if isCastle {
		if move == "O-O" {
			piece = "K"
			if isWhite {
				square = "h1"
			} else {
				square = "h8"
			}
		} else if move == "O-O-O" {
			piece = "K"
			if isWhite {
				square = "a1"
			} else {
				square = "a8"
			}
		}
	} else { // No x
		chars := len(move)
		if chars == 2 {
			piece = "P"
			square = res[2]
		} else if chars == 3 && move != "0-0" {
			// Breaks when e44 is entered...
			piece = res[1]
			square = res[2] //move[0]
		} else if chars == 4 {
			piece = res[1] // remove second char
			//precise = move
			square = res[2]
		} else {
			return errors.New("Not enough input")
		}
	}
	// the presumed destination
	dest := b.pgnMap[square]
	// The piece will be saved as case sensitive byte
	if b.toMove == "b" {
		target = []byte(strings.ToLower(piece))[0]
	} else {
		target = []byte(piece)[0]
	}
	switch {
	case piece == "P": // Pawn Parse
		var possibilities [2]int // two potentional origins
		// TODO: Allow for empassant take
		if b.toMove == "w" {
			if isCapture {
				possibilities[0],
					possibilities[1] = dest-9,
					dest-11
			} else {
				possibilities[0],
					possibilities[1] = dest-10,
					dest-20
			}
		} else { // is black to move
			if isCapture {
				possibilities[0],
					possibilities[1] = dest+9,
					dest+11
			} else {
				possibilities[0],
					possibilities[1] = dest+10,
					dest+20
			}
		}
		if b.board[possibilities[0]] == target {
			orig = possibilities[0]
		} else if b.board[possibilities[1]] == target {
			orig = possibilities[1]
		}
	case piece == "N": // Knight Parse
		var possibilities [8]int
		// TODO: assume no precision
		// Change to possibilities[]
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+21,
			dest+19, dest+12, dest+8, dest-8,
			dest-12, dest-19, dest-21
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	case piece == "B": // Bishop Parse
		var possibilities [14]int
		ticker := 0
		// a8 - h1
		for i := dest + 9; i < 90; i += 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 9; i > 10; i -= 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// a1 - h8 Vector
		for i := dest + 11; i < 90; i += 11 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 11; i > 10; i -= 11 {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// Find piece origin
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	case piece == "R": // Rook Parse
		var possibilities [14]int
		ticker := 0
		// Horizontal Vector
		for i := dest + 10; i < 90; i += 10 {
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 10; i > 10; i -= 10 {
			possibilities[ticker] = i
			ticker++
		}
		// Vertical Vector
		for i := dest + 1; i < 90; i++ {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 1; i > 10; i-- {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
	Looposs:
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				err := b.validRook(orig, dest)
				if err != nil {
					continue
				}
				break Looposs
			}
		}
	case piece == "Q": // Queen Parse
		for idx, possibility := range b.board {
			if possibility == target {
				orig = idx
				break
			}
		}
	case piece == "K": // King Parse
		var possibilities [8]int
		if isCastle {
			if isWhite {
				orig = 14
			} else {
				orig = 84
			}
			break
		}
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+10,
			dest+11, dest+1, dest+9, dest-10,
			dest-11, dest-1, dest-9
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	}
	// Move the Piece
	// - Validate Move in Board.Move()
	if b.board[dest] != '.' && !isCapture && !isCastle {
		return errors.New("Not the proper capture syntax")

	}
	if orig != 0 && dest != 0 {
		err := b.Move(orig, dest)
		if err == nil {
			// Update pgn History
			if b.toMove == "b" {
				b.pgn += strconv.Itoa(b.moves) + ". "
			}
			b.pgn += (move)
			if b.check {
				// TODO
				// if move doesn't already have check..
				b.pgn += "+ "
				// check for checkmate?
			} else {
				b.pgn += " " // add space
			}
		}
		return err
	} else {
		return errors.New("No such move")
	}
}

// Read a pgn match
func (b *Board) LoadPgn(match string) (Board, error) {
	// TODO: ignore header strings
	game := NewBoard()
	result := game.pattern.FindAllString(match, -1)
	for _, val := range result {
		fmt.Println("Move: ", game.moves)
		err := game.ParseMove(val)
		if err != nil {
			return game, err
		}
	}
	return game, nil
}

func (b *Board) ParseFen() {
	// Parse Fen
}

// Get FEN position
func (b *Board) Position() string {
	// b.board -> Fen
	return b.fen
}

/*
Main thread
*/
func main() {
	board := NewBoard()
	PlayGame(board)
}

// Take user input and commands
// TODO: make a method of board
func PlayGame(board Board) { // TODO Rotate Board
	var turn string
	welcome := `
********
go-chess
    Enter /help for more options

    /~ |_ _  _ _
    \_|||(/__\_\

`
	manuel := `Help:
    Prefix commands with / - slash

Commands:
	quit - exit game
	new - new game
	coordinates - print board coordinates
	pgn - print PGN history
	fen - print FEN position
	set-headers - set PGN headers
	headers - print game info
Tests:
	test-castle
        test-pgn - load a pgn game
`
	reader := bufio.NewReader(os.Stdin)
	// welcome message
	fmt.Println(welcome)
	fmt.Print(board.String())

Loop:
	for {
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Print(turn, " to move: ")
		input, _ := reader.ReadString('\n')
		isCmd, _ := regexp.MatchString(`/`, input)
		if isCmd {
			input = strings.TrimRight(input, "\r\n")
			switch {
			case input == "/help":
				fmt.Print("\n", manuel)
			case input == "/quit":
				break Loop //os.Exit(1)
			case input == "/new":
				board = NewBoard()
				fmt.Print(board.String())
			case input == "/print":
				fmt.Print(board.String())
			case input == "/coordinates":
				fmt.Println("Coordinates:")
				board.Coordinates()
			case input == "/pgn":
				fmt.Println("PGN history:")
				fmt.Println(board.headers,
					board.pgn, "\n")
			case input == "/load-pgn":
				var err error
				fmt.Print("Enter PGN history: ")
				history, _ := reader.ReadString('\n')
				board, err = board.LoadPgn(history)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(board.String())
			case input == "/fen":
				fmt.Println("FEN position:")
			case input == "/set-headers":
				fmt.Println("Set Headers:")
				fmt.Print("White: ")
				inWhite, _ := reader.ReadString('\n')
				fmt.Print("Black: ")
				inBlack, _ := reader.ReadString('\n')
				board.setHeaders(inWhite, inBlack)
			case input == "/headers":
				fmt.Println(board.headers)
			case input == "/test-pgn":
				hist := `1. Nf3 Nc6 2. d4 d5 3. c4 e6 4. e3 Nf6 5. Nc3 Be7 6. a3 O-O 7. b4 a6 8. Be2 Re8 9. O-O Bf8 10. c5 g6 11. b5 axb5 12. Bxb5 Bd7 13. h3 Na5 14. Bd3 Nc6 15. Rb1 Qc8 16. Nb5 e5 17. Be2 e4 18. Ne1 h6 19. Nc2 g5 20. f3 exf3 21. Bxf3 g4 22. hxg4 Bxg4 23. Nxc7 Qxc7 24. Bxg4 Nxg4 25. Qxg4+ Bg7 26. Nb4 Nxb4 27. Rxb4 Ra6 28. Rf5 Re4 29. Qh5 Rg6 30. Qh3 Qc8 31. Qf3 Qd7 32. Rb2 Bxd4 33. exd4 Re1+ 34. Kh2 Rxc1 35. Qxd5 Qe7 36. g3 Qc7 37. Rf4 b6 38. a4 Rg5 39. cxb6 Rxd5 40. bxc7 Rxc7 41. Rb5 Rc2+`
				var err error
				board, err = board.LoadPgn(hist)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(board.String())
				if board.check {
					fmt.Println("****Check!****")
				}
			case input == "/test-castle-check":
				hist := `1. e4 e5 2. d4 d5 3. Bh6 Bh3 4. Nxh3 Nxh6 5. Qg4 Qg5 6. Na3 Na6 7. Qf5 Qg4 8. Qg5 Bc5 9. Bc4 dxe4 10. dxe5 f6 11. exf6 gxf6 12. Qxf6 Bb4+ 13. c3 e3`
				var err error
				board, err = board.LoadPgn(hist)
				if err != nil {
					fmt.Println(err)
				}
			case input == "/test-check":
				hist := `1. e4 e5 2. Qf3 Qg5 3. Qxf7`
				var err error
				board, err = board.LoadPgn(hist)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(board.String())
			default:
				fmt.Println("Mysterious input")
			}
			continue
		}
		e := board.ParseMove(input)
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Println("\n-------------------")
		panel := "Debug Mode:\nMove: " + strconv.Itoa(board.moves) + " | Castle: " + string(board.castle) + "\nCheck: " + strconv.FormatBool(board.check) + " | Turn: " + string(turn)
		// TODO use formats.
		fmt.Println(panel)
		if e != nil {
			fmt.Printf("   [Error: %v]\n", e)
		}
		fmt.Print(board.String())
		if board.check {
			fmt.Println("****Check!****")
		}
	}
	fmt.Println("\nGood Game.")
}

// Set pgnHeaders for a pgn export
func (b *Board) setHeaders(w, bl string) {
	w = strings.TrimRight(w, "\r\n")
	bl = strings.TrimRight(bl, "\r\n")
	y, m, d := time.Now().Date()
	ye, mo, da := strconv.Itoa(y), strconv.Itoa(int(m)),
		strconv.Itoa(d)
	white := "[White \"" + w + "\"]"
	black := "[Black \"" + bl + "\"]"
	date := "[Date \"" + ye + "." + mo + "." + da + "\"]"
	result := `[Result "*"]`
	b.headers = white + "\n" + black + "\n" + date + "\n" + result + "\n"
}

// Print the Board.Move() coordinate
func (b *Board) Coordinates() {
	// TODO Rotate Board
	game := b.board
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i]) + ": "
			continue
		}
		printBoard += "|" + strconv.Itoa(i) + "|"
	}
	printBoard += "\n"
	printBoard += "   :a ::b ::c ::d ::e ::f ::g ::h :\n"
	fmt.Println(printBoard)
}

// Check if byte in board is upper case.
func (b Board) isUpper(x int) bool {
	//compare = []byte(bytes.ToLower(b))[0]
	compare := byte(unicode.ToUpper(rune(b.board[x])))
	if b.board[x] == compare {
		return true
	} else {
		return false
	}
}
