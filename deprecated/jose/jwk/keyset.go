package jwk

import (
	"encoding/json"
	"errors"
	"gopkg.in/square/go-jose.v2"
	"io"
)

var (
	ErrKeyNotFound = errors.New("json web key is not found")
)

// NewKeySet creates a new key set with the given keys.
func NewKeySet(keys ...*Key) *KeySet {
	s := &KeySet{
		index: map[string]*Key{},
		keys:  []*Key{},
	}
	for _, key := range keys {
		s.index[key.KeyID] = key
		s.keys = append(s.keys, key)
	}
	return s
}

// ReadKeySet create new KeySet with data from the reader
func ReadKeySet(reader io.Reader) (*KeySet, error) {
	gojoseJwks := new(jose.JSONWebKeySet)
	err := json.NewDecoder(reader).Decode(&gojoseJwks)
	if err != nil {
		return nil, err
	}

	set := &KeySet{index: map[string]*Key{}, keys: []*Key{}}
	for _, k := range gojoseJwks.Keys {
		/*
		 * !Important!
		 * -----------
		 * copy k value onto local stack
		 * before k changes in the next iteration.
		 */
		jwk := k
		set.index[jwk.KeyID] = (*Key)(&jwk)
		set.keys = append(set.keys, (*Key)(&jwk))
	}

	return set, nil
}

type KeySet struct {
	index map[string]*Key
	keys  []*Key
}

// Count returns the number of keys in the set.
func (s *KeySet) Count() int {
	return len(s.keys)
}

// KeyById finds a Key by its id value.
func (s *KeySet) KeyById(kid string) (*Key, error) {
	k, ok := s.index[kid]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return k, nil
}

// KeyForSigning find a key for signing with the given algorithm. If multiple signing keys with the same
// algorithm exists in the set, the last one is picked.
func (s *KeySet) KeyForSigning(alg string) (*Key, error) {
	var candidates []*Key

	for _, k := range s.keys {
		if k.Use == UseSig && k.Algorithm == alg {
			candidates = append(candidates, k)
		}
	}

	switch len(candidates) {
	case 0:
		return nil, ErrKeyNotFound
	case 1:
		return candidates[0], nil
	default:
		return candidates[len(candidates)-1], nil
	}
}

// KeyForEncryption find a key for encryption with the given algorithm. The returned key may be a private key, in
// which case, caller needs to convert to a public key before use. If multiple encryption keys with the same algorithm
// exists in the set, the last one based on the current time.
func (s *KeySet) KeyForEncryption(alg string) (*Key, error) {
	var candidates []*Key

	for _, k := range s.keys {
		if k.Use == UseEnc && k.Algorithm == alg {
			candidates = append(candidates, k)
		}
	}

	switch len(candidates) {
	case 0:
		return nil, ErrKeyNotFound
	case 1:
		return candidates[0], nil
	default:
		return candidates[len(candidates)-1], nil
	}
}

// ToPublic returns a new KeySet with only public asymmetric keys so that it is read to be shared.
func (s *KeySet) ToPublic() *KeySet {
	var pubKeys []*Key
	for _, each := range s.keys {
		if !each.IsSymmetric() {
			pubKeys = append(pubKeys, each.Public())
		}
	}
	return NewKeySet(pubKeys...)
}

func (s *KeySet) MarshalJSON() ([]byte, error) {
	gojoseJwks := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, each := range s.keys {
		gojoseJwks.Keys = append(gojoseJwks.Keys, (jose.JSONWebKey)(*each))
	}

	raw, err := json.Marshal(gojoseJwks)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (s *KeySet) UnmarshalJSON(bytes []byte) error {
	goJoseJwks := new(jose.JSONWebKeySet)
	if err := json.Unmarshal(bytes, &goJoseJwks); err != nil {
		return err
	}

	s.index = map[string]*Key{}
	s.keys = []*Key{}

	for _, each := range goJoseJwks.Keys {
		jwk := each // copy it onto the stack before referencing it, this is important!
		s.index[each.KeyID] = (*Key)(&jwk)
		s.keys = append(s.keys, (*Key)(&jwk))
	}

	return nil
}
