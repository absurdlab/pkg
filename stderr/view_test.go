package stderr_test

import (
	"errors"
	"github.com/absurdlab/pkg/stderr"
	"testing"
)

func TestView(t *testing.T) {
	err := stderr.Chain(
		stderr.Status(400),
		stderr.Code("invalid_item"),
		stderr.Message("item is not found"),
		errors.New("not found"),
		stderr.Message("another message"),
	)

	view := stderr.ToView(err)
	{
		if actual, expect := view.Status, 400; actual != expect {
			t.Errorf("expect %d, actual %d", expect, actual)
		}
		if actual, expect := view.Code, "invalid_item"; actual != expect {
			t.Errorf("expect %s, actual %s", expect, actual)
		}
		if actual, expect := view.Message, "item is not found"; actual != expect {
			t.Errorf("expect %s, actual %s", expect, actual)
		}
	}

	err2 := stderr.FromView(view)
	{
		if !errors.Is(err2, stderr.Status(400)) {
			t.Error("expect recovered error to have status error in chain")
		}
		if !errors.Is(err2, stderr.Code("invalid_item")) {
			t.Error("expect recovered error to have code error in chain")
		}
	}
}
