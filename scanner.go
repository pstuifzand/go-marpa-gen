package main

import (
	//"io"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eof = -1
)

const (
	itemNone itemType = iota
	itemEOF
	itemName
	itemBnfOp
	itemPlus
	itemStar
	itemLeftCode
	itemRightCode
	itemInnerCode
)

var (
	tokenNames = []string{
		"itemNone",
		"itemEOF",
		"name",
		"bnfop",
		"plus",
		"star",
		"leftcode",
		"rightcode",
		"innercode",
	}
)

type stateFn func(*scanner) stateFn

type itemType int

type item struct {
	typ itemType
	val string
}

type scanner struct {
	input string
	//item  item
	items chan item
	start int
	pos   int
	width int
}

func lexExpr(s *scanner) stateFn {
	for {
		if strings.HasPrefix(s.input[s.pos:], "::=") {
			return lexBnfOp
		} else if strings.HasPrefix(s.input[s.pos:], "+") {
			return lexPlus
		} else if strings.HasPrefix(s.input[s.pos:], "*") {
			return lexStar
		} else if strings.HasPrefix(s.input[s.pos:], "{{") {
			return lexLeftCode
		}

		r := s.next()
		if r == eof {
			break
		} else if unicode.IsSpace(r) {
			s.ignore()
		} else if unicode.IsLetter(r) {
			s.backup()
			return lexName
		} else {
			fmt.Printf("Unknown\n")
			break
		}
	}
	s.emit(itemEOF)
	return nil
}

func (s *scanner) ignore() {
	s.start = s.pos
}

func (s *scanner) backup() {
	s.pos -= s.width
}

func lexBnfOp(s *scanner) stateFn {
	s.pos += len("::=")
	s.emit(itemBnfOp)
	return lexExpr
}

func lexPlus(s *scanner) stateFn {
	s.pos += len("+")
	s.emit(itemPlus)
	return lexExpr
}

func lexStar(s *scanner) stateFn {
	s.pos += len("*")
	s.emit(itemStar)
	return lexExpr
}

func lexInnerCode(s *scanner) stateFn {
	for {
		if strings.HasPrefix(s.input[s.pos:], "}}") {
			s.emit(itemInnerCode)
			return lexRightCode
		}

		r := s.next()
		if r == eof {
			fmt.Printf("Unterminated code\n")
			break
		}
	}
	s.emit(itemEOF)
	return nil
}

func lexLeftCode(s *scanner) stateFn {
	s.pos += len("{{")
	s.emit(itemLeftCode)
	return lexInnerCode
}

func lexRightCode(s *scanner) stateFn {
	s.pos += len("}}")
	s.emit(itemRightCode)
	return lexExpr
}

func lexName(s *scanner) stateFn {
	for {
		r := s.next()
		if r == eof {
			s.backup()
			s.emit(itemName)
			s.emit(itemEOF)
			return nil
		} else if !unicode.IsLetter(r) {
			s.backup()
			s.emit(itemName)
			return lexExpr
		}
	}
	return lexExpr
}

func NewScanner(r string) (*scanner, chan item) {
	items := make(chan item)
	scanner := &scanner{input: r, items: items}
	go scanner.run()
	return scanner, scanner.items
}

func (s *scanner) run() {
	for state := lexExpr; state != nil; {
		state = state(s)
	}
	close(s.items)
}

func (s *scanner) emit(t itemType) {
	s.items <- item{t, s.input[s.start:s.pos]}
	s.start = s.pos
}

func (s *scanner) next() (r rune) {
	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	return r
}
