package types

// StringHashSet is a type for a string hash set.
type StringHashSet struct {
	m map[string]struct{}
}

// NewStringHashSet creates a new string hash set.
func NewStringHashSet() *StringHashSet {
	return &StringHashSet{
		m: make(map[string]struct{}),
	}
}

// Add adds a value to the hash set.
func (s *StringHashSet) Add(v string) {
	s.m[v] = struct{}{}
}

// Contains checks whether the specified value is in the hash set.
func (s *StringHashSet) Contains(v string) bool {
	_, exists := s.m[v]
	return exists
}

// IsEmpty checks whether the hash set is empty.
func (s *StringHashSet) IsEmpty() bool {
	return len(s.m) == 0
}

// Remove removes a value from the hash set.
func (s *StringHashSet) Remove(v string) {
	delete(s.m, v)
}

// Values gets the values from the hash set.
func (s *StringHashSet) Values() []string {
	var values []string
	for k, _ := range s.m {
		values = append(values, k)
	}
	return values
}

// UInt64HashSet is a type for a uint64 hash set.
type UInt64HashSet struct {
	m map[uint64]struct{}
}

// NewUInt64HashSet creates a new string hash set.
func NewUInt64HashSet() *UInt64HashSet {
	return &UInt64HashSet{
		m: make(map[uint64]struct{}),
	}
}

// Add adds a value to the hash set.
func (s *UInt64HashSet) Add(v uint64) {
	s.m[v] = struct{}{}
}

// Contains checks whether the specified value is in the hash set.
func (s *UInt64HashSet) Contains(v uint64) bool {
	_, exists := s.m[v]
	return exists
}

// IsEmpty checks whether the hash set is empty.
func (s *UInt64HashSet) IsEmpty() bool {
	return len(s.m) == 0
}

// Remove removes a value from the hash set.
func (s *UInt64HashSet) Remove(v uint64) {
	delete(s.m, v)
}

// Values gets the values from the hash set.
func (s *UInt64HashSet) Values() []uint64 {
	var values []uint64
	for k, _ := range s.m {
		values = append(values, k)
	}
	return values
}
