package main

const (
	maxEval  = +10000
	minEval  = -maxEval
	mateEval = maxEval + 1
	noScore  = minEval - 1
)

var pieceVal = [16]int{100, -100, 300, -300, 350, -350, 500, -500, 1000, -1000, 10000, -10000, 0, 0, 0, 0}

var knightFile = [8]int{-4, -3, -2, +2, +2, 0, -2, -4}
var knightRank = [8]int{-15, 0, +5, +6, +7, +8, +2, -4}
var centerFile = [8]int{-8, -1, 0, +1, +1, 0, -1, -3}
var kingFile = [8]int{+1, +2, 0, -2, -2, 0, +2, +1}
var kingRank = [8]int{+1, 0, -2, -4, -6, -8, -10, -12}
var pawnRank = [8]int{0, 0, 0, 0, +2, +6, +25, 0}
var pawnFile = [8]int{0, 0, +1, +10, +10, +8, +10, +8}

const longDiag = 10

// Piece Square Table
var pSqTab [12][64]int

//TODO: eval hash
//TODO: pawn hash
//TODO: pawn structures. isolated, backward, duo, passed (guarded and not), double and more...
//TODO: bishop pair
//TODO: King safety. pawn shelter, guarding pieces
//TODO: King attack. Attacking area surrounding the enemy king, closeness to the enemy king
//TODO: space, center control, knight outposts, connected rooks, 7th row and more
//TODO: combine middle game and end game values

// evaluate returns score from white pov
func evaluate(b *boardStruct) int {
	ev := 0
	for sq := A1; sq <= H8; sq++ {
		p12 := b.sq[sq]
		if p12 == empty {
			continue
		}
		ev += pieceVal[p12]
		ev += pcSqScore(p12, sq)
	}
	return ev
}

// Score returns the pc2pt square table value for a given pc2pt on a given square. Stage = MG/EG
func pcSqScore(p12, sq int) int {
	return pSqTab[p12][sq]
}

// PstInit inits the pieces-square-tables when the program starts
func pcSqInit() {
	tell("info string pStInit starter")
	// for p12 := 0; p12 < 12; p12++ {
	// 	for sq := 0; sq < 64; sq++ {
	// 		pSqTab[p12][sq] = 0
	// 	}
	// }

	// for sq := 0; sq < 64; sq++ {

	// 	fl := sq % 8
	// 	rk := sq / 8

	// 	pSqTab[wP][sq] = pawnFile[fl] + pawnRank[rk]

	// 	pSqTab[wN][sq] = knightFile[fl] + knightRank[rk]
	// 	pSqTab[wB][sq] = centerFile[fl] + centerFile[rk]*2

	// 	pSqTab[wR][sq] = centerFile[fl] * 5

	// 	pSqTab[wQ][sq] = centerFile[fl] + centerFile[rk]

	// 	pSqTab[wK][sq] = (kingFile[fl] + kingRank[rk]) * 8
	// }

	// // bonus for e4 d5 and c4
	// pSqTab[wP][E2], pSqTab[wP][D2], pSqTab[wP][E3], pSqTab[wP][D3], pSqTab[wP][E4], pSqTab[wP][D4], pSqTab[wP][C4] = 0, 0, 6, 6, 24, 20, 12

	// // long diagonal
	// for sq := A1; sq <= H8; sq += NE {
	// 	pSqTab[wB][sq] += longDiag - 2
	// }
	// for sq := H1; sq <= A8; sq += NW {
	// 	pSqTab[wB][sq] += longDiag
	// }

	pawnScore := [64]int{
		90, 90, 90, 90, 90, 90, 90, 90,
		30, 30, 30, 40, 40, 30, 30, 30,
		20, 20, 20, 30, 30, 30, 20, 20,
		10, 10, 10, 20, 20, 10, 10, 10,
		5, 5, 10, 20, 20, 5, 5, 5,
		0, 0, 0, 5, 5, 0, 0, 0,
		0, 0, 0, -10, -10, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	knightScore := [64]int{
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 10, 10, 0, 0, -5,
		-5, 5, 20, 20, 20, 20, 5, -5,
		-5, 10, 20, 30, 30, 20, 10, -5,
		-5, 10, 20, 30, 30, 20, 10, -5,
		-5, 5, 20, 10, 10, 20, 5, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, -10, 0, 0, 0, 0, -10, -5,
	}
	bishopScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 10, 10, 0, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 10, 0, 0, 0, 0, 10, 0,
		0, 30, 0, 0, 0, 0, 30, 0,
		0, 0, -10, 0, 0, -10, 0, 0,
	}
	rookScore := [64]int{
		50, 50, 50, 50, 50, 50, 50, 50,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
	}
	kingScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 5, 5, 5, 0, 0,
		0, 5, 5, 10, 10, 5, 5, 0,
		0, 5, 10, 20, 20, 10, 5, 0,
		0, 5, 10, 20, 20, 10, 5, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 5, 5, -5, -5, 0, 5, 0,
		0, 0, 5, 0, -15, 0, 10, 0,
	}
	queenScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	pSqTab[wP] = pawnScore
	pSqTab[wN] = knightScore
	pSqTab[wR] = rookScore
	pSqTab[wK] = kingScore
	pSqTab[wB] = bishopScore
	pSqTab[wQ] = queenScore

	// for Black
	for pc := Pawn; pc <= King; pc++ {

		wP12 := pt2pc(pc, WHITE)
		bP12 := pt2pc(pc, BLACK)

		for bSq := 0; bSq < 64; bSq++ {
			wSq := oppRank(bSq)
			pSqTab[bP12][bSq] = -pSqTab[wP12][wSq]
		}
	}
}

// mirror the rank_sq
func oppRank(sq int) int {
	fl := sq % 8
	rk := sq / 8
	rk = 7 - rk
	return rk*8 + fl
}
