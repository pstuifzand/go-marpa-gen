package main

import (
	"fmt"
	marpa "github.com/pstuifzand/go-marpa"
	"io"
	"io/ioutil"
	"os"
)

type MarpaParser struct {
	grammar *marpa.Grammar
	re      *marpa.Recognizer
}

type Rule struct {
	Lhs string
	Rhs Rhs
}

func (r Rule) String() string {
	return fmt.Sprintf("%s ::= %s", r.Lhs, r.Rhs)
}

func (r Rhs) String() string {
	min_str := ""
	if r.Min == 0 {
		min_str = "*"
	} else if r.Min == 1 {
		min_str = "+"
	}

	s := r.Names[0]
	for _, name := range r.Names[1:] {
		s += " " + name
	}
	return fmt.Sprintf("%s%s", s, min_str)
}

type Rhs struct {
	Names []string
	Min   int
}

func ActionRules(args []interface{}) interface{} {
	rules := []Rule{}
	for _, n := range args {
		rule := n.(Rule)
		rules = append(rules, rule)
	}
	return rules
}
func ActionRule(args []interface{}) interface{} {
	lhs := args[0].(string)
	rhs := args[2].(Rhs)
	return Rule{Lhs: lhs, Rhs: rhs}
}
func ActionPlus(args []interface{}) interface{} {
	return Rhs{Names: []string{args[0].(string)}, Min: 1}
}
func ActionStar(args []interface{}) interface{} {
	return Rhs{Names: []string{args[0].(string)}, Min: 0}
}
func ActionLhs(args []interface{}) interface{} {
	return args[0]
}
func ActionRhs(args []interface{}) interface{} {
	return Rhs{Names: args[0].([]string), Min: -1}
}
func ActionNames(args []interface{}) interface{} {
	names := []string{}
	for _, n := range args {
		name := n.(string)
		names = append(names, name)
	}
	return names
}

func NewParser() *MarpaParser {
	g := marpa.NewGrammar()

	g.StartRule("rules")
	g.AddSequence("rules", "rule", marpa.Seq{Min: 1}, ActionRules)
	g.AddRule("rule", []string{"lhs", "bnfop", "rhs"}, ActionRule)
	g.AddRule("lhs", []string{"name"}, ActionLhs)
	g.AddRule("rhs", []string{"names"}, ActionRhs)
	g.AddRule("rhs", []string{"name", "plus"}, ActionPlus)
	g.AddRule("rhs", []string{"name", "star"}, ActionStar)
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

func main() {
	p := NewParser()
	rules := p.Parse(os.Stdin)
	for _, rule := range rules {
		fmt.Printf("%s\n", rule)
	}
}
