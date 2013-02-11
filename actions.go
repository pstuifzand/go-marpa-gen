package main

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

	var code string

	if len(args) == 4 {
		if c, ok := args[3].(string); ok {
			code = c
		}
	}
	return Rule{Lhs: lhs, Rhs: rhs, Code: code}
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
func ActionCode(args []interface{}) interface{} {
	return args[1]
}
func ActionNames(args []interface{}) interface{} {
	names := []string{}
	for _, n := range args {
		name := n.(string)
		names = append(names, name)
	}
	return names
}
