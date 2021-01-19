package page

// exists is an empty struct that holds 0-byte in memory.
// It's a cleaner way of using an empty struct which we use
// as the value in our internal set map
var exists = struct{}{}

// Set is a thread-safe Set data structure
type Set struct {
	// data is an internal map which holds our Set data
	data map[string]struct{}
}

// NewSet initializes a new set
func NewSet() *Set {
	data := make(map[string]struct{})

	return &Set{data: data}
}

// Add adds a new item to the set
func (s *Set) Add(key string) {
	s.data[key] = exists
}

// Has checks for the existence of a key in the set
func (s *Set) Has(key string) bool {
	_, ok := s.data[key]
	return ok
}
