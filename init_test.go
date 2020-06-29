package togo_test

import lhm "github.com/emirpasic/gods/maps/linkedhashmap"

var namesLevel *lhm.Map
var names []string

func init() {
	names = []string{
		"foo", "bar", "baz", "qux", "quux", "quuz", "corge", "grault",
		"garply", "waldo", "fred", "plugh", "xyzzy", "thud",
		"Wibble", "wobble", "wubble", "flob",
	}
	namesLevel = lhm.New()
	for i := 0; i < 9; i++ {
		lvl := i / 3
		namesLevel.Put(names[i], lvl)
	}
	lhm.New()
}
