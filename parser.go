package main

import (
	marpa "github.com/pstuifzand/go-marpa"
	"io"
	"io/ioutil"
)

type MarpaParser struct {
	grammar *marpa.Grammar
	re      *marpa.Recognizer
}

func NewParser() *MarpaParser {
	g := marpa.NewGrammar()

	g.StartRule("rules")

	g.AddSequence("rules", "rule", marpa.Seq{Min: 1}, ActionRules)
	g.AddRule("rule", []string{"lhs", "bnfop", "rhs"}, ActionRule)
	g.AddRule("rule", []string{"lhs", "bnfop", "rhs", "code"}, ActionRule)
	g.AddRule("lhs", []string{"name"}, ActionLhs)
	g.AddRule("rhs", []string{"names"}, ActionRhs)
	g.AddRule("rhs", []string{"name", "plus"}, ActionPlus)
	g.AddRule("rhs", []string{"name", "star"}, ActionStar)
	g.AddRule("code", []string{"leftcode", "innercode", "rightcode"}, ActionCode)
	g.AddSequence("names", "name", marpa.Seq{Min: 1}, ActionNames)

	g.Precompute()

	re, err := marpa.NewRecognizer(g)
	if err != nil {
		return nil
	}

	p := &MarpaParser{g, re}
	return p
}

func (p *MarpaParser) read(token, value string) {
	p.re.Read(token, value)
}

func (p *MarpaParser) Parse(r io.Reader) []Rule {
	// call read
	input, _ := ioutil.ReadAll(r)
	_, items := NewScanner(string(input))

TOKENS:
	for {
		item := <-items
		switch item.typ {
		case itemNone, itemEOF:
			break TOKENS
		default:
			sym := tokenNames[item.typ]
			p.read(sym, item.val)
		}
	}

	val := p.re.Value()
	if val.Next() {
		return val.Value().([]Rule)
	}
	return nil
}
