package main

import (
	"os"
	"text/template"
)

var (
	codeTemplate *template.Template = template.Must(template.ParseFiles("code.templ"))
)

func main() {
	p := NewParser()

	rules := p.Parse(os.Stdin)

	for i, _ := range rules {
		rules[i].Count = i
	}

	err := codeTemplate.Execute(os.Stdout, rules)

	if err != nil {
		panic(err)
	}
}
