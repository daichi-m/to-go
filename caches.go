package togo

import (
	"fmt"
	"math"
	"sync"
)

// LevelCache is a structure for level ordering of the structs.
// It wraps around map[int][]string and exposes utility functions that
// makes it easy for different operations.
type levelCache struct {
	internalCache map[int][]string
	maxLevel      int
}

// Search a string inside the level. Returns the index at that level if found
// or -1 if not found
func (l levelCache) searchInLevel(lvl int, name string) int {
	lvlCache, ok := l.internalCache[lvl]
	if !ok {
		return -1
	}
	for idx, s := range lvlCache {
		if s == name {
			return idx
		}
	}
	return -1
}

// Find the level of the name in the cache. Returns the level if found,
// or -1 if not found
func (l levelCache) findLevel(name string) int {
	for lvl := range l.internalCache {
		idx := l.searchInLevel(lvl, name)
		if idx != -1 {
			return lvl
		}
	}
	return -1
}

// Jumps a name from src level to dest level. If the src and dest level
// are same, then do nothing. If the jump is to the lower level, ignore the
// jump in that case. Otherwise jump to a higher level.
func (l levelCache) jump(name string, src, dest int) error {

	if src == dest {
		if _, ok := l.internalCache[src]; !ok {
			return GenericCacheError{
				cacheName: "levelCache",
				element:   name,
				message: fmt.Sprintf("%s does not exist in cache for jump",
					name),
			}
		}
		return nil
	}

	if src < dest {
		srcSl, okSrc := l.internalCache[src]
		destSl, okDest := l.internalCache[dest]

		if !okSrc {
			return GenericCacheError{
				cacheName: "levelCache",
				element:   name,
				message: fmt.Sprintf("%s not present in level %d",
					name, src),
			}
		}

		if !okDest {
			destSl = make([]string, 0, 10)
		}
		destSl = append(destSl, name)
		l.internalCache[dest] = destSl

		idx := l.searchInLevel(src, name)
		if idx == -1 {
			return GenericCacheError{
				cacheName: "levelCache",
				element:   name,
				message: fmt.Sprintf("%s is not present in level %d",
					name, src),
			}
		}
		srcSl = append(srcSl[:idx], srcSl[idx+1:]...)
		l.internalCache[src] = srcSl
	}

	if src > dest {
		return nil
	}

	return nil
}

// Adds an instance of GoStruct to the levelCache
func (l levelCache) cache(gs *GoStruct) error {
	lvl := gs.Level
	existLvl := l.findLevel(gs.Name)

	if existLvl == -1 {
		lvlCache, ok := l.internalCache[lvl]
		if !ok {
			lvlCache = make([]string, 0, 10)
			lvlCache = append(lvlCache, gs.Name)
			l.internalCache[lvl] = lvlCache
			l.maxLevel = int(math.Max(float64(l.maxLevel), float64(lvl)))
			return nil
		}
	} else if existLvl >= lvl {
		l.maxLevel = int(math.Max(float64(l.maxLevel), float64(lvl)))
	} else if existLvl < lvl {
		err := l.jump(gs.Name, existLvl, lvl)
		if err != nil {
			return err
		}
		l.maxLevel = int(math.Max(float64(l.maxLevel), float64(lvl)))
	}
	return nil
}

// LevelNames is the struct that is used to iterate over the level cache
// in a level-order view.
type LevelNames struct {
	level int
	name  string
}

// Returns a slice of LevelNames in decreasing order of level so that
// iteration is simpler over the levels.
func (l levelCache) Iterate() []LevelNames {
	var slc []LevelNames
	max := l.maxLevel
	slc = make([]LevelNames, 0, 10)
	for level := max; level >= 0; level-- {
		names := l.internalCache[level]
		for _, n := range names {
			ln := LevelNames{
				level: level,
				name:  n,
			}
			slc = append(slc, ln)
		}
	}
	return slc
}

// Caches structure for storing various data structures.
type Caches struct {
	NameCache     map[string]bool
	GoStructCache map[string]*GoStruct
	LevelCache    levelCache
	MaxLevel      int
}

var instance *Caches
var lock sync.Mutex

// GetOrCreate creates a new instance of Caches or gets the existing instance
// if it was already created.
func GetOrCreate() *Caches {

	if instance != nil {
		return instance
	}

	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Caches{
				NameCache:     make(map[string]bool),
				GoStructCache: make(map[string]*GoStruct),
				LevelCache: levelCache{
					internalCache: make(map[int][]string),
					maxLevel:      -1,
				},
				MaxLevel: -1,
			}
		}
	}
	return instance
}

// CacheStruct caches the struct represent by the GoStruct pointer into the
// name and level caches.
func (cache *Caches) CacheStruct(gs *GoStruct) error {
	name := gs.Name
	cache.GoStructCache[name] = gs
	err := cache.LevelCache.cache(gs)
	if err != nil {
		return err
	}
	cache.MaxLevel = int(math.Max(float64(cache.MaxLevel), float64(gs.Level)))
	return nil
}

// CacheName caches a name into the name cache. If the name already exists
// it adds a number to its end (1 to 100) until it can find a unique name.
// In the unlikely scenario all 100 numbers are used up, it will return
// a NameClashError
func (cache *Caches) CacheName(name string) string {
	ex := cache.Exist(name)
	if !ex {
		cache.NameCache[name] = true
		return name
	}

	i := 1
	for i < 100 {
		nm := name + "_" + string(i)
		if !cache.Exist(nm) {
			cache.NameCache[nm] = true
			return nm
		}
	}
	return ""
}

// Exist returns true if the name has already been used.
func (cache *Caches) Exist(name string) bool {
	ex := cache.NameCache[name]
	return ex
}

/*
func (cache *Caches) CacheGoStruct(gs *GoStruct) error {
    levelSlice, okLvl := cache.LevelCache[gs.Level]
    if !okLvl {
        levelSlice = make([]string, 1)
        cache.LevelCache[gs.Level] = levelSlice
    }
	cachedGs, okGs := cache.GoStructCache[gs.Name]




	levelSlice, ok1 := cache.LevelCache[gs.Level]

    if !okLvl && levelSlice contains gs {
        return NewCacheError("GoStructCache", gs.Name)
    }

	if !okLvl {
		cache.GoStructCache[gsName] = gs
		if !ok1 {
			levelSlice = make([]*GoStruct, 1)
		}
		levelSlice = append(levelSlice, gs)
		cache.LevelCache[gs.Level] = levelSlice
		return nil
	}

	if cachedGs.Level == gs.Level {
		cache.GoStructCache[gs.Name] = gs
		return nil
	} else {

	}

}*/
