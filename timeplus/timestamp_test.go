package timeplus_test

import (
	"github.com/absurdlab/pkg/timeplus"
	"testing"
	"time"
)

func TestTimestamp_Before(t *testing.T) {
	now := timeplus.On(time.Now())
	past := timeplus.On(time.Now().Add(-5 * time.Hour))
	future := timeplus.On(time.Now().Add(5 * time.Hour))
	var zero timeplus.Timestamp

	if before := past.Before(now); !before {
		t.Error("expected past is before now")
	}

	if before := now.Before(future); !before {
		t.Error("expected now is before future")
	}

	if before := zero.Before(now); before {
		t.Error("expect zero time to always fail Before")
	}
}

func TestTimestamp_After(t *testing.T) {
	now := timeplus.On(time.Now())
	past := timeplus.On(time.Now().Add(-5 * time.Hour))
	future := timeplus.On(time.Now().Add(5 * time.Hour))
	var zero timeplus.Timestamp

	if after := now.After(past); !after {
		t.Error("expected now is after past")
	}

	if after := future.After(now); !after {
		t.Error("expected future is after now")
	}

	if after := zero.After(now); after {
		t.Error("expect zero time to always fail After")
	}
}

func TestTimestamp_Equals(t *testing.T) {
	t1 := time.Now()

	ts1 := timeplus.On(t1)
	ts2 := timeplus.On(t1)
	ts3 := timeplus.On(time.Now().Add(time.Minute))
	var zero timeplus.Timestamp

	if eq := ts1.Equals(ts2); !eq {
		t.Error("expect ts1 to equal ts2")
	}

	if eq := ts1.Equals(ts3); eq {
		t.Error("expect ts1 to not equal ts3")
	}

	if eq := ts1.Equals(zero); eq {
		t.Error("expect zero to always fail Equals")
	}
}
