package togo

import (
	"math/rand"
)

// Various utility methods used for testing

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

func createName() string {
	names := []string{
		"Foo", "Bar", "Baz", "Blot", "Fry",
	}
	idx := rand.Intn(5)
	return names[idx]
}

func createField(fdt FieldDT, strType string, nest int) Field {
	name := createName()
	field := Field{
		name:         name,
		annotation:   "`json:`" + name,
		dataType:     fdt,
		dtStruct:     strType,
		sliceNesting: nest,
	}
	return field
}

func createPrimitiveField(fdt FieldDT) Field {
	return createField(fdt, "", 0)
}

func createIntegerField() Field {
	return createPrimitiveField(Int)
}

func createStringField() Field {
	return createPrimitiveField(String)
}

func createFloatField() Field {
	return createPrimitiveField(Float64)
}

func createMapField(gs GoStruct) Field {
	return createField(Map, gs.Name, 0)
}

func creteSliceField(gs GoStruct, nest int) Field {
	return createField(Slice, gs.Name, nest)
}

func createSimpleGoStruct() GoStruct {
	name := createName()
	gs := GoStruct{
		Name:   name,
		Fields: make(map[string]*Field),
		Level:  1,
	}
	for i := 0; i < 5; i++ {
		fld := createIntegerField()
		gs.Fields[fld.name] = &fld
	}
	return gs
}

func createCaches() Caches {
	caches := Caches{
		NameCache:     make(map[string]bool),
		GoStructCache: make(map[string]*GoStruct),
		LevelCache:    createLevelCache(),
		MaxLevel:      -1,
	}
	var gs []GoStruct
	for i := 0; i < 5; i++ {
		gs = append(gs, createSimpleGoStruct())
	}

	for _, g := range gs {
		caches.NameCache[g.Name] = true
		caches.GoStructCache[g.Name] = &g
		caches.LevelCache.cache(&g)
		if g.Level > caches.MaxLevel {
			caches.MaxLevel = g.Level
		}
	}
	return caches
}
