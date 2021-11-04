package stderr

import (
	"encoding/json"
	"errors"
)

var (
	// ErrCorruptedView is the default error returned when a View cannot be reconstructed into a chain of errors.
	ErrCorruptedView = Message("error view corrupted")
)

// View is the standard payload to transmit error information over the wire. It is usually constructed
// from an error chain and serialized to the wire format by the sender, and then deserialized and optionally
// reconstructed back to an error chain by the receiver.
type View struct {
	Status  int     `json:"status,omitempty" yaml:"status,omitempty"`
	Code    string  `json:"error,omitempty" yaml:"error,omitempty"`
	Message string  `json:"message,omitempty" yaml:"message,omitempty"`
	Context []*node `json:"context,omitempty" yaml:"context,omitempty"`
}

// With defaults the View with information suggested in the error chain. Traversing down the error chain, the first
// status error is used as View.Status; the first code error is used as View.Code; the first message error is used
// as View.Message. And context is collected by all errors in chain, as long as they don't generate an error during
// collection. By default, this method does not touch the status, code, message and context if they are not zero valued.
func (v *View) With(err error) *View {
	if v.Status == 0 {
		var status *StatusError
		if errors.As(err, &status) {
			v.Status = status.Status()
		}
	}

	if len(v.Code) == 0 {
		var code *CodeError
		if errors.As(err, &code) {
			v.Code = code.Code()
		}
	}

	if len(v.Message) == 0 {
		var message *MessageError
		if errors.As(err, &message) {
			v.Message = message.Message()
		}
	}

	if len(v.Context) == 0 {
		nodes, err := collectNodes(err)
		if err == nil {
			v.Context = nodes
		}
	}

	return v
}

// ToView fills error into a new View using With.
func ToView(err error) *View {
	return new(View).With(err)
}

// FromView attempts to restore the error chain using data from the context. If context is empty, or an error
// occurred during the recovery process, it defaults to FromViewWithoutContext.
func FromView(v *View) error {
	if len(v.Context) == 0 {
		return FromViewWithoutContext(v)
	}

	var chain []error
	{
		for _, n := range v.Context {
			if len(n.Data) == 0 {
				continue
			}

			var target Error

			switch n.Type {
			case typeStatus:
				target = new(StatusError)
			case typeCode:
				target = new(CodeError)
			case typeMessage:
				target = new(MessageError)
			case typeParams:
				target = new(ParamsError)
			case typeGeneric:
				target = new(GenericError)
			default:
				continue
			}

			if e := json.Unmarshal(n.Data, target); e != nil {
				return FromViewWithoutContext(v)
			}

			chain = append(chain, target)
		}
	}

	return Chain(chain...)
}

// FromViewWithoutContext attempts to recover error chain using only status, code and message. If none of these
// are available, it returns an error chain containing ErrCorruptedView.
func FromViewWithoutContext(v *View) error {
	var chain []error
	{
		if v.Status > 0 {
			chain = append(chain, Status(v.Status))
		}

		if len(v.Code) > 0 {
			chain = append(chain, Code(v.Code))
		}

		if len(v.Message) > 0 {
			chain = append(chain, Message(v.Message))
		}
	}

	if len(chain) == 0 {
		return Chain(ErrCorruptedView, errors.New("no status, code or message"))
	}

	return Chain(chain...)
}

type node struct {
	Type string          `json:"type,omitempty" yaml:"type,omitempty"`
	Data json.RawMessage `json:"data,omitempty" yaml:"data,omitempty"`
}

func collectNodes(err error) ([]*node, error) {
	if err == nil {
		return []*node{}, nil
	}

	var results []*node
	for cur := normalize(err); cur != nil; cur = normalize(cur.Unwrap()) {
		n, e := cur.asNode()
		if e != nil {
			return nil, e
		}
		results = append(results, n)
	}

	return results, nil
}
