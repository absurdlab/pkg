package timeplus_test

import (
	"github.com/absurdlab/pkg/timeplus"
	"testing"
)

func TestSecond_SQL(t *testing.T) {
	var s = timeplus.Second(60)

	dv, err := s.Value()
	if err != nil {
		t.Error(err)
	}

	var s2 timeplus.Second
	err = s2.Scan(dv)
	if err != nil {
		t.Error(err)
	}

	if s != s2 {
		t.Log("expected seconds to be equal")
	}
}
