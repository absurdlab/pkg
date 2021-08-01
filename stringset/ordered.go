package stringset

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"strings"
)

var (
	_ Interface = (*Ordered)(nil)
)

// NewOrdered returns a new order string set.
func NewOrdered() *Ordered {
	return &Ordered{
		elem:  []string{},
		index: map[string]struct{}{},
	}
}

// NewOrderedWith returns a new order string set with given elements.
func NewOrderedWith(elements ...string) *Ordered {
	set := NewOrdered()
	AddAll(set, elements...)
	return set
}

// NewOrderedBy returns a new order string set with delimited elements from the value string.
func NewOrderedBy(value string, delimiter rune) *Ordered {
	elements := strings.FieldsFunc(value, func(r rune) bool {
		return r == delimiter
	})
	return NewOrderedWith(elements...)
}

// NewOrderedBySpace is NewOrderedBy with space delimiter.
func NewOrderedBySpace(value string) *Ordered {
	return NewOrderedBy(value, ' ')
}

// Ordered is a ordered set data structure for strings.
type Ordered struct {
	elem  []string
	index map[string]struct{}
}

func (s *Ordered) Size() int {
	if s == nil {
		return 0
	}
	return len(s.elem)
}

func (s *Ordered) Contains(element string) bool {
	if IsEmpty(s) {
		return false
	}

	_, ok := s.index[element]
	return ok
}

func (s *Ordered) Add(element string) {
	if s == nil {
		return
	}

	if !s.Contains(element) {
		s.elem = append(s.elem, element)
		s.index[element] = struct{}{}
	}
}

func (s *Ordered) Remove(element string) {
	if IsEmpty(s) {
		return
	}

	if i := s.indexOf(element); i >= 0 {
		s.elem = append(s.elem[:i], s.elem[i+1:]...)
		delete(s.index, element)
	}
}

func (s *Ordered) indexOf(elem string) int {
	for i, it := range s.elem {
		if it == elem {
			return i
		}
	}
	return -1
}

func (s *Ordered) All(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for _, each := range s.elem {
		if ok := criteria(each); !ok {
			return false
		}
	}

	return true
}

func (s *Ordered) Any(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for _, each := range s.elem {
		if ok := criteria(each); ok {
			return true
		}
	}

	return false
}

func (s *Ordered) None(criteria Criteria) bool {
	if IsEmpty(s) {
		return false
	}

	for _, each := range s.elem {
		if ok := criteria(each); ok {
			return false
		}
	}

	return true
}

func (s *Ordered) Array() []string {
	if s == nil || len(s.elem) == 0 {
		return []string{}
	}

	dup := make([]string, len(s.elem))
	copy(dup, s.elem)

	return dup
}

func (s *Ordered) One() string {
	if IsEmpty(s) {
		panic("order string set is empty")
	}
	return s.elem[0]
}

func (s *Ordered) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(s.elem)
}

func (s *Ordered) UnmarshalJSON(bytes []byte) error {
	var elem []string
	if err := json.Unmarshal(bytes, &elem); err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	set := NewOrdered()
	AddAll(set, elem...)
	*s = *set

	return nil
}

func (s *Ordered) MarshalYAML() (interface{}, error) {
	if s == nil {
		return []byte("null\n"), nil
	}
	return yaml.Marshal(s.elem)
}

func (s *Ordered) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var elem []string
	if err := unmarshal(&elem); err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	set := NewOrdered()
	AddAll(set, elem...)
	*s = *set

	return nil
}

func (s *Ordered) IsZero() bool {
	return IsEmpty(s)
}

func (s *Ordered) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte
	switch src.(type) {
	case []byte:
		source = src.([]byte)
	case string:
		source = []byte(src.(string))
	default:
		return errors.New("incompatible type for ordered string set")
	}

	return json.Unmarshal(source, s)
}

func (s *Ordered) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}
