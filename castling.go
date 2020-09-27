package main

import (
	"strings"
)

type castling uint

const (
	shortW = uint(0x1) // white can casle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can casle short
	longB  = uint(0x8) // black can castle long
)

func (c *castling) on(val uint) {
	(*c) |= castling(val)
}

func (c *castling) off(val uint) {
	(*c) &= castling(^val)
}

func (c castling) String() string {
	flags := ""

	if uint(c)&shortW != 0 {
		flags += "K"
	}
	if uint(c)&longB != 0 {
		flags += "Q"
	}
	if uint(c)&shortB != 0 {
		flags += "k"
	}
	if uint(c)&longB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}

// parse castling rights in fenstring
func parseCastling(fenCastl string) castling {
	c := uint(0)
	if fenCastl == "-" {
		return castling(0)
	}

	if strings.Index(fenCastl, "K") >= 0 {
		c |= shortW
	}

	if strings.Index(fenCastl, "Q") >= 0 {
		c |= longW
	}

	if strings.Index(fenCastl, "k") >= 0 {
		c |= shortB
	}

	if strings.Index(fenCastl, "q") >= 0 {
		c |= longB
	}

	return castling(c)
}