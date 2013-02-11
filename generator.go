package main

import (
	"fmt"
)

type Rule struct {
	Lhs   string
	Rhs   Rhs
	Code  string
	Count int
}

type Rhs struct {
	Names []string
	Min   int
}

func (r Rule) Action() string {
	return fmt.Sprintf("_action_%s_%d", r.Lhs, r.Count)
}

func (r Rule) String() string {
	return fmt.Sprintf("%s ::= %s\n\t%s", r.Lhs, r.Rhs, r.Code)
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

func (rhs Rhs) Sequence() bool {
	return rhs.Min >= 0
}
