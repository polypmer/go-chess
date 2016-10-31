package ghess

var (
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
		21: 0,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: 0,
		31: 0,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: 0,
		41: 0,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: 0,
		51: 0,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: 0,
		61: 0,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: 0,
		71: 0,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: 0,
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
		21: 0,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: 0,
		31: 0,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: 0,
		41: 0,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: 0,
		51: 0,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: 0,
		61: 0,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: 0,
		71: 0,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: 0,
		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}

	bishopMap = map[int]int{
		11: -20,
		12: -10,
		13: -10,
		14: -10,
		15: -10,
		16: -10,
		17: -10,
		18: -10,
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
	knightMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,
		21: 0,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: 0,
		31: 0,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: 0,
		41: 0,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: 0,
		51: 0,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: 0,
		61: 0,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: 0,
		71: 0,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: 0,
		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}
	rookMap = map[int]int{
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,
		21: 0,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: 0,
		31: 0,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: 0,
		41: 0,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: 0,
		51: 0,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: 0,
		61: 0,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
		67: 0,
		68: 0,
		71: 0,
		72: 0,
		73: 0,
		74: 0,
		75: 0,
		76: 0,
		77: 0,
		78: 0,
		81: 0,
		82: 0,
		83: 0,
		84: 0,
		85: 0,
		86: 0,
		87: 0,
		88: 0,
	}
)