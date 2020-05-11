package togo

import "fmt"

// CacheAddError is an error type describing an error while adding an element to a cache
type CacheAddError struct {
	cacheName   string
	elementName string
}

func (c CacheAddError) Error() string {
	return fmt.Sprintf("Error in caching %s to cache %s", c.elementName, c.cacheName)
}

// GenericCacheError is a generic error in handling of the cache as verified by the message
type GenericCacheError struct {
	cacheName string
	element   string
	message   string
}

func (c GenericCacheError) Error() string {
	return fmt.Sprintf("Error while working with cache %s and element %s: %s", c.cacheName, c.element, c.message)
}
