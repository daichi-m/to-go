package togo

import (
	"fmt"
	"sync"
)

type IntIntPair struct {
	Left, Right int
}

// CacheIterator is an iterator to iterate over the values in a cache
type CacheIterator interface {
	// HasNext indicates if the iterator have a next element
	HasNext() bool
	// Next returns the next element of the iterator
	Next() (interface{}, bool)
	// Close ends the iterator
	Close()
}

// LevelCache is a structure for level ordering of the structs.
// It wraps around map[int][]string and exposes utility functions that
// makes it easy for different operations.
type LevelCache struct {
	levelToName map[int][]string
	nameToLevel map[string]int
	maxLevel    int
}

type levelCacheIterator struct {
	LevelCache
	current IntIntPair
	active  bool
}

func (lci *levelCacheIterator) HasNext() bool {
	if !lci.active {
		return false
	}

	currLevel := lci.current.Left
	if currLevel == -1 {
		return true
	}

	if currLevel == lci.maxLevel {
		ln := lci.levelToName[lci.maxLevel]
		return lci.current.Right < (len(ln) - 1)
	}
	return true
}

func (lci *levelCacheIterator) Next() (interface{}, bool) {

	if !lci.active {
		return "", false
	}

	currLevel := lci.current.Left
	if currLevel == -1 {
		lci.current.Left = 0
		currLevel = 0
	}
	ln := lci.levelToName[currLevel]
	if lci.current.Right < (len(ln) - 1) {
		lci.current.Right++
		x := ln[lci.current.Right]
		return x, true
	}

	currLevel++
	ln, ok := lci.levelToName[currLevel]
	for !ok || len(ln) == 0 {
		currLevel++
		if currLevel > lci.maxLevel {
			break
		}
		ln, ok = lci.levelToName[currLevel]

	}

	if ok {
		lci.current.Left = currLevel
		x := ln[0]
		lci.current.Right = 0
		return x, true
	}

	return "", false
}

func (lci *levelCacheIterator) Close() {
	lci.active = false
	lci.current = IntIntPair{lci.maxLevel + 1, -1}
}

// LevelNames is the struct that is used to iterate over the level cache
// in a level-order view.
type LevelNames struct {
	Level int
	Name  string
}

// NewLevelCache creates a new instance of LevelCache
func NewLevelCache() *LevelCache {
	return &LevelCache{
		levelToName: make(map[int][]string),
		nameToLevel: make(map[string]int),
		maxLevel:    -1,
	}
}

func (lc *LevelCache) getNameSlice(level int) []string {
	sl, ok := lc.levelToName[level]
	if !ok {
		sl = make([]string, 0, 10)
		lc.levelToName[level] = sl
	}
	return sl
}

// Cache a (name, level) entry into the LevelCache
func (lc *LevelCache) Cache(name string, level int) error {

	srcLvl, ok := lc.nameToLevel[name]
	if lc.maxLevel < level {
		lc.maxLevel = level
	}
	if !ok {
		// New data addition
		nameSlc := lc.getNameSlice(level)
		nameSlc = append(nameSlc, name)
		lc.nameToLevel[name] = level
		lc.levelToName[level] = nameSlc
		return nil
	}
	return lc.Jump(name, srcLvl, level)
}

// DeCache removes a name from the level cache
func (lc *LevelCache) DeCache(name string) {

	lvl, ok := lc.nameToLevel[name]
	if !ok {
		return
	}
	slc, ok := lc.levelToName[lvl]
	if !ok {
		return
	}
	slc2 := make([]string, 0, len(slc)-1)
	for _, s := range slc {
		if s == name {
			continue
		}
		slc2 = append(slc2, s)
	}

	if len(slc2) == 0 {
		delete(lc.levelToName, lvl)
	} else {
		lc.levelToName[lvl] = slc2
	}
	delete(lc.nameToLevel, name)
}

// Jump moves a name from src level to dest level. If the src and dest level
// are same, then do nothing. If the Jump is to the lower level, ignore the
// jump in that case. Otherwise jump to a higher level.
func (lc *LevelCache) Jump(name string, src, dest int) error {
	lvl, ok := lc.nameToLevel[name]
	if !ok {
		return fmt.Errorf("Name %s not found in LevelCache", name)
	}
	if lvl != src {
		return fmt.Errorf("Name %s not found at level %d", name, src)
	}

	if src >= dest {
		return nil
	}

	lc.DeCache(name)
	return lc.Cache(name, dest)
}

// Location gets the (level, index) of a name in the cache
func (lc *LevelCache) Location(name string) (IntIntPair, bool) {
	lvl, ok := lc.nameToLevel[name]
	if !ok {
		return IntIntPair{-1, -1}, false
	}
	slc := lc.levelToName[lvl]
	for i, s := range slc {
		if s == name {
			return IntIntPair{lvl, i}, true
		}
	}
	return IntIntPair{-1, -1}, false
}

// Iterator returns an CacheIterator for the LevelCache
func (lc *LevelCache) Iterator() CacheIterator {
	return &levelCacheIterator{
		LevelCache: *lc,
		current:    IntIntPair{-1, -1},
		active:     true,
	}
}

// MaxLevel gets the maximum level the LevelCache has
func (lc *LevelCache) MaxLevel() int {
	return lc.maxLevel
}

// StructCache structure for storing various data structures.
type StructCache struct {
	nameCache  map[string]bool
	igsCache   map[string]IGoStruct
	levelCache *LevelCache
	maxLevel   int
}

type structCacheIterator struct {
	*StructCache
	lci    CacheIterator
	active bool
}

var _ CacheIterator = (*structCacheIterator)(nil)

func (sci *structCacheIterator) HasNext() bool {
	if !sci.active {
		return false
	}
	return sci.lci.HasNext()
}

func (sci *structCacheIterator) Next() (interface{}, bool) {
	if !sci.active {
		return nil, false
	}

	ln, ok := sci.lci.Next()
	if !ok {
		return nil, false
	}
	gs, ok := sci.igsCache[ln.(string)]
	if !ok {
		return nil, false
	}
	return gs, true
}

func (sci *structCacheIterator) Close() {
	sci.active = false
	sci.lci.Close()
}

var instance *StructCache
var lock sync.Mutex

// GetOrCreateCaches creates a new instance of Caches or gets the existing instance
// if it was already created.
func GetOrCreateCaches() *StructCache {

	if instance != nil {
		return instance
	}

	instance = &StructCache{
		nameCache:  make(map[string]bool),
		igsCache:   make(map[string]IGoStruct),
		levelCache: NewLevelCache(),
		maxLevel:   -1,
	}
	return instance
}

// CacheStruct caches the struct represent by the GoStruct pointer into the
// name and level caches.
func (cache *StructCache) CacheStruct(gs IGoStruct, uniq bool) error {
	name := gs.Name()
	var err error
	if uniq {
		name, err = cache.CacheName(name)
		if err != nil {
			return err
		}
		gs.ChangeName(name)
	} else {
		cache.CacheNameWithConflict(name)
	}

	err = cache.levelCache.Cache(gs.Name(), gs.Level())
	if err != nil {
		return err
	}
	if gs.Level() > cache.maxLevel {
		cache.maxLevel = gs.Level()
	}
	cache.igsCache[gs.Name()] = gs
	return nil
}

// CacheName caches a name into the name cache. If the name already exists
// it adds a number to its end (1 to 100) until it can find a unique name.
// In the unlikely scenario all 100 numbers are used up, it will return
// a NameClashError
func (cache *StructCache) CacheName(name string) (string, error) {
	ex := cache.Exist(name)
	if !ex {
		cache.nameCache[name] = true
		return name, nil
	}

	for i := 1; i < 10; i++ {
		nm := fmt.Sprintf("%s_%d", name, i)
		if !cache.Exist(nm) {
			cache.nameCache[nm] = true
			return nm, nil
		}
	}
	return "", fmt.Errorf("NameConflict: Cannot find suitable alternative for %s", name)
}

// CacheNameWithConflict caches the name irrespective of a clash. If the name already exists,
// this method just returns silently
func (cache *StructCache) CacheNameWithConflict(name string) {
	ex := cache.Exist(name)
	if ex {
		return
	}
	cache.CacheName(name)
}

// Exist returns true if the name has already been used.
func (cache *StructCache) Exist(name string) bool {
	ex := cache.nameCache[name]
	return ex
}

// MaxLevel returns the max level in the Caches
func (cache *StructCache) MaxLevel() int {
	return cache.maxLevel
}

// Iterator returns a CacheIterator to iterate over the cache
func (cache *StructCache) Iterator() CacheIterator {
	return &structCacheIterator{
		StructCache: cache,
		lci:         cache.levelCache.Iterator(),
		active:      true,
	}
}
