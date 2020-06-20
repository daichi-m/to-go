package togo

import (
	"fmt"
	"log"
	"testing"
)

func search(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestLevelCache_searchInLevel(t *testing.T) {

	l := createLevelCache()
	table := []struct {
		tc     string
		level  int
		name   string
		result int
	}{
		{"0th Level", 0, "Foo", 0},
		{"1st Level", 1, "World", 1},
		{"2nd Level", 2, "Few", 0},
		{"Non Existing Level", 3, "Example", -1},
		{"Not Exist in Level", 0, "Few", -1},
	}

	for _, tb := range table {
		t.Run(tb.tc, func(t *testing.T) {
			res := l.searchInLevel(tb.level, tb.name)
			if res != tb.result {
				t.Errorf("TC: %s: Expected %v but found %v \n", tb.tc, tb.result, res)
			}
		})
	}
}

func TestLevelCache_findLevel(t *testing.T) {

	l := createLevelCache()
	table := []struct {
		tc    string
		name  string
		level int
	}{
		{"0th Level", "Foo", 0},
		{"1st Level", "World", 1},
		{"2nd Level", "Few", 2},
		{"Not Foun", "Random", -1},
	}

	for _, tb := range table {
		t.Run(tb.tc, func(t *testing.T) {
			res := l.findLevel(tb.name)
			if res != tb.level {
				t.Errorf("TC: %s: Expected %v but found %v\n", tb.tc, tb.level, res)
			}
		})
	}
}

func TestLevelCache_jump(t *testing.T) {

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
		tc        string
		name      string
		src, dest int
		verify    func() bool
		isError   bool
	}{
		{"Jump higher", "Foo", 0, 1, func() bool { return verFunc("Foo", 0, 1) }, false},
		{"Jump same", "World", 1, 1, func() bool { return verFunc("World", 1, 1) }, false},
		{"Jump lower", "Strings", 2, 1, func() bool { return verFunc("Strings", 2, 2) }, false},
		{"Jump non-existent", "Random", 1, 2, func() bool { return verFunc("Random", 1, 2) }, true},
	}

	for _, tb := range table {
		t.Run(tb.tc, func(t *testing.T) {
			l = createLevelCache()
			err := l.jump(tb.name, tb.src, tb.dest)
			if err != nil && tb.isError != true {
				t.Errorf("TC: %s, Found error %+v where it was not expected \n", tb.tc, err)
			}
			if !tb.verify() {
				t.Errorf("TC: %s, Jump failed for %v\n", tb.tc, tb)
				t.Errorf("Status of cache: %v\n", l.internalCache)
			}
		})
	}
}

func TestLevelCache_cache(t *testing.T) {
	var l levelCache
	dummyFields := make(map[string]*Field)

	cases := []struct {
		tc     string
		gs     *GoStruct
		verify func() bool
	}{
		{
			tc: "Level 0 Struct",
			gs: &GoStruct{name: "Foo", level: 0, fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[0], "Foo")
			},
		},
		{
			tc: "Level 1 struct",
			gs: &GoStruct{name: "Foo", level: 1, fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[1], "Foo")
			},
		},
		{
			tc: "Level 0 Struct already in level 1",
			gs: &GoStruct{name: "World", level: 0, fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[1], "World")
			},
		},
		{
			tc: "New level 3 struct",
			gs: &GoStruct{name: "New", level: 3, fields: dummyFields},
			verify: func() bool {
				return search(l.internalCache[3], "New")
			},
		},
	}

	for _, tb := range cases {
		t.Run(tb.tc, func(t *testing.T) {
			l = createLevelCache()
			l.cache(tb.gs)
			if !tb.verify() {
				t.Errorf("TC: %s: Error in caching %s", tb.tc, tb.gs.name)
			}
		})
	}
}

func TestLevelCache_iterate(t *testing.T) {

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

func TestCaches_CacheStruct(t *testing.T) {
	caches := createCaches()
	type testCase struct {
		tc   string
		gs   GoStruct
		uniq bool
	}
	verify := func(tc testCase) bool {
		v, ok := caches.GoStructCache[tc.gs.name]
		if !ok || v.name != tc.gs.name {
			return false
		}
		vx, ok := caches.NameCache[tc.gs.name]
		if !ok || vx != true {
			return false
		}
		return true
	}

	gs := createSimpleGoStruct()
	table := []testCase{
		{
			tc:   "New Struct",
			gs:   gs,
			uniq: true,
		},
		{
			tc:   "Existing Struct",
			gs:   gs,
			uniq: false,
		},
	}
	for _, tc := range table {
		t.Run(tc.tc, func(t *testing.T) {
			_ = caches.CacheStruct(&tc.gs, tc.uniq)
			if !verify(tc) {
				t.Errorf("TC: %s: verification function failed", tc.tc)
			}
		})
	}
}

func TestCaches_CacheName(t *testing.T) {

	caches := createCaches()
	setup := func(mp map[string]bool, nm string) {
		for i := 1; i < 10; i++ {
			name := fmt.Sprintf("%s_%d", nm, i)
			mp[name] = true
		}
	}
	setup(caches.NameCache, "Bar")

	verify := func(name string) bool {
		v, ok := caches.NameCache[name]
		return ok == true && v == true
	}

	table := []struct {
		tc          string
		name        string
		expected    string
		expectError bool
	}{
		{
			tc:          "New Name",
			name:        "Pickle",
			expected:    "Pickle",
			expectError: false,
		},
		{
			tc:          "Existing Name",
			name:        "Foo",
			expected:    "Foo_1",
			expectError: false,
		},
		{
			tc:          "Name Conflict Error",
			name:        "Bar",
			expected:    "",
			expectError: true,
		},
	}
	for _, tc := range table {
		t.Run(tc.tc, func(t *testing.T) {
			got, err := caches.CacheName(tc.name)
			if err == nil && !verify(tc.expected) {
				t.Errorf("TC: %s: verification function failed", tc.tc)
			}
			if got != tc.expected && (!tc.expectError && err != nil) {
				t.Errorf("TC: %s: Expected %v but got %v instead", tc.tc, tc.expected, got)
			}
		})
	}
}

func TestCaches_CacheNameErrorFree(t *testing.T) {
	caches := createCaches()
	verify := func(name string) bool {
		v, ok := caches.NameCache[name]
		return ok == true && v == true
	}

	table := []struct {
		tc   string
		name string
	}{
		{
			tc:   "New Name",
			name: "Pickle",
		},
		{
			tc:   "Existing Name",
			name: "Foo",
		},
	}
	for _, tb := range table {
		t.Run(tb.name, func(t *testing.T) {
			caches.CacheNameErrorFree(tb.name)
			if !verify(tb.name) {
				t.Errorf("TC: %s: verification function failed", tb.tc)
			}
		})
	}
}

func TestCaches_Exist(t *testing.T) {
	caches := createCaches()
	table := []struct {
		tc     string
		name   string
		verify func() bool
		want   bool
	}{
		{
			tc:   "Non Existing Name",
			name: "Pickle",
			verify: func() bool {
				_, ok := caches.NameCache["Pickle"]
				return !ok
			},
			want: false,
		},
		{
			tc:   "Existing Name",
			name: "Foo",
			verify: func() bool {
				v, ok := caches.NameCache["Foo"]
				return ok == true && v == true
			},
			want: true,
		},
	}
	for _, tb := range table {
		t.Run(tb.tc, func(t *testing.T) {
			got := caches.Exist(tb.name)
			if !tb.verify() {
				t.Errorf("TC: %s: verification function failed", tb.tc)
			}
			if got != tb.want {
				t.Errorf("TC: %s: Expected %v but got %v instead", tb.tc, tb.want, got)
			}
		})
	}
}

func TestGetOrCreate(t *testing.T) {

	var retCache *Caches
	tests := []struct {
		tc     string
		verify func(c *Caches) bool
	}{
		{
			tc: "First Call",
			verify: func(c *Caches) bool {
				if len(c.GoStructCache) != 0 ||
					len(c.NameCache) != 0 ||
					len(c.LevelCache.internalCache) != 0 {
					return false
				}
				return true
			},
		},
		{
			tc: "Subsequent Call",
			verify: func(c *Caches) bool {
				if c != retCache {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			c := GetOrCreate()
			if !tt.verify(c) {
				t.Errorf("TC: %s: Verfication function failed", tt.tc)
			}
			if retCache == nil {
				retCache = c
			}
		})
	}
}
