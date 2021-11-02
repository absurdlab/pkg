package timeplus

import "time"

// Now returns the current Timestamp.
func Now() Timestamp {
	return On(time.Now())
}

// On returns the Timestamp at given time.
func On(t time.Time) Timestamp {
	return Timestamp(t.Unix())
}
