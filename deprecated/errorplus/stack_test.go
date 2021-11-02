package errorplus_test

import (
	"errors"
	"fmt"
	"github.com/absurdlab/pkg/errorplus"
	"testing"
)

func TestGetDetailStack(t *testing.T) {
	err := errors.New("root")
	err = errorplus.Wrap(err).DecorateF("first").Field("level", 1)
	err = errorplus.DecorateF("second").Wrap(err).Field("level", 2)

	stack := errorplus.GetDetailStack(err)

	if len(stack) != 3 {
		t.Errorf("stack should be 3 layers")
	}

	if stack[0]["level"] != 2 || stack[0]["error"] != "second" {
		t.Errorf("stack error")
	}

	if stack[1]["level"] != 1 || stack[1]["error"] != "first" {
		t.Errorf("stack error")
	}

	if stack[2]["error"] != "root" {
		t.Errorf("stack error")
	}
}

func TestGetStatusHint(t *testing.T) {
	err := errors.New("root")
	err = errorplus.Wrap(err).DecorateF("first").StatusHint(404)
	err = errorplus.Wrap(err).DecorateF("second")
	err = errorplus.Wrap(err).DecorateF("third").StatusHint(400)
	err = errorplus.Wrap(err).DecorateF("forth")

	if errorplus.GetStatusHint(err) != 400 {
		t.Error("status hint mismatch")
	}
}

func TestGetMessage(t *testing.T) {
	err := errors.New("root")
	err = errorplus.Wrap(err).DecorateF("first")
	err = errorplus.Wrap(err)
	err = errorplus.Wrap(err).DecorateF("second")
	err = errorplus.Wrap(err)

	if "second" != errorplus.GetMessage(err) {
		t.Errorf("get message error")
	}
}

func TestAs(t *testing.T) {
	err := errors.New("root")
	err = errorplus.Decorate(&userError{userID: "foo"}).Wrap(err)
	err = errorplus.Wrap(err).DecorateF("first")

	recoveredErr := new(userError)
	errors.As(err, &recoveredErr)

	if recoveredErr.userID != "foo" {
		t.Error("did not recover error")
	}
}

type userError struct {
	userID string
}

func (e *userError) Error() string {
	return fmt.Sprintf("user %s has error", e.userID)
}
