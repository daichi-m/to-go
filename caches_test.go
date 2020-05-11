package togo

import (
	"log"
	"testing"
)

func createLevelCache() levelCache {
	internalCache := make(map[int][]string)
	internalCache[0] = []string{"Foo", "Bar", "Baz"}
	internalCache[1] = []string{"Hello", "World"}
	internalCache[2] = []string{"Few", "More", "Strings"}

	l := levelCache{
		internalCache: internalCache,
		maxLevel:      2,
	}
	return l
}

func TestSearchInLevel(t *testing.T) {

	l := createLevelCache()
	table := []struct {
		level  int
		name   string
		result int
	}{
		{0, "Foo", 0},
		{1, "World", 1},
		{2, "Few", 0},
		{3, "Example", -1},
		{0, "Few", -1},
	}

	for _, tb := range table {
		res := l.searchInLevel(tb.level, tb.name)
		if res != tb.result {
			t.Errorf("Expected %v but found %v \n", tb.result, res)
		}
	}
}

func TestFindLevel(t *testing.T) {

	l := createLevelCache()
	table := []struct {
		name  string
		level int
	}{
		{"Foo", 0},
		{"World", 1},
		{"Few", 2},
		{"Random", -1},
	}

	for _, tb := range table {
		res := l.findLevel(tb.name)
		if res != tb.level {
			t.Errorf("Expected %v but found %v\n", tb.level, res)
		}
	}
}

func search(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestJump(t *testing.T) {

	var l levelCache
	verFunc := func(name string, abs, pres int) bool {
		rm := search(l.internalCache[abs], name)
		added := search(l.internalCache[pres], name)
		if abs != pres {
			return rm == false && added == true
		}
		return added
	}

	table := []struct {
		name      string
		src, dest int
		verify    func() bool
		isError   bool
	}{
		{"Foo", 0, 1, func() bool { return verFunc("Foo", 0, 1) }, false},
		{"World", 1, 1, func() bool { return verFunc("World", 1, 1) }, false},
		{"Strings", 2, 1, func() bool { return verFunc("Strings", 2, 2) }, false},
		{"Random", 1, 2, func() bool { return verFunc("Random", 1, 2) }, true},
	}

	for _, tb := range table {
		l = createLevelCache()
		err := l.jump(tb.name, tb.src, tb.dest)
		if err != nil && tb.isError != true {
			t.Errorf("Found error %+v where it was not expected \n", err)
		}
		if !tb.verify() {
			t.Errorf("Jump failed for %v\n", tb)
			t.Errorf("Status of cache: %v\n", l.internalCache)
		}
	}
}

func TestCache(t *testing.T) {
	var l levelCache
	dummyFields := make(map[string]*Field)

	cases := []struct {
		gs     *GoStruct
		verify func() bool
	}{
		{
			gs: &GoStruct{Name: "Foo", Level: 0, Fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[0], "Foo")
			},
		},
		{
			gs: &GoStruct{Name: "Foo", Level: 1, Fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[1], "Foo")
			},
		},
		{
			gs: &GoStruct{Name: "World", Level: 0, Fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[1], "World")
			},
		},
		{
			gs: &GoStruct{Name: "New", Level: 3, Fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[3], "New")
			},
		},
	}

	for _, tc := range cases {
		l = createLevelCache()
		l.cache(tc.gs)
		if !tc.verify() {
			t.Errorf("Error in caching %s", tc.gs.Name)
		}
	}
}

func TestIterate(t *testing.T) {

	l := createLevelCache()
	iter := []LevelNames{
		{2, "Few"},
		{2, "More"},
		{2, "Strings"},
		{1, "Hello"},
		{1, "World"},
		{0, "Foo"},
		{0, "Bar"},
		{0, "Baz"},
	}

	resIter := l.Iterate()
	log.Printf("Result: %+v\n", resIter)

	for idx, ln := range resIter {
		eln := iter[idx]
		if ln.name != eln.name {
			t.Errorf("Expected to get %s, but got %s instead. \n",
				eln.name, ln.name)
		}
	}

}
