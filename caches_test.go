package togo_test

import (
	"log"
	"testing"

	. "github.com/daichi-m/togo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func createLevelCache() *LevelCache {
	lc := NewLevelCache()
	iter := namesLevel.Iterator()
	for iter.Next() {
		lc.Cache(iter.Key().(string), iter.Value().(int))
	}
	return lc
}

type ToGoSuite interface {
	suite.SetupTestSuite
	suite.BeforeTest
	suite.AfterTest
	suite.TearDownTestSuite
}

type CacheSuite struct {
	suite.Suite
	name       string
	levelCache *LevelCache
	caches     *StructCache
}

var _ ToGoSuite = (*CacheSuite)(nil)

func (suite *CacheSuite) SetupTest() {
	suite.name = "CachesTestSuite"
	suite.caches = GetOrCreateCaches()
}

func (suite *CacheSuite) BeforeTest(suiteName, testName string) {
	suite.levelCache = createLevelCache()
}

func (suite *CacheSuite) AfterTest(suiteName, testName string) {
	suite.levelCache = nil
	suite.caches = nil
}

func (suite *CacheSuite) TearDownTest() {
	suite.levelCache = nil
	suite.caches = nil
}

func (suite *CacheSuite) Test_LevelCache_Cache() {

	tests := []struct {
		name   string
		level  int
		expLoc IntIntPair
	}{
		{"new", 3, IntIntPair{3, 0}},
		{"new_same_level", 0, IntIntPair{0, 3}},
		{"foo", 0, IntIntPair{0, 0}},
		{"bar", 2, IntIntPair{2, 3}},
	}

	for _, test := range tests {
		suite.levelCache.Cache(test.name, test.level)
		loc, _ := suite.levelCache.Location(test.name)
		assert.Equal(suite.T(), test.expLoc.Left, loc.Left)
		assert.Equal(suite.T(), test.expLoc.Right, loc.Right)
		assert.Equal(suite.T(), 3, suite.levelCache.MaxLevel())
	}
}

func (suite *CacheSuite) Test_LevelCache_Jump() {
	tests := []struct {
		name          string
		src           int
		dest          int
		expectedLevel int
		expectedErr   bool
	}{
		{"foo", 0, 0, 0, false},
		{"qux", 1, 0, 1, false},
		{"bar", 0, 1, 1, false},
		{"unknown", 0, 1, -1, true},
		{"foo", 1, 2, -1, true},
	}

testLoop:
	for _, test := range tests {
		err := suite.levelCache.Jump(test.name, test.src, test.dest)
		if err != nil {
			assert.True(suite.T(), test.expectedErr, "Did not expect error, but found %s",
				err.Error())
			continue testLoop
		}
		loc, _ := suite.levelCache.Location(test.name)
		if !test.expectedErr {
			assert.Equal(suite.T(), test.expectedLevel, loc.Left)
		} else {
			assert.Fail(suite.T(), "Expected error, did not get error")
		}
	}
}

func (suite *CacheSuite) Test_LevelCache_DeCache() {
	tests := []struct {
		name []string
	}{
		{[]string{"foo"}},
		{[]string{"rand"}},
		{[]string{"qux"}},
		{[]string{"foo", "bar", "baz"}},
	}

	for _, test := range tests {
		for _, n := range test.name {
			suite.levelCache.DeCache(n)
		}

		for _, n := range test.name {
			_, ok := suite.levelCache.Location(n)
			assert.Falsef(suite.T(), ok, "Not expected to find %s", n)
		}
	}
}

func (suite *CacheSuite) Test_LevelCache_Location() {
	tests := []struct {
		name   string
		expLoc IntIntPair
		expOk  bool
	}{
		{"foo", IntIntPair{0, 0}, true},
		{"quux", IntIntPair{1, 1}, true},
		{"wobble", IntIntPair{-1, -1}, false},
	}

	for _, test := range tests {
		loc, ok := suite.levelCache.Location(test.name)
		if !ok {
			assert.Equalf(suite.T(), test.expOk, ok, "Expected ok %v found %v", test.expOk, ok)
		} else {
			assert.Equal(suite.T(), test.expLoc.Left, loc.Left)
			assert.Equal(suite.T(), test.expLoc.Right, loc.Right)
		}
	}
}

func (suite *CacheSuite) Test_LevelCache_Iterator() {
	expectedOrder := make([]string, 9)
	copy(expectedOrder, names[0:9])
	n2 := expectedOrder[5]
	expectedOrder = append(expectedOrder, "hobb")
	expectedOrder = append(expectedOrder[0:5], expectedOrder[6:]...)
	expectedOrder = append(expectedOrder, n2)
	suite.levelCache.Cache("hobb", 6)
	suite.levelCache.Jump(n2, 1, 6)

	actualOrder := make([]string, 0, 10)
	iter := suite.levelCache.Iterator()
	for iter.HasNext() {
		n, ok := iter.Next()
		assert.True(suite.T(), ok)
		log.Printf("Got next in iteration: %v, %v\n", n, ok)
		actualOrder = append(actualOrder, n.(string))
		if len(actualOrder) > len(expectedOrder) {
			assert.LessOrEqual(suite.T(), len(actualOrder), len(expectedOrder))
			return
		}
	}

	for i := range expectedOrder {
		if len(actualOrder) <= i {
			assert.Fail(suite.T(), "Length of order is less than %d", i)
		}
		assert.Equal(suite.T(), expectedOrder[i], actualOrder[i])
	}

	iter.Close()
	_, ok := iter.Next()
	assert.False(suite.T(), ok)

}

type TestGoStruct struct {
	mock.Mock
	name string
}

var _ IGoStruct = (*TestGoStruct)(nil)

func (tgs *TestGoStruct) String() string {
	return "TestGoStruct"
}

func (tgs *TestGoStruct) GoString() string {
	return tgs.String()
}

func (tgs *TestGoStruct) Clone() IGoStruct {
	args := tgs.Called()
	return args.Get(0).(*TestGoStruct)
}

func (tgs *TestGoStruct) Grow(gs IGoStruct) (IGoStruct, error) {
	args := tgs.Called(gs)
	return tgs, args.Error(1)
}

func (tgs *TestGoStruct) Equals(gs IGoStruct) bool {
	args := tgs.Called(gs)
	return args.Bool(0)
}

func (tgs *TestGoStruct) Name() string {
	args := tgs.Called()
	if tgs.name == "" {
		return args.String(0)
	}
	return tgs.name
}

func (tgs *TestGoStruct) ChangeName(name string) {
	tgs.name = name
}

func (tgs *TestGoStruct) AddField(IField) (IGoStruct, error) {
	args := tgs.Called()
	return tgs, args.Error(0)
}

func (tgs *TestGoStruct) Level() int {
	args := tgs.Called()
	return args.Int(0)
}

func (suite *CacheSuite) Test_StructCache_CacheStruct() {

	tests := []struct {
		structName  string
		structLevel int
		uniq        bool
		expErr      bool
	}{
		{"Foo", 0, true, false},
		{"Bar", 1, true, false},
		{"Baz", 2, true, false},
		{"Qux", 1, true, false},
		{"Quux", 1, true, false},
		{"Quz", 0, true, false},
		{"Quux", 2, true, false},
		{"Quux", 1, false, false},
	}
	expectOrder := []string{
		"Foo", "Quz", "Bar", "Qux", "Quux", "Baz", "Quux_1",
	}

	for _, test := range tests {
		mgs := new(TestGoStruct)
		mgs.On("Name").Return(test.structName)
		mgs.On("Level").Return(test.structLevel)
		err := suite.caches.CacheStruct(mgs, test.uniq)
		if err != nil {
			assert.True(suite.T(), test.expErr)
		}
		if test.expErr {
			assert.NotNil(suite.T(), err)
		}
	}

	iter := suite.caches.Iterator()
	actualOrder := make([]string, 0, 10)
outer:
	for iter.HasNext() {
		it, ok := iter.Next()
		if !ok {
			assert.Fail(suite.T(), "Expected ok, but received not ok")
		}
		gs := it.(IGoStruct)
		actualOrder = append(actualOrder, gs.Name())

		if len(actualOrder) > len(expectOrder) {
			assert.LessOrEqual(suite.T(), len(actualOrder), len(expectOrder))
			break outer
		}
	}

	for i := range expectOrder {
		assert.Equal(suite.T(), expectOrder[i], actualOrder[i])
	}
	log.Println(actualOrder)

	assert.Equal(suite.T(), suite.caches.MaxLevel(), 2)
}

func Test_CacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
