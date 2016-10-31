package ghess

import "regexp"

var (

	/* ***************************************************
	   Scoring Tables
	   ***************************************************  */
	matMap = map[byte]int{
		'P': 100,
		'N': 320,
		'B': 330,
		'R': 500,
		'Q': 900,
		'K': 20000,
	} // material map

	blackPawnMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,

		71: 5,
		72: 10,
		73: 10,
		74: -20,
		75: -20,
		76: 10,
		77: 10,
		78: 5,

		61: 5,
		62: -5,
		63: -10,
		64: 0,
		65: 0,
		66: -10,
		67: -5,
		68: 5,

		51: 0,
		52: 0,
		53: 0,
		54: 20,
		55: 20,
		56: 0,
		57: 0,
		58: 0,

		41: 5,
		42: 5,
		43: 10,
		44: 25,
		45: 25,
		46: 10,
		47: 5,
		48: 5,

		31: 10,
		32: 10,
		33: 20,
		34: 30,
		35: 30,
		36: 20,
		37: 10,
		38: 10,

		21: 50,
		22: 50,
		23: 50,
		24: 50,
		25: 50,
		26: 50,
		27: 50,
		28: 50,

		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}
	whitePawnMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,

		21: 5,
		22: 10,
		23: 10,
		24: -20,
		25: -20,
		26: 10,
		27: 10,
		28: 5,

		31: 5,
		32: -5,
		33: -10,
		34: 0,
		35: 0,
		36: -10,
		37: -5,
		38: 5,

		41: 0,
		42: 0,
		43: 0,
		44: 20,
		45: 20,
		46: 0,
		47: 0,
		48: 0,

		51: 5,
		52: 5,
		53: 10,
		54: 25,
		55: 25,
		56: 10,
		57: 5,
		58: 5,

		61: 10,
		62: 10,
		63: 20,
		64: 30,
		65: 30,
		66: 20,
		67: 10,
		68: 10,

		71: 50,
		72: 50,
		73: 50,
		74: 50,
		75: 50,
		76: 50,
		77: 50,
		78: 50,

		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}

	blackBishopMap = map[int]int{
		81: -20,
		82: -10,
		83: -10,
		84: -10,
		85: -10,
		86: -10,
		87: -10,
		88: -20,

		71: -10,
		72: 5,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 5,
		78: -10,

		61: -10,
		62: 10,
		63: 10,
		64: 10,
		65: 10,
		66: 10,
		67: 10,
		68: -10,

		51: -10,
		52: 0,
		53: 10,
		54: 10,
		55: 10,
		56: 10,
		57: 0,
		58: -10,

		41: -10,
		42: 5,
		43: 5,
		44: 10,
		45: 10,
		46: 5,
		47: 5,
		48: -10,

		31: -10,
		32: 0,
		33: 5,
		34: 10,
		35: 10,
		36: 5,
		37: 0,
		38: -10,

		21: -20,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: -20,

		11: -20,
		12: -10,
		13: -10,
		14: -10,
		15: -10,
		16: -10,
		17: -10,
		18: -20,
	}
	whiteBishopMap = map[int]int{
		11: -20,
		12: -10,
		13: -10,
		14: -10,
		15: -10,
		16: -10,
		17: -10,
		18: -20,
		21: -10,
		22: 5,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 5,
		28: -10,
		31: -10,
		32: 10,
		33: 10,
		34: 10,
		35: 10,
		36: 10,
		37: 10,
		38: -10,
		41: -10,
		42: 0,
		43: 10,
		44: 10,
		45: 10,
		46: 10,
		47: 0,
		48: -10,
		51: -10,
		52: 5,
		53: 5,
		54: 10,
		55: 10,
		56: 5,
		57: 5,
		58: -10,
		61: -10,
		62: 0,
		63: 5,
		64: 10,
		65: 10,
		66: 5,
		67: 0,
		68: -10,
		71: -20,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: -20,
		81: -20,
		82: -10,
		83: -10,
		84: -10,
		85: -10,
		86: -10,
		87: -10,
		88: -20,
	}
	whiteKnightMap = map[int]int{
		11: -50,
		12: -40,
		13: -30,
		14: -30,
		15: -30,
		16: -30,
		17: -40,
		18: -50,
		21: -40,
		22: -20,
		23: 0,
		24: 5,
		25: 5,
		26: 0,
		27: -20,
		28: -40,
		31: -30,
		32: 5,
		33: 10,
		34: 15,
		35: 15,
		36: 10,
		37: 5,
		38: -30,
		41: -30,
		42: 0,
		43: 15,
		44: 20,
		45: 20,
		46: 15,
		47: 0,
		48: -30,
		51: -30,
		52: 5,
		53: 15,
		54: 20,
		55: 20,
		56: 15,
		57: 5,
		58: -30,
		61: -30,
		62: 0,
		63: 10,
		64: 15,
		65: 15,
		66: 10,
		67: 0,
		68: -30,
		71: -40,
		72: -20,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: -20,
		78: -40,
		81: -50,
		82: -40,
		83: -30,
		84: -30,
		85: -30,
		86: -30,
		87: -40,
		88: -50,
	}

	blackKnightMap = map[int]int{
		11: -50,
		12: -40,
		13: -30,
		14: -30,
		15: -30,
		16: -30,
		17: -40,
		18: -50,

		21: -40,
		22: -20,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: -20,
		28: -40,

		31: -30,
		32: 0,
		33: 10,
		34: 15,
		35: 15,
		36: 10,
		37: 0,
		38: -30,

		41: -30,
		42: 5,
		43: 15,
		44: 20,
		45: 20,
		46: 15,
		47: 5,
		48: -30,

		51: -30,
		52: 0,
		53: 15,
		54: 20,
		55: 20,
		56: 15,
		57: 0,
		58: -30,

		61: -30,
		62: 5,
		63: 10,
		64: 15,
		65: 15,
		66: 10,
		67: 5,
		68: -30,

		71: -40,
		72: -20,
		73: 0,
		74: 5,
		75: 5,
		76: 0,
		77: -20,
		78: -40,

		81: -50,
		82: -40,
		83: -30,
		84: -30,
		85: -30,
		86: -30,
		87: -40,
		88: -50,
	}

	whiteRookMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 5,
		15: 5,
		16: 0,
		17: 0,
		18: 0,

		21: -5,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: -5,

		31: -5,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: -5,

		41: -5,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: -5,

		51: -5,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: -5,

		61: -5,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: -5,

		71: 5,
		72: 10,
		73: 10,
		74: 10,
		75: 10,
		76: 10,
		77: 10,
		78: 5,

		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}
	blackRookMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,

		21: 5,
		22: 10,
		23: 10,
		24: 10,
		25: 10,
		26: 10,
		27: 10,
		28: 5,

		31: -5,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: -5,

		41: -5,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: -5,

		51: -5,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: -5,

		61: -5,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: -5,

		71: -5,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: -5,

		81: 0,
		82: 0,
		83: 0,
		84: 5,
		85: 5,
		86: 0,
		87: 0,
		88: 0,
	}
	// TODO: Queen and King tables, King special for end game

	/* ***************************************************
	   Opening Dictionaries
	   ***************************************************  */
	dict = map[string][2]int{
		// e4 e5 Bc4
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w":     [2]int{24, 44},
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b":   [2]int{74, 54},
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w": [2]int{13, 46},
		//2. Bc4 Nf6
		"rnbqkbnr/pppp1ppp/8/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR b": [2]int{82, 63},
		// 3. Nc3
		"rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR w": [2]int{17, 36},
		// 3. Nc3 Nc6
		"rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/2N5/PPPP1PPP/R1BQK1NR b": [2]int{87, 66},

		// alternative white move 2. Nf3
		// 2 Nf3 Nc6
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b": [2]int{87, 66},
		// 3 Bc4
		"r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R b": [2]int{83, 56},

		/* d4 d5 */
		//1 d4 d5
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR b": [2]int{75, 55},
		//2 Nf3
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR w": [2]int{12, 33},

		/*  Nf3 */
		//1 Nf3 d5
		"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b": [2]int{75, 55},
		// 2. d4
		"rnbqkbnr/ppp1pppp/8/3p4/8/5N2/PPPPPPPP/RNBQKB1R w": [2]int{25, 45},
		/* Sicilian


		1   e4 c5
		2.   Nf3 d6
		3.   d4 cxd4
		4.   Nxd4 Nf6
		5.   Nc3 g6
		6.   Be3 Bg7
		7.   f3 O-O
		*/
		// after c5 -> Nf3
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w": [2]int{12, 33},
		// 2. Nf3 -> d6
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b": [2]int{75, 65},
		// 3. -> d4
		"rnbqkbnr/pp2pppp/3p4/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R w": [2]int{25, 45},
		// 3. d4 -> cxd4
		"rnbqkbnr/pp2pppp/3p4/2p5/3PP3/5N2/PPP2PPP/RNBQKB1R b": [2]int{56, 45},
		// 4. -> Nxd4
		"rnbqkbnr/pp2pppp/3p4/8/3pP3/5N2/PPP2PPP/RNBQKB1R w": [2]int{33, 45},
		// 4. Nxd4 -> g6
		"rnbqkbnr/pp2pppp/3p4/8/3NP3/8/PPP2PPP/RNBQKB1R b": [2]int{72, 62},
	}

	/* ***************************************************
	   Board Maps
	   ***************************************************  */
	// Regex patterns for parsing
	PgnPattern = /* const */ regexp.MustCompile(`([PNBRQK]?[a-h]?[1-8]?)x?([a-h][1-8])([\+\?\!]?)|O(-?O){1,2}`)
	FenPattern = /* const */ regexp.MustCompile(`([PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8})\s(w|b)\s([KQkq-]{1,4})\s([a-h][36]|-)\s\d\s([1-9]?[1-9])`)
	// TODO: Enter the map values in NewBoard here
	PgnRowMap = map[int][8]int{
		1: [8]int{18, 17, 16, 15, 14, 13, 12, 11},
		2: [8]int{28, 27, 26, 25, 24, 23, 22, 21},
		3: [8]int{38, 37, 36, 35, 34, 33, 32, 31},
		4: [8]int{48, 47, 46, 45, 44, 43, 42, 41},
		5: [8]int{58, 57, 56, 55, 54, 53, 52, 51},
		6: [8]int{68, 67, 66, 65, 64, 63, 62, 61},
		7: [8]int{78, 77, 76, 75, 74, 73, 72, 71},
		8: [8]int{88, 87, 86, 85, 84, 83, 82, 81},
	}

	UnicodeMap = map[string]string{
		"p": "\u2659",
		"P": "\u265F",
		"b": "\u2657",
		"B": "\u265D",
		"n": "\u2658",
		"N": "\u265E",
		"r": "\u2656",
		"R": "\u265C",
		"q": "\u2655",
		"Q": "\u265B",
		"k": "\u2654",
		"K": "\u265A",
		".": "\u00B7",
	}

	PgnToCoordMap = map[string]int{
		"a1": 18,
		"b1": 17,
		"c1": 16,
		"d1": 15,
		"e1": 14,
		"f1": 13,
		"g1": 12,
		"h1": 11,
		"a2": 28,
		"b2": 27,
		"c2": 26,
		"d2": 25,
		"e2": 24,
		"f2": 23,
		"g2": 22,
		"h2": 21,
		"a3": 38,
		"b3": 37,
		"c3": 36,
		"d3": 35,
		"e3": 34,
		"f3": 33,
		"g3": 32,
		"h3": 31,
		"a4": 48,
		"b4": 47,
		"c4": 46,
		"d4": 45,
		"e4": 44,
		"f4": 43,
		"g4": 42,
		"h4": 41,
		"a5": 58,
		"b5": 57,
		"c5": 56,
		"d5": 55,
		"e5": 54,
		"f5": 53,
		"g5": 52,
		"h5": 51,
		"a6": 68,
		"b6": 67,
		"c6": 66,
		"d6": 65,
		"e6": 64,
		"f6": 63,
		"g6": 62,
		"h6": 61,
		"a7": 78,
		"b7": 77,
		"c7": 76,
		"d7": 75,
		"e7": 74,
		"f7": 73,
		"g7": 72,
		"h7": 71,
		"a8": 88,
		"b8": 87,
		"c8": 86,
		"d8": 85,
		"e8": 84,
		"f8": 83,
		"g8": 82,
		"h8": 81,
	}

	//CoordToPgnMap = map[int]string

)
