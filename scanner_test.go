package main

import (
	"testing"
)

func IsItem(t *testing.T, itm item, typ itemType, token string) {
	if itm.typ != typ {
		t.Fatalf("Item type is wrong: expected \"%s\", but is \"%s\" (%d != %d)", tokenNames[typ], tokenNames[itm.typ], typ, itm.typ)
	}
	if itm.val != token {
		t.Fatalf("Item is wrong: is %s, but expected %s (%d != %d)", itm.val, token)
	}
}

func TestScanner(t *testing.T) {
	_, items := NewScanner("test")

	var x item

	x = <-items
	IsItem(t, x, itemName, "test")

	x = <-items
	IsItem(t, x, itemEOF, "")

}

func TestScanner2(t *testing.T) {
	_, items := NewScanner("test ::= test+")

	var x item

	x = <-items
	IsItem(t, x, itemName, "test")

	x = <-items
	IsItem(t, x, itemBnfOp, "::=")

	x = <-items
	IsItem(t, x, itemName, "test")

	x = <-items
	IsItem(t, x, itemPlus, "+")

	x = <-items
	IsItem(t, x, itemEOF, "")
}
