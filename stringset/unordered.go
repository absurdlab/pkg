package stringset

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
)

// NewUnOrdered returns a new unordered string set.
func NewUnOrdered() *UnOrdered {
	return &UnOrdered{
		index: map[string]struct{}{},
	}
}

// NewUnOrderedWith returns a new unordered string set with given elements.
func NewUnOrderedWith(elements ...string) *UnOrdered {
	set := NewUnOrdered()
	AddAll(set, elements...)
	return set
}

type UnOrdered struct {
	index map[string]struct{}
}

func (s *UnOrdered) Size() int {
	if s == nil {
		return 0
	}
	return len(s.index)
}

func (s *UnOrdered) Contains(element string) bool {
	if IsEmpty(s) {
		return false
	}

	_, ok := s.index[element]
	return ok
}

func (s *UnOrdered) Add(element string) {
	if s == nil {
		*s = *NewUnOrdered()
	}

	s.index[element] = struct{}{}
}

func (s *UnOrdered) Remove(element string) {
	if s == nil {
		return
	}

	delete(s.index, element)
}

func (s *UnOrdered) All(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for each := range s.index {
		if ok := criteria(each); !ok {
			return false
		}
	}

	return true
}

func (s *UnOrdered) Any(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for each := range s.index {
		if ok := criteria(each); ok {
			return true
		}
	}

	return false
}

func (s *UnOrdered) None(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for each := range s.index {
		if ok := criteria(each); ok {
			return false
		}
	}

	return true
}

func (s *UnOrdered) Array() []string {
	if IsEmpty(s) {
		return []string{}
	}

	dup := make([]string, 0)

	for each := range s.index {
		dup = append(dup, each)
	}

	return dup
}

func (s *UnOrdered) One() string {
	if IsEmpty(s) {
		panic("unordered string set is empty")
	}

	for each := range s.index {
		return each
	}

	panic("unreachable code")
}

func (s *UnOrdered) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte{}, nil
	}
	return json.Marshal(s.Array())
}

func (s *UnOrdered) UnmarshalJSON(bytes []byte) error {
	var elem []string
	if err := json.Unmarshal(bytes, &elem); err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	set := NewUnOrdered()
	AddAll(set, elem...)
	*s = *set

	return nil
}

func (s *UnOrdered) MarshalYAML() (interface{}, error) {
	if s == nil {
		return []byte{}, nil
	}
	return yaml.Marshal(s.Array())
}

func (s *UnOrdered) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var elem []string
	if err := unmarshal(&elem); err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	set := NewUnOrdered()
	AddAll(set, elem...)
	*s = *set

	return nil
}

func (s *UnOrdered) IsZero() bool {
	return IsEmpty(s)
}

func (s *UnOrdered) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case []byte:
		source = src.([]byte)
	case string:
		source = []byte(src.(string))
	default:
		return errors.New("incompatible type for unordered string set")
	}

	return json.Unmarshal(source, s)
}

func (s *UnOrdered) Value() (driver.Value, error) {
	return json.Marshal(s)
}
